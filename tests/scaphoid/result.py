"""
file:
    name: result.py
    id: 019e9840-7d74-73f3-901b-6633d079fa23
    defines: __all__
links:
    using:
    - name: streams.py
      id: 019e9865-02c7-7f31-b723-ccb0525fb6fb
      for: [OutputStream, StreamSearch]
    tests:
    - name: test_result.py
      id: 019e9846-fbdb-7c12-a0c2-8504fa6fcdd8
---
ExecutionResult — exit-level assertions over a completed command run.
"""

from __future__ import annotations

import shlex
from dataclasses import dataclass
from enum import Enum
from typing import Callable, Union

from .streams import OutputStream, StreamSearch

__all__ = ["ExitCodeType", "ExitCodeRange", "ExecutionResult"]

# Type alias accepted by succeeded() / failed() — either the standard enum or
# a custom range built with ExitCodeRange.
ExitCodeSpec = "ExitCodeType | ExitCodeRange"


class ExitCodeType(Enum):
    """Exit-code categories with built-in matching logic.

    NOTE: that some executables do not use all results above 0 as error. (e.g. `robot` from robot-framework)
    """

    SUCCESS = "success"  # == 0
    FAILURE = "failure"  # > 0  (app-level error)
    SIGNAL = "signal"  #   < 0  (killed by OS signal)
    RESULT = "result"  #  user defined via ExitCodeRange
    NONZERO = "nonzero"  # != 0 (failure or signal)

    def matches(self, code: int) -> bool:
        """Return True if *code* falls within this category."""
        match self:
            case ExitCodeType.SUCCESS:
                return code == 0
            case ExitCodeType.FAILURE:
                return code > 0
            case ExitCodeType.SIGNAL:
                return code < 0
            case ExitCodeType.NONZERO:
                return code != 0

    @classmethod
    def classify(cls, code: int) -> "ExitCodeType":
        """Return the most specific category for *code*."""
        if code == 0:
            return cls.SUCCESS
        if code > 0:
            return cls.FAILURE
        return cls.SIGNAL


class ExitCodeRange:
    """Custom exit code matcher for tools with non-standard exit code conventions.

    Implements the same ``matches(code)`` protocol as :class:`ExitCodeType` so
    both are accepted wherever an exit-code spec is expected::

        # Robot Framework: exit codes 0–250 are valid test outcomes
        result.succeeded(ExitCodeRange.between(0, 250))

        # assert the runner itself crashed (RF codes 251–255)
        result.failed(ExitCodeRange.above(250))

        # specific known sentinel values
        result.failed(ExitCodeRange.exactly(1, 2))
    """

    __slots__ = ("_predicate", "_label")

    def __init__(self, predicate: "Callable[[int], bool]", label: str) -> None:
        object.__setattr__(self, "_predicate", predicate)
        object.__setattr__(self, "_label", label)

    def __setattr__(self, name: str, value: object) -> None:
        raise AttributeError("ExitCodeRange is immutable")

    # --- factories -----------------------------------------------------------

    @classmethod
    def exactly(cls, *codes: int) -> "ExitCodeRange":
        """Match one or more specific exit codes."""
        frozen = frozenset(codes)
        return cls(lambda c: c in frozen, f"in {set(codes)}")

    @classmethod
    def between(cls, low: int, high: int) -> "ExitCodeRange":
        """Match any code in the closed interval [low, high]."""
        return cls(lambda c: low <= c <= high, f"in [{low}, {high}]")

    @classmethod
    def above(cls, threshold: int) -> "ExitCodeRange":
        """Match any code strictly greater than *threshold*."""
        return cls(lambda c: c > threshold, f"> {threshold}")

    @classmethod
    def below(cls, threshold: int) -> "ExitCodeRange":
        """Match any code strictly less than *threshold*."""
        return cls(lambda c: c < threshold, f"< {threshold}")

    # --- protocol ------------------------------------------------------------

    def matches(self, code: int) -> bool:
        """Return True if *code* satisfies this range."""
        return self._predicate(code)

    def __repr__(self) -> str:
        return f"ExitCodeRange({self._label})"


@dataclass(frozen=True)
class ExecutionResult:
    """Outcome of one command run, plus exit-level assertions."""

    exit_code: int  # the integer code from the binary
    stdout: OutputStream  # standard output (string) in a convenience class
    stderr: OutputStream  # standard error (string) in a convenience class
    elapsed_ms: int  # time used for execution
    command: list[str]  # command with all arguments, options and flags
    path: str | None = (
        None  # the path the executable was executed in (current working dir/cwd)
    )

    # --- factory -------------------------------------------------------------

    @classmethod
    def from_completed(
        cls,
        *,
        exit_code: int,
        stdout: str,
        stderr: str,
        elapsed_ms: int,
        command: list[str],
        path: str | None = None,
    ) -> "ExecutionResult":
        """Build from raw strings, wrapping each stream with a label."""
        return cls(
            exit_code=exit_code,
            stdout=OutputStream(stdout, "stdout"),
            stderr=OutputStream(stderr, "stderr"),
            elapsed_ms=elapsed_ms,
            command=list(command),
            path=path,
        )

    # --- introspection -------------------------------------------------------

    @property
    def ok(self) -> bool:
        return self.exit_code == 0

    @property
    def output(self) -> StreamSearch:
        """Cross-stream view for assertions spanning stdout and stderr."""
        return StreamSearch(self.stdout, self.stderr)

    # --- failure helper ------------------------------------------------------

    def _context(self) -> str:
        return (
            f"\n    command : {shlex.join(self.command)}"
            f"\n    path    : {self.path}"
            f"\n    exit    : {self.exit_code}  ({self.elapsed_ms} ms)"
            f"\n    stdout  : {self.stdout.content!r}"
            f"\n    stderr  : {self.stderr.content!r}"
        )

    def _fail(self, what: str) -> "ExecutionResult":
        raise AssertionError(what + self._context())

    # --- exit-level assertions (chainable) -----------------------------------

    def succeeded(
        self, spec: Union[ExitCodeType, "ExitCodeRange"] = ExitCodeType.SUCCESS
    ) -> "ExecutionResult":
        """Assert the process exited with a code matching *spec*."""
        if not spec.matches(self.exit_code):
            label = spec.value if isinstance(spec, ExitCodeType) else repr(spec)
            self._fail(f"expected {label} exit, got {self.exit_code}")
        return self

    def failed(
        self, spec: Union[ExitCodeType, "ExitCodeRange"] = ExitCodeType.NONZERO
    ) -> "ExecutionResult":
        """Assert the process exited with the given kind of non-success code."""
        if spec is ExitCodeType.SUCCESS:
            raise ValueError("ExitCodeType.SUCCESS cannot be used with failed()")
        if not spec.matches(self.exit_code):
            label = spec.value if isinstance(spec, ExitCodeType) else repr(spec)
            self._fail(f"expected {label} exit, got {self.exit_code}")
        return self

    def has_exit_code(self, code: int) -> "ExecutionResult":
        """check for a specific exit code"""
        if self.exit_code != code:
            self._fail(f"expected exit {code}, got {self.exit_code}")
        return self
