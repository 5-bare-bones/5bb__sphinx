import os

import pytest

from scaphoid.binary_execution import BinaryExecution, RunConfig, _resolve_env, run
from scaphoid.environment import EnvironmentConfig, EnvironmentMode
from scaphoid.result import ExecutionResult

# ---------------------------------------------------------------------------
# RunConfig — fields and defaults
# ---------------------------------------------------------------------------


class TestRunConfigDefaults:
    def test_path_is_none(self):
        assert RunConfig().path is None

    def test_environment_is_none(self):
        assert RunConfig().environment is None

    def test_timeout_is_30(self):
        assert RunConfig().timeout == 30.0

    def test_shell_is_false(self):
        assert RunConfig().shell is False

    def test_check_is_false(self):
        assert RunConfig().check is False

    def test_encoding_is_utf8(self):
        assert RunConfig().encoding == "utf-8"


class TestRunConfigFields:
    def test_path_stored(self, tmp_path):
        assert RunConfig(path=tmp_path).path == tmp_path

    def test_environment_stored(self):
        env = EnvironmentConfig.merge({"KEY": "val"})
        assert RunConfig(environment=env).environment == env

    def test_timeout_stored(self):
        assert RunConfig(timeout=5.0).timeout == 5.0

    def test_timeout_none_disables_it(self):
        assert RunConfig(timeout=None).timeout is None

    def test_shell_stored(self):
        assert RunConfig(shell=True).shell is True

    def test_check_stored(self):
        assert RunConfig(check=True).check is True

    def test_encoding_stored(self):
        assert RunConfig(encoding="latin-1").encoding == "latin-1"


class TestRunConfigReplace:
    def test_returns_new_instance(self):
        cfg = RunConfig(timeout=10.0)
        assert cfg.replace(timeout=5.0) is not cfg

    def test_original_unchanged(self):
        cfg = RunConfig(timeout=10.0)
        cfg.replace(timeout=5.0)
        assert cfg.timeout == 10.0

    def test_changed_field_applied(self):
        assert RunConfig().replace(timeout=99.0).timeout == 99.0

    def test_unchanged_fields_preserved(self):
        cfg = RunConfig(timeout=5.0, encoding="latin-1")
        assert cfg.replace(timeout=1.0).encoding == "latin-1"


# ---------------------------------------------------------------------------
# _resolve_env — delegates to EnvironmentConfig.resolve()
# ---------------------------------------------------------------------------


class TestResolveEnv:
    def test_none_returns_none(self):
        assert _resolve_env(RunConfig()) is None

    def test_replace_environment_returns_only_given_vars(self):
        env = EnvironmentConfig.replace({"FOO": "bar"})
        result = _resolve_env(RunConfig(environment=env))
        assert result == {"FOO": "bar"}
        assert "PATH" not in result

    def test_merge_environment_merges_os_environ(self):
        env = EnvironmentConfig.merge({"SCAPHOID_RESOLVE_TEST": "yes"})
        result = _resolve_env(RunConfig(environment=env))
        assert result["SCAPHOID_RESOLVE_TEST"] == "yes"
        assert "PATH" in result

    def test_merge_overrides_os_environ_key(self):
        env = EnvironmentConfig.merge({"PATH": "overridden"})
        result = _resolve_env(RunConfig(environment=env))
        assert result["PATH"] == "overridden"


# ---------------------------------------------------------------------------
# run() — command forms
# ---------------------------------------------------------------------------


class TestRunCommandForms:
    def test_list_command(self):
        run(["echo", "hello"]).stdout.contains("hello")

    def test_string_command_is_tokenised(self):
        run("echo hello").stdout.contains("hello")

    def test_string_command_with_shell(self):
        run("echo hello", RunConfig(shell=True)).stdout.contains("hello")


class TestRunOutput:
    def test_captures_stdout(self):
        run(["echo", "captured"]).stdout.contains("captured")

    def test_stdout_empty_on_no_output(self):
        run(["true"]).stdout.is_empty()

    def test_captures_stderr(self):
        run(["sh", "-c", "echo err >&2"]).stderr.contains("err")

    def test_stdin_forwarded(self):
        run(["show"], stdin="piped content\n").stdout.contains("piped content")


class TestRunExitCode:
    def test_zero_on_success(self):
        run(["true"]).succeeded()

    def test_nonzero_on_failure(self):
        run(["false"]).failed()

    def test_specific_exit_code(self):
        run(["sh", "-c", "exit 3"]).has_exit_code(3)


