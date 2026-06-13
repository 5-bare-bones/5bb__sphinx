"""
file:
    name: streams.py
    id: 019e9865-02c7-7f31-b723-ccb0525fb6fb
    defines: __all__
links:
    using:
    - none
    tests:
    - name: test_result.py
      id: 019e986b-1372-7862-982e-0bfdf9744afc
---
OutputStream and StreamSearch — stream-level assertion wrappers.
"""

from __future__ import annotations

import re
from dataclasses import dataclass
from enum import Enum

__all__ = ["OutputStream", "StreamSearch", "OutputType"]


class OutputType(Enum):
    """defining some output types.
    TODO: use OutputType in the classes
    TODO: add tests for OutputType
    """

    # System defaults:
    STDOUT = "stdout"
    STDERR = "stderr"
    FULL = "full"


@dataclass(frozen=True)
class OutputStream:
    """A single captured stream and the assertions that read it."""

    content: str  # the captured data
    label: str = "output"  # "stdout" / "stderr" — used in failure messages

    def _fail(self, what: str) -> "OutputStream":
        raise AssertionError(
            f"{self.label}: {what}\n    {self.label} was: {self.content!r}"
        )

    # --- assertions (chainable) ----------------------------------------------

    def contains(self, needle: str) -> "OutputStream":
        """"""
        if needle not in self.content:
            self._fail(f"does not contain {needle!r}")
        return self

    def excludes(self, needle: str) -> "OutputStream":
        if needle in self.content:
            self._fail(f"unexpectedly contains {needle!r}")
        return self

    def matches(self, pattern: str) -> "OutputStream":
        if re.search(pattern, self.content) is None:
            self._fail(f"does not match /{pattern}/")
        return self

    def equals(self, expected: str, *, strip: bool = True) -> "OutputStream":
        actual = self.content.strip() if strip else self.content
        wanted = expected.strip() if strip else expected
        if actual != wanted:
            self._fail(f"!= {wanted!r}")
        return self

    def is_empty(self) -> "OutputStream":
        if self.content.strip():
            self._fail("expected empty")
        return self

    def has_line_count(self, n: int) -> "OutputStream":
        actual = self.line_count()
        if actual != n:
            self._fail(f"expected {n} lines, found {actual}")
        return self

    # --- accessors -----------------------------------------------------------

    def lines(self) -> list[str]:
        return self.content.splitlines()

    def line_count(self) -> int:
        return len(self.lines())

    def __contains__(self, needle: str) -> bool:
        return needle in self.content

    def __str__(self) -> str:
        return self.content


@dataclass(frozen=True)
class StreamSearch:
    """Assertions that span both stdout and stderr."""

    stdout: OutputStream
    stderr: OutputStream

    def _fail(self, what: str) -> "StreamSearch":
        raise AssertionError(
            f"{what}\n    stdout was: {self.stdout.content!r}"
            f"\n    stderr was: {self.stderr.content!r}"
        )

    def contains(self, needle: str) -> "StreamSearch":
        if needle not in self.stdout.content and needle not in self.stderr.content:
            self._fail(f"{needle!r} not found in stdout or stderr")
        return self

    def excludes(self, needle: str) -> "StreamSearch":
        if needle in self.stdout.content or needle in self.stderr.content:
            self._fail(f"{needle!r} unexpectedly found in stdout or stderr")
        return self

    def in_stdout_or_stderr(self, needle: str) -> tuple[bool, bool]:
        """Return ``(in_stdout, in_stderr)`` for conditional logic."""
        return (needle in self.stdout.content, needle in self.stderr.content)
