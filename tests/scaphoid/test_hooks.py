"""Tests for tests/environment.py — the behave before/after_scenario hooks."""

import os
from types import SimpleNamespace

import pytest

from environment import after_scenario, before_scenario


def make_context(**kwargs):
    return SimpleNamespace(**kwargs)


def run_before(extra=None):
    ctx = make_context(**(extra or {}))
    before_scenario(ctx, scenario=None)
    return ctx


# ---------------------------------------------------------------------------
# before_scenario — field types and initial values
# ---------------------------------------------------------------------------

class TestBeforeScenarioFields:
    def test_returncode_is_none(self):
        assert run_before().returncode is None

    def test_stdout_is_empty_string(self):
        assert run_before().stdout == ""

    def test_stderr_is_empty_string(self):
        assert run_before().stderr == ""

    def test_extra_env_isolates_home_to_tmpdir(self):
        ctx = run_before()
        assert ctx.extra_env == {"HOME": ctx.tmpdir}

    def test_extra_env_is_a_dict(self):
        assert isinstance(run_before().extra_env, dict)

    def test_tmpdir_is_a_string(self):
        ctx = run_before()
        assert isinstance(ctx.tmpdir, str)
        os.rmdir(ctx.tmpdir)

    def test_tmpdir_exists_on_disk(self):
        ctx = run_before()
        assert os.path.isdir(ctx.tmpdir)
        os.rmdir(ctx.tmpdir)

    def test_tmpdir_has_sphinx_prefix(self):
        ctx = run_before()
        assert os.path.basename(ctx.tmpdir).startswith("sphinx_test_")
        os.rmdir(ctx.tmpdir)

    def test_each_call_creates_a_unique_tmpdir(self):
        ctx1 = run_before()
        ctx2 = run_before()
        assert ctx1.tmpdir != ctx2.tmpdir
        os.rmdir(ctx1.tmpdir)
        os.rmdir(ctx2.tmpdir)

    def test_overwrites_existing_returncode(self):
        ctx = make_context(returncode=99)
        before_scenario(ctx, scenario=None)
        assert ctx.returncode is None
        os.rmdir(ctx.tmpdir)


# ---------------------------------------------------------------------------
# after_scenario — cleanup
# ---------------------------------------------------------------------------

class TestAfterScenario:
    def test_removes_tmpdir(self):
        ctx = run_before()
        tmpdir = ctx.tmpdir
        after_scenario(ctx, scenario=None)
        assert not os.path.exists(tmpdir)

    def test_removes_tmpdir_with_files_inside(self):
        ctx = run_before()
        open(os.path.join(ctx.tmpdir, "file.txt"), "w").close()
        after_scenario(ctx, scenario=None)
        assert not os.path.exists(ctx.tmpdir)

    def test_removes_tmpdir_with_subdirectories(self):
        ctx = run_before()
        os.mkdir(os.path.join(ctx.tmpdir, "sub"))
        after_scenario(ctx, scenario=None)
        assert not os.path.exists(ctx.tmpdir)

    def test_noop_when_tmpdir_attribute_missing(self):
        ctx = make_context()
        after_scenario(ctx, scenario=None)  # must not raise

    def test_noop_when_tmpdir_already_deleted(self):
        ctx = run_before()
        os.rmdir(ctx.tmpdir)
        after_scenario(ctx, scenario=None)  # must not raise

    def test_noop_when_tmpdir_is_a_file_not_a_dir(self, tmp_path):
        f = tmp_path / "not_a_dir.txt"
        f.write_text("hello")
        ctx = make_context(tmpdir=str(f))
        after_scenario(ctx, scenario=None)  # must not raise
        assert f.exists()
