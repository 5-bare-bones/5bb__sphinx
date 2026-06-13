"""
file:
    name: test_streams.py
    id: 019e986b-1372-7862-982e-0bfdf9744afc
links:
    code:
    - name: streams.py
      id: 019e9865-02c7-7f31-b723-ccb0525fb6fb
---
"""

import pytest

from scaphoid.streams import OutputStream, StreamSearch

# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------


def out(content="", label="stdout"):
    return OutputStream(content, label)


def search(stdout="", stderr=""):
    return StreamSearch(out(stdout, "stdout"), out(stderr, "stderr"))


# ---------------------------------------------------------------------------
# OutputType — enum
# ---------------------------------------------------------------------------


# ---------------------------------------------------------------------------
# OutputStream — fields
# ---------------------------------------------------------------------------


class TestOutputStreamFields:
    def test_content_stored(self):
        assert out("hello").content == "hello"

    def test_label_stored(self):
        assert OutputStream("x", "stderr").label == "stderr"

    def test_label_defaults_to_output(self):
        assert OutputStream("x").label == "output"

    def test_frozen_content(self):
        with pytest.raises(Exception):
            o = out("hello")
            o.content = "changed"

    def test_frozen_label(self):
        with pytest.raises(Exception):
            o = out("hello")
            o.label = "changed"


# ---------------------------------------------------------------------------
# OutputStream — _fail message format
# ---------------------------------------------------------------------------


class TestOutputStreamFail:
    def test_message_contains_label(self):
        with pytest.raises(AssertionError, match="stderr"):
            OutputStream("x", "stderr")._fail("boom")

    def test_message_contains_what(self):
        with pytest.raises(AssertionError, match="boom"):
            out()._fail("boom")

    def test_message_contains_content(self):
        with pytest.raises(AssertionError, match="my content"):
            out("my content")._fail("something wrong")


# ---------------------------------------------------------------------------
# OutputStream — contains
# ---------------------------------------------------------------------------


class TestOutputStreamContains:
    def test_happy_substring(self):
        out("hello world").contains("world")

    def test_happy_full_string(self):
        out("exact").contains("exact")

    def test_chains(self):
        o = out("foo bar baz")
        assert o.contains("foo").contains("baz") is o

    def test_error_absent_needle(self):
        with pytest.raises(AssertionError, match="does not contain"):
            out("hello").contains("missing")

    def test_error_empty_content(self):
        with pytest.raises(AssertionError, match="does not contain"):
            out("").contains("anything")


# ---------------------------------------------------------------------------
# OutputStream — excludes
# ---------------------------------------------------------------------------


class TestOutputStreamExcludes:
    def test_happy_absent(self):
        out("hello").excludes("world")

    def test_happy_empty_content(self):
        out("").excludes("anything")

    def test_chains(self):
        o = out("foo")
        assert o.excludes("bar") is o

    def test_error_needle_present(self):
        with pytest.raises(AssertionError, match="unexpectedly contains"):
            out("hello world").excludes("world")

    def test_error_needle_is_full_content(self):
        with pytest.raises(AssertionError, match="unexpectedly contains"):
            out("exact").excludes("exact")


# ---------------------------------------------------------------------------
# OutputStream — matches
# ---------------------------------------------------------------------------


class TestOutputStreamMatches:
    def test_happy_simple_pattern(self):
        out("Error: file not found").matches(r"Error:")

    def test_happy_pattern_with_groups(self):
        out("exit code 42").matches(r"exit code (\d+)")

    def test_chains(self):
        o = out("abc 123")
        assert o.matches(r"\d+") is o

    def test_error_no_match(self):
        with pytest.raises(AssertionError, match="does not match"):
            out("hello").matches(r"\d+")

    def test_error_partial_pattern_mismatch(self):
        with pytest.raises(AssertionError, match="does not match"):
            out("Error: ").matches(r"Error: \w+")


# ---------------------------------------------------------------------------
# OutputStream — equals
# ---------------------------------------------------------------------------


class TestOutputStreamEquals:
    def test_happy_exact(self):
        out("hello").equals("hello")

    def test_happy_strips_by_default(self):
        out("  hello\n").equals("hello")

    def test_happy_no_strip(self):
        out("hello").equals("hello", strip=False)

    def test_chains(self):
        o = out("hi")
        assert o.equals("hi") is o

    def test_error_mismatch(self):
        with pytest.raises(AssertionError, match="!="):
            out("hello").equals("world")

    def test_error_whitespace_matters_without_strip(self):
        with pytest.raises(AssertionError):
            out("  hello  ").equals("hello", strip=False)


# ---------------------------------------------------------------------------
# OutputStream — is_empty
# ---------------------------------------------------------------------------


class TestOutputStreamIsEmpty:
    def test_happy_empty_string(self):
        out("").is_empty()

    def test_happy_whitespace_only(self):
        out("   \n\t  ").is_empty()

    def test_chains(self):
        o = out("")
        assert o.is_empty() is o

    def test_error_has_content(self):
        with pytest.raises(AssertionError, match="expected empty"):
            out("something").is_empty()

    def test_error_single_char(self):
        with pytest.raises(AssertionError, match="expected empty"):
            out("x").is_empty()


