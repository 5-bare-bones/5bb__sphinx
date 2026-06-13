"""
file:
    name: test_result.py
    id: 019e9846-fbdb-7c12-a0c2-8504fa6fcdd8
links:
    code:
    - name: result.py
      id: 019e9840-7d74-73f3-901b-6633d079fa23
---
"""

import pytest

from scaphoid.result import ExitCodeRange, ExitCodeType, ExecutionResult
from scaphoid.streams import OutputStream, StreamSearch

# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------


def make(**kwargs):
    defaults = dict(
        exit_code=0,
        stdout="",
        stderr="",
        elapsed_ms=1,
        command=["alchemist"],
        path="/tmp",
    )
    defaults.update(kwargs)
    return ExecutionResult.from_completed(**defaults)


# ---------------------------------------------------------------------------
# from_completed — all fields
# ---------------------------------------------------------------------------


class TestFromCompleted:
    def test_exit_code_stored(self):
        assert make(exit_code=42).exit_code == 42

    def test_elapsed_ms_stored(self):
        assert make(elapsed_ms=250).elapsed_ms == 250

    def test_command_stored(self):
        assert make(command=["alchemist", "brew"]).command == ["alchemist", "brew"]

    def test_cwd_stored(self):
        assert make(path="/some/dir").path == "/some/dir"

    def test_cwd_defaults_to_none(self):
        r = ExecutionResult.from_completed(
            exit_code=0, stdout="", stderr="", elapsed_ms=0, command=["x"]
        )
        assert r.path is None

    def test_stdout_wrapped_as_output_stream(self):
        r = make(stdout="hello")
        assert isinstance(r.stdout, OutputStream)
        assert r.stdout.content == "hello"
        assert r.stdout.label == "stdout"

    def test_stderr_wrapped_as_output_stream(self):
        r = make(stderr="oops")
        assert isinstance(r.stderr, OutputStream)
        assert r.stderr.content == "oops"
        assert r.stderr.label == "stderr"

    def test_command_is_copied(self):
        cmd = ["alchemist", "start"]
        r = make(command=cmd)
        cmd.append("--extra")
        assert r.command == ["alchemist", "start"]

    def test_frozen(self):
        r = make()
        with pytest.raises(Exception):
            r.exit_code = 1


# ---------------------------------------------------------------------------
# ok property
# ---------------------------------------------------------------------------


class TestOk:
    def test_true_when_zero(self):
        assert make(exit_code=0).ok is True

    def test_false_when_one(self):
        assert make(exit_code=1).ok is False

    def test_false_when_nonzero(self):
        assert make(exit_code=127).ok is False


# ---------------------------------------------------------------------------
# output property
# ---------------------------------------------------------------------------


class TestOutput:
    def test_returns_stream_search(self):
        assert isinstance(make().output, StreamSearch)

    def test_stdout_delegate(self):
        make(stdout="found it").output.contains("found it")

    def test_stderr_delegate(self):
        make(stderr="error here").output.contains("error here")

    def test_each_call_returns_fresh_object(self):
        r = make()
        assert r.output is not r.output


# ---------------------------------------------------------------------------
# _context format
# ---------------------------------------------------------------------------


class TestContext:
    def _ctx(self, **kwargs):
        return make(**kwargs)._context()

    def test_contains_command(self):
        assert "alchemist brew" in self._ctx(command=["alchemist", "brew"])

    def test_contains_cwd(self):
        assert "/my/dir" in self._ctx(path="/my/dir")

    def test_contains_exit_code(self):
        assert "7" in self._ctx(exit_code=7)

    def test_contains_elapsed_ms(self):
        assert "42" in self._ctx(elapsed_ms=42)

    def test_contains_stdout(self):
        assert "brewed" in self._ctx(stdout="brewed")

    def test_contains_stderr(self):
        assert "failed hard" in self._ctx(stderr="failed hard")


# ---------------------------------------------------------------------------
# _fail
# ---------------------------------------------------------------------------


class TestFail:
    def test_raises_assertion_error(self):
        with pytest.raises(AssertionError):
            make()._fail("something broke")

    def test_message_includes_what(self):
        with pytest.raises(AssertionError, match="something broke"):
            make()._fail("something broke")

    def test_message_includes_context(self):
        with pytest.raises(AssertionError, match="alchemist"):
            make(command=["alchemist"])._fail("x")


# ---------------------------------------------------------------------------
# succeeded
# ---------------------------------------------------------------------------


