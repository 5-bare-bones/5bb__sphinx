"""Result — chainable assertion wrapper around a completed subprocess call."""

from __future__ import annotations

import re
import shlex
from dataclasses import dataclass

__all___ = ["ExecutionResult"]


@dataclass
class ExecutionResult:
    """Outcome of one command run, with fluent assertion methods.

    Every assertion returns ``self`` so calls can be chained::

        run("alchemist brew", cfg).succeeded().stdout_contains("Brewed")
    """

    argv: list[str]
    return_code: int
    stdout: str
    stderr: str
    duration: float
    path: str

    # --- introspection -------------------------------------------------------

    @property
    def output(self) -> str:
        """stdout and stderr concatenated, in that order."""
        return self.stdout + self.stderr

    @property
    def ok(self) -> bool:
        return self.return_code == 0

    def _context(self) -> str:
        return (
            f"\n    command : {shlex.join(self.argv)}"
            f"\n    path    : {self.path}"
            f"\n    exit    : {self.return_code}  ({self.duration * 1000:.0f} ms)"
            f"\n    stdout  : {self.stdout!r}"
            f"\n    stderr  : {self.stderr!r}"
        )

    def _fail(self, what: str) -> "ExecutionResult":
        raise AssertionError(what + self._context())

    # --- exit-code assertions ------------------------------------------------

    def succeeded(self) -> "ExecutionResult":
        if self.return_code != 0:
            self._fail(f"expected success, got exit {self.return_code}")
        return self

    def failed(self) -> "ExecutionResult":
        if self.return_code == 0:
            self._fail("expected a non-zero exit, but command succeeded")
        return self

    def exit_code(self, code: int) -> "ExecutionResult":
        if self.return_code != code:
            self._fail(f"expected exit {code}, got {self.return_code}")
        return self

    # --- stdout assertions ---------------------------------------------------

    def stdout_contains(self, text: str) -> "ExecutionResult":
        if text not in self.stdout:
            self._fail(f"stdout does not contain {text!r}")
        return self

    def stdout_excludes(self, text: str) -> "ExecutionResult":
        if text in self.stdout:
            self._fail(f"stdout unexpectedly contains {text!r}")
        return self

    def stdout_matches(self, pattern: str) -> "ExecutionResult":
        if re.search(pattern, self.stdout) is None:
            self._fail(f"stdout does not match /{pattern}/")
        return self

    def stdout_equals(self, text: str, *, strip: bool = True) -> "ExecutionResult":
        actual = self.stdout.strip() if strip else self.stdout
        expected = text.strip() if strip else text
        if actual != expected:
            self._fail(f"stdout != {expected!r}")
        return self

    def stdout_empty(self) -> "ExecutionResult":
        if self.stdout.strip():
            self._fail("expected empty stdout")
        return self

    # --- stderr assertions ---------------------------------------------------

    def stderr_contains(self, text: str) -> "ExecutionResult":
        if text not in self.stderr:
            self._fail(f"stderr does not contain {text!r}")
        return self

    def stderr_empty(self) -> "ExecutionResult":
        if self.stderr.strip():
            self._fail("expected empty stderr")
        return self

    # --- combined output -----------------------------------------------------

    def output_contains(self, text: str) -> "ExecutionResult":
        if text not in self.output:
            self._fail(f"combined output does not contain {text!r}")
        return self