class TestRunMetadata:
    def test_returns_execution_result(self):
        assert isinstance(run(["true"]), ExecutionResult)

    def test_elapsed_ms_non_negative(self):
        assert run(["true"]).elapsed_ms >= 0

    def test_command_field_matches_argv(self):
        assert run(["echo", "hi"]).command == ["echo", "hi"]

    def test_cwd_field_populated(self, tmp_path):
        assert run(["pwd"], RunConfig(path=tmp_path)).path == str(tmp_path)

    def test_cwd_falls_back_to_getcwd(self):
        assert run(["true"]).path == os.getcwd()


class TestRunCheck:
    def test_check_passes_on_success(self):
        run(["true"], RunConfig(check=True))

    def test_check_raises_on_failure(self):
        with pytest.raises(AssertionError, match="expected success"):
            run(["false"], RunConfig(check=True))


class TestRunEnvironment:
    def test_merge_environment_visible_to_process(self):
        env = EnvironmentConfig.merge({"SCAPHOID_RUN_VAR": "present"})
        run(
            ["sh", "-c", "echo $SCAPHOID_RUN_VAR"], RunConfig(environment=env)
        ).stdout.contains("present")

    def test_replace_environment_hides_os_environ(self):
        env = EnvironmentConfig.replace({"PATH": os.environ["PATH"]})
        run(
            ["sh", "-c", "echo ${SCAPHOID_RUN_VAR:-absent}"],
            RunConfig(environment=env),
        ).stdout.contains("absent")


# ---------------------------------------------------------------------------
# Shell — init and attributes
# ---------------------------------------------------------------------------


class TestShellInit:
    def test_default_config_is_run_config(self):
        assert isinstance(BinaryExecution().config, RunConfig)

    def test_provided_config_stored(self):
        cfg = RunConfig(timeout=5.0)
        assert BinaryExecution(cfg).config is cfg

    def test_last_starts_as_none(self):
        assert BinaryExecution().last is None


# ---------------------------------------------------------------------------
# Shell — run
# ---------------------------------------------------------------------------


class TestShellRun:
    def test_returns_execution_result(self):
        assert isinstance(BinaryExecution().run(["true"]), ExecutionResult)

    def test_updates_last(self):
        sh = BinaryExecution()
        r = sh.run(["echo", "hi"])
        assert sh.last is r

    def test_last_overwritten_on_second_call(self):
        sh = BinaryExecution()
        sh.run(["true"])
        r2 = sh.run(["true"])
        assert sh.last is r2

    def test_scoped_to_path(self, tmp_path):
        sh = BinaryExecution(RunConfig(path=tmp_path))
        sh.run(["pwd"]).stdout.contains(str(tmp_path))

    def test_stdin_forwarded(self):
        BinaryExecution().run(["show"], stdin="hello\n").stdout.contains("hello")

    def test_override_does_not_mutate_base_config(self):
        sh = BinaryExecution(RunConfig(timeout=30.0))
        sh.run(["true"], timeout=1.0)
        assert sh.config.timeout == 30.0


# ---------------------------------------------------------------------------
# Shell — set_env
# ---------------------------------------------------------------------------


class TestShellSetEnv:
    def test_variable_visible_to_process(self):
        sh = BinaryExecution()
        sh.set_env("SCAPHOID_SHELL_VAR", "hello")
        sh.run(["sh", "-c", "echo $SCAPHOID_SHELL_VAR"]).stdout.contains("hello")

    def test_accumulates_multiple_keys(self):
        sh = BinaryExecution()
        sh.set_env("VAR_A", "a")
        sh.set_env("VAR_B", "b")
        sh.run(["sh", "-c", "echo $VAR_A $VAR_B"]).stdout.contains("a b")

    def test_stores_environment_object_on_config(self):
        sh = BinaryExecution()
        sh.set_env("X", "1")
        assert isinstance(sh.config.environment, EnvironmentConfig)

    def test_does_not_lose_previous_vars(self):
        sh = BinaryExecution()
        sh.set_env("A", "1")
        sh.set_env("B", "2")
        assert sh.config.environment.variables == {"A": "1", "B": "2"}

    def test_defaults_to_merge_mode(self):
        sh = BinaryExecution()
        sh.set_env("X", "1")
        assert sh.config.environment.mode is EnvironmentMode.MERGE