class TestSucceeded:
    def test_happy_zero_exit(self):
        make(exit_code=0).succeeded()

    def test_chains(self):
        r = make(exit_code=0)
        assert r.succeeded() is r

    def test_error_nonzero_exit(self):
        with pytest.raises(AssertionError, match="expected success exit, got 1"):
            make(exit_code=1).succeeded()

    def test_error_message_includes_command(self):
        with pytest.raises(AssertionError, match="alchemist start"):
            make(exit_code=1, command=["alchemist", "start"]).succeeded()

    def test_error_message_includes_stdout(self):
        with pytest.raises(AssertionError, match="some output"):
            make(exit_code=2, stdout="some output").succeeded()


# ---------------------------------------------------------------------------
# failed
# ---------------------------------------------------------------------------


class TestFailed:
    def test_happy_nonzero_exit(self):
        make(exit_code=1).failed()

    def test_happy_any_nonzero(self):
        make(exit_code=127).failed()

    def test_chains(self):
        r = make(exit_code=1)
        assert r.failed() is r

    def test_error_zero_exit(self):
        with pytest.raises(AssertionError, match="expected nonzero exit"):
            make(exit_code=0).failed()

    def test_error_message_includes_stderr(self):
        with pytest.raises(AssertionError, match="error text"):
            make(exit_code=0, stderr="error text").failed()


# ---------------------------------------------------------------------------
# has_exit_code
# ---------------------------------------------------------------------------


# ---------------------------------------------------------------------------
# ExitCodeType
# ---------------------------------------------------------------------------

class TestExitCodeTypeMatches:
    def test_success_matches_zero(self):
        assert ExitCodeType.SUCCESS.matches(0) is True

    def test_success_rejects_positive(self):
        assert ExitCodeType.SUCCESS.matches(1) is False

    def test_success_rejects_negative(self):
        assert ExitCodeType.SUCCESS.matches(-15) is False

    def test_failure_matches_positive(self):
        assert ExitCodeType.FAILURE.matches(1) is True
        assert ExitCodeType.FAILURE.matches(127) is True

    def test_failure_rejects_zero(self):
        assert ExitCodeType.FAILURE.matches(0) is False

    def test_failure_rejects_negative(self):
        assert ExitCodeType.FAILURE.matches(-15) is False

    def test_signal_matches_negative(self):
        assert ExitCodeType.SIGNAL.matches(-15) is True
        assert ExitCodeType.SIGNAL.matches(-1) is True

    def test_signal_rejects_zero(self):
        assert ExitCodeType.SIGNAL.matches(0) is False

    def test_signal_rejects_positive(self):
        assert ExitCodeType.SIGNAL.matches(1) is False

    def test_nonzero_matches_positive(self):
        assert ExitCodeType.NONZERO.matches(1) is True

    def test_nonzero_matches_negative(self):
        assert ExitCodeType.NONZERO.matches(-15) is True

    def test_nonzero_rejects_zero(self):
        assert ExitCodeType.NONZERO.matches(0) is False


class TestExitCodeTypeClassify:
    def test_zero_is_success(self):
        assert ExitCodeType.classify(0) is ExitCodeType.SUCCESS

    def test_positive_is_failure(self):
        assert ExitCodeType.classify(1) is ExitCodeType.FAILURE
        assert ExitCodeType.classify(127) is ExitCodeType.FAILURE

    def test_negative_is_signal(self):
        assert ExitCodeType.classify(-15) is ExitCodeType.SIGNAL
        assert ExitCodeType.classify(-1) is ExitCodeType.SIGNAL

    def test_classify_never_returns_nonzero(self):
        for code in (0, 1, 2, -1, -15):
            assert ExitCodeType.classify(code) is not ExitCodeType.NONZERO


