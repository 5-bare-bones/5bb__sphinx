"""binary_execution — subprocess runner for CLI integration tests.

Why not raw ``subprocess.run``:

* Does **not** raise on a non-zero exit by default, so a test can assert on
  failures just as easily as on successes.
* Results are :class:`~scaphoid.result.ExecutionResult` objects with chainable
  assertion methods and rich failure messages.
* Never uses ``shell=True`` implicitly; string commands are tokenised with
  ``shlex.split``.
* :class:`~scaphoid.environment.EnvironmentConfig` keeps env-var injection explicit and
  immutable, so one test cannot pollute another.
* A timeout by default so one hung binary can't wedge the whole run.

Typical use inside a behave step::

    from scaphoid import BinaryExecution, RunConfig, EnvironmentConfig

    context.sh = BinaryExecution(RunConfig(path=tmpdir))
    context.sh.run("alchemist brew") \\
              .succeeded() \\
              .stdout.contains("Brewed")
"""

from __future__ import annotations

import os
import shlex
import subprocess
import time
from dataclasses import dataclass
from dataclasses import replace as _replace
from typing import Sequence

from .environment import EnvironmentConfig
from .result import ExecutionResult

__all__ = ["RunConfig", "BinaryExecution", "run"]


@dataclass(frozen=True)
class RunConfig:
    """How a command is run. Immutable; derive variants with ``.replace(...)``."""

    path: str | os.PathLike | None = None
    environment: EnvironmentConfig | None = None  # None → inherit from parent process
    timeout: float | None = 30.0
    shell: bool = False
    check: bool = False
    encoding: str = "utf-8"

    def replace(self, **changes) -> "RunConfig":
        return _replace(self, **changes)


_DEFAULT = RunConfig()
MS_PER_SECONDS = 1000


def _resolve_env(config: RunConfig) -> dict[str, str] | None:
    if config.environment is not None:
        return config.environment.resolve()
    return None


def run(
    command: str | Sequence[str],
    config: RunConfig | None = None,
    *,
    stdin: str | None = None,
) -> ExecutionResult:
    """Run *command* under *config* and return an :class:`ExecutionResult`.

    *command* may be a string (tokenised with ``shlex.split`` unless
    ``config.shell`` is set) or an already-split sequence.
    """
    cfg = config or _DEFAULT

    if cfg.shell:
        argv = command if isinstance(command, str) else shlex.join(command)
        display = [argv] if isinstance(argv, str) else list(argv)
    else:
        argv = shlex.split(command) if isinstance(command, str) else list(command)
        display = list(argv)

    start = time.monotonic()
    completed = subprocess.run(
        argv,
        cwd=os.fspath(cfg.path) if cfg.path is not None else None,
        env=_resolve_env(cfg),
        input=stdin,
        capture_output=True,
        text=True,
        encoding=cfg.encoding,
        timeout=cfg.timeout,
        shell=cfg.shell,
    )
    duration = time.monotonic() - start

    result = ExecutionResult.from_completed(
        exit_code=completed.returncode,
        stdout=completed.stdout or "",
        stderr=completed.stderr or "",
        elapsed_ms=int(duration * MS_PER_SECONDS),
        command=display,
        path=os.fspath(cfg.path) if cfg.path is not None else os.getcwd(),
    )
    if cfg.check:
        result.succeeded()
    return result


class BinaryExecution:
    """A session that remembers a :class:`RunConfig` and the last :class:`ExecutionResult`.

    Designed to live on behave's ``context``::

        context.sh = BinaryExecution(RunConfig(path=tmpdir))

    Steps call ``context.sh.run(...)`` and the result is both returned and
    stored as ``context.sh.last`` for subsequent Then-step assertions.
    Per-call deviations are passed as keyword overrides::

        context.sh.run("alchemist start", timeout=5)
    """

    def __init__(self, config: RunConfig | None = None) -> None:
        self.config = config or RunConfig()
        self.last: ExecutionResult | None = None

    def set_env(self, key: str, value: str) -> None:
        """Add a single environment variable, extending the current environment."""
        base = self.config.environment or EnvironmentConfig.merge()
        self.config = self.config.replace(environment=base.with_var(key, value))

    def run(
        self,
        command: str | Sequence[str],
        *,
        stdin: str | None = None,
        **overrides,
    ) -> ExecutionResult:
        cfg = self.config.replace(**overrides) if overrides else self.config
        self.last = run(command, cfg, stdin=stdin)
        return self.last