# ---------------------------------------------------------------------------
# OutputStream — has_line_count
# ---------------------------------------------------------------------------


class TestOutputStreamHasLineCount:
    def test_happy_zero_lines(self):
        out("").has_line_count(0)

    def test_happy_one_line(self):
        out("hello").has_line_count(1)

    def test_happy_multiple_lines(self):
        out("a\nb\nc").has_line_count(3)

    def test_chains(self):
        o = out("x")
        assert o.has_line_count(1) is o

    def test_error_too_few(self):
        with pytest.raises(AssertionError, match="expected 3 lines, found 1"):
            out("hello").has_line_count(3)

    def test_error_too_many(self):
        with pytest.raises(AssertionError, match="expected 1 lines, found 3"):
            out("a\nb\nc").has_line_count(1)


# ---------------------------------------------------------------------------
# OutputStream — accessors
# ---------------------------------------------------------------------------


class TestOutputStreamLines:
    def test_single_line(self):
        assert out("hello").lines() == ["hello"]

    def test_multiple_lines(self):
        assert out("a\nb\nc").lines() == ["a", "b", "c"]

    def test_empty_string(self):
        assert out("").lines() == []

    def test_trailing_newline_not_counted(self):
        assert out("a\nb\n").lines() == ["a", "b"]


class TestOutputStreamLineCount:
    def test_zero(self):
        assert out("").line_count() == 0

    def test_one(self):
        assert out("hello").line_count() == 1

    def test_many(self):
        assert out("a\nb\nc\nd").line_count() == 4


class TestOutputStreamDunder:
    def test_contains_operator_present(self):
        assert "hello" in out("hello world")

    def test_contains_operator_absent(self):
        assert "missing" not in out("hello world")

    def test_str_returns_content(self):
        assert str(out("hello")) == "hello"

    def test_str_empty(self):
        assert str(out("")) == ""


# ---------------------------------------------------------------------------
# StreamSearch — fields
# ---------------------------------------------------------------------------


class TestStreamSearchFields:
    def test_stdout_stored(self):
        o = out("from stdout", "stdout")
        s = StreamSearch(o, out("", "stderr"))
        assert s.stdout is o

    def test_stderr_stored(self):
        e = out("from stderr", "stderr")
        s = StreamSearch(out("", "stdout"), e)
        assert s.stderr is e

    def test_frozen(self):
        with pytest.raises(Exception):
            s = search("a", "b")
            s.stdout = out("other")


# ---------------------------------------------------------------------------
# StreamSearch — _fail message format
# ---------------------------------------------------------------------------


class TestStreamSearchFail:
    def test_message_contains_stdout_content(self):
        with pytest.raises(AssertionError, match="stdout content"):
            search("stdout content", "stderr content")._fail("boom")

    def test_message_contains_stderr_content(self):
        with pytest.raises(AssertionError, match="stderr content"):
            search("stdout content", "stderr content")._fail("boom")

    def test_message_contains_what(self):
        with pytest.raises(AssertionError, match="boom"):
            search()._fail("boom")


# ---------------------------------------------------------------------------
# StreamSearch — contains
# ---------------------------------------------------------------------------


class TestStreamSearchContains:
    def test_happy_found_in_stdout(self):
        search("Task started", "").contains("Task started")

    def test_happy_found_in_stderr(self):
        search("", "not a git repository").contains("not a git repository")

    def test_happy_found_in_both(self):
        search("word", "word").contains("word")

    def test_chains(self):
        s = search("foo", "bar")
        assert s.contains("foo").contains("bar") is s

    def test_error_absent_from_both(self):
        with pytest.raises(AssertionError, match="not found in stdout or stderr"):
            search("foo", "bar").contains("missing")

    def test_error_empty_streams(self):
        with pytest.raises(AssertionError, match="not found"):
            search().contains("anything")


# ---------------------------------------------------------------------------
# StreamSearch — excludes
# ---------------------------------------------------------------------------


class TestStreamSearchExcludes:
    def test_happy_absent_from_both(self):
        search("foo", "bar").excludes("missing")

    def test_happy_empty_streams(self):
        search().excludes("anything")

    def test_chains(self):
        s = search("foo", "bar")
        assert s.excludes("missing") is s

    def test_error_present_in_stdout(self):
        with pytest.raises(AssertionError, match="unexpectedly found"):
            search("secret", "").excludes("secret")

    def test_error_present_in_stderr(self):
        with pytest.raises(AssertionError, match="unexpectedly found"):
            search("", "secret").excludes("secret")


# ---------------------------------------------------------------------------
# StreamSearch — in_stdout_or_stderr
# ---------------------------------------------------------------------------


class TestStreamSearchInStdoutOrStderr:
    def test_in_stdout_only(self):
        assert search("hello", "world").in_stdout_or_stderr("hello") == (True, False)

    def test_in_stderr_only(self):
        assert search("hello", "world").in_stdout_or_stderr("world") == (False, True)

    def test_in_both(self):
        assert search("x", "x").in_stdout_or_stderr("x") == (True, True)

    def test_in_neither(self):
        assert search("foo", "bar").in_stdout_or_stderr("missing") == (False, False)