class TestExitCodeRange:
    # --- exactly -------------------------------------------------------------

    def test_exactly_matches_listed_code(self):
        assert ExitCodeRange.exactly(1, 2).matches(1)
        assert ExitCodeRange.exactly(1, 2).matches(2)

    def test_exactly_rejects_unlisted_code(self):
        assert not ExitCodeRange.exactly(1, 2).matches(3)

    def test_exactly_single_value(self):
        assert ExitCodeRange.exactly(42).matches(42)
        assert not ExitCodeRange.exactly(42).matches(0)

    # --- between -------------------------------------------------------------

    def test_between_matches_low_bound(self):
        assert ExitCodeRange.between(0, 250).matches(0)

    def test_between_matches_high_bound(self):
        assert ExitCodeRange.between(0, 250).matches(250)

    def test_between_matches_midpoint(self):
        assert ExitCodeRange.between(0, 250).matches(125)

    def test_between_rejects_below_low(self):
        assert not ExitCodeRange.between(1, 250).matches(0)

    def test_between_rejects_above_high(self):
        assert not ExitCodeRange.between(0, 250).matches(251)

    # --- above ---------------------------------------------------------------

    def test_above_matches_value_over_threshold(self):
        assert ExitCodeRange.above(250).matches(251)
        assert ExitCodeRange.above(250).matches(255)

    def test_above_rejects_threshold_itself(self):
        assert not ExitCodeRange.above(250).matches(250)

    def test_above_rejects_below_threshold(self):
        assert not ExitCodeRange.above(250).matches(0)

    # --- below ---------------------------------------------------------------

    def test_below_matches_value_under_threshold(self):
        assert ExitCodeRange.below(0).matches(-1)
        assert ExitCodeRange.below(0).matches(-15)

    def test_below_rejects_threshold_itself(self):
        assert not ExitCodeRange.below(0).matches(0)

    def test_below_rejects_above_threshold(self):
        assert not ExitCodeRange.below(0).matches(1)

    # --- immutability --------------------------------------------------------

    def test_is_immutable(self):
        r = ExitCodeRange.exactly(1)
        with pytest.raises(AttributeError):
            r._label = "changed"

    # --- repr ----------------------------------------------------------------

    def test_repr_contains_label(self):
        assert "250" in repr(ExitCodeRange.above(250))

    # --- used with succeeded / failed ----------------------------------------

    def test_succeeded_with_range(self):
        make(exit_code=1).succeeded(ExitCodeRange.between(0, 2))

    def test_succeeded_with_range_fails_outside(self):
        with pytest.raises(AssertionError, match=r"\[0, 2\]"):
            make(exit_code=3).succeeded(ExitCodeRange.between(0, 2))

    def test_failed_with_range(self):
        make(exit_code=251).failed(ExitCodeRange.above(250))

    def test_failed_with_range_fails_inside(self):
        with pytest.raises(AssertionError, match="> 250"):
            make(exit_code=1).failed(ExitCodeRange.above(250))


class TestFailedWithKind:
    def test_failure_kind_passes_on_positive(self):
        make(exit_code=1).failed(ExitCodeType.FAILURE)

    def test_failure_kind_error_on_zero(self):
        with pytest.raises(AssertionError, match="failure"):
            make(exit_code=0).failed(ExitCodeType.FAILURE)

    def test_failure_kind_error_on_negative(self):
        with pytest.raises(AssertionError, match="failure"):
            make(exit_code=-15).failed(ExitCodeType.FAILURE)

    def test_signal_kind_passes_on_negative(self):
        make(exit_code=-15).failed(ExitCodeType.SIGNAL)

    def test_signal_kind_error_on_zero(self):
        with pytest.raises(AssertionError, match="signal"):
            make(exit_code=0).failed(ExitCodeType.SIGNAL)

    def test_signal_kind_error_on_positive(self):
        with pytest.raises(AssertionError, match="signal"):
            make(exit_code=1).failed(ExitCodeType.SIGNAL)

    def test_success_kind_raises_value_error(self):
        with pytest.raises(ValueError):
            make(exit_code=0).failed(ExitCodeType.SUCCESS)


# ---------------------------------------------------------------------------
# has_exit_code
# ---------------------------------------------------------------------------

class TestHasExitCode:
    def test_happy_zero(self):
        make(exit_code=0).has_exit_code(0)

    def test_happy_specific_code(self):
        make(exit_code=2).has_exit_code(2)

    def test_chains(self):
        r = make(exit_code=0)
        assert r.has_exit_code(0) is r

    def test_error_wrong_code(self):
        with pytest.raises(AssertionError, match="expected exit 2, got 1"):
            make(exit_code=1).has_exit_code(2)

    def test_error_message_includes_context(self):
        with pytest.raises(AssertionError, match="alchemist"):
            make(exit_code=0, command=["alchemist"]).has_exit_code(1)


# ---------------------------------------------------------------------------
# cross-level chaining
# ---------------------------------------------------------------------------


class TestChaining:
    def test_succeeded_into_stdout(self):
        make(exit_code=0, stdout="Task started: My Task").succeeded().stdout.contains(
            "Task started"
        )

    def test_failed_into_stderr(self):
        make(exit_code=1, stderr="not a git repository").failed().stderr.contains(
            "not a git repository"
        )

    def test_has_exit_code_into_output(self):
        make(exit_code=0, stdout="ok", stderr="warn").has_exit_code(0).output.contains(
            "ok"
        )
