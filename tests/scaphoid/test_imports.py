"""Verify the public API surface of the scaphoid package.

IMPORTANT: Our safe-guard, this also fires when the classes change.

These tests exist to catch accidental renames, missing re-exports, and
circular-import regressions. They import every public name from every entry
point and check __all__ lists are accurate and complete.
"""

import importlib
import inspect

# ---------------------------------------------------------------------------
# Top-level package
# ---------------------------------------------------------------------------


class TestPackagePublicApi:
    def test_all_names_are_importable(self):
        import scaphoid

        for name in scaphoid.__all__:
            assert hasattr(scaphoid, name), (
                f"scaphoid.__all__ lists {name!r} but it is not importable"
            )

    def test_all_is_complete(self):
        import scaphoid

        expected = {
            "EnvironmentConfig",
            "EnvironmentMode",
            "ExitCodeRange",
            "ExitCodeType",
            "ExecutionResult",
            "OutputStream",
            "RunConfig",
            "BinaryExecution",
            "StreamSearch",
            "run",
        }
        assert set(scaphoid.__all__) == expected

    def test_environment_importable(self):
        from scaphoid import EnvironmentConfig

        assert inspect.isclass(EnvironmentConfig)

    def test_execution_result_importable(self):
        from scaphoid import ExecutionResult

        assert inspect.isclass(ExecutionResult)

    def test_output_stream_importable(self):
        from scaphoid import OutputStream

        assert inspect.isclass(OutputStream)

    def test_stream_search_importable(self):
        from scaphoid import StreamSearch

        assert inspect.isclass(StreamSearch)

    def test_run_config_importable(self):
        from scaphoid import RunConfig

        assert inspect.isclass(RunConfig)

    def test_shell_importable(self):
        from scaphoid import BinaryExecution

        assert inspect.isclass(BinaryExecution)

    def test_run_importable(self):
        from scaphoid import run

        assert callable(run)


# ---------------------------------------------------------------------------
# scaphoid.environment
# ---------------------------------------------------------------------------


class TestEnvModule:
    def test_module_is_importable(self):
        importlib.import_module("scaphoid.environment")

    def test_all_is_complete(self):
        from scaphoid import environment as env

        assert set(env.__all__) == {"EnvironmentConfig", "EnvironmentMode"}

    def test_environment_importable_directly(self):
        from scaphoid.environment import EnvironmentConfig

        assert inspect.isclass(EnvironmentConfig)

    def test_environment_is_same_object(self):
        from scaphoid import EnvironmentConfig
        from scaphoid.environment import EnvironmentConfig as _Direct

        assert EnvironmentConfig is _Direct


# ---------------------------------------------------------------------------
# scaphoid.streams
# ---------------------------------------------------------------------------


class TestStreamsModule:
    def test_module_is_importable(self):
        importlib.import_module("scaphoid.streams")

    def test_all_is_complete(self):
        from scaphoid import streams

        assert set(streams.__all__) == {"OutputStream", "OutputType", "StreamSearch"}

    def test_all_names_are_importable(self):
        from scaphoid import streams

        for name in streams.__all__:
            assert hasattr(streams, name)

    def test_output_stream_importable_directly(self):
        from scaphoid.streams import OutputStream

        assert inspect.isclass(OutputStream)

    def test_stream_search_importable_directly(self):
        from scaphoid.streams import StreamSearch

        assert inspect.isclass(StreamSearch)


# ---------------------------------------------------------------------------
# scaphoid.result
# ---------------------------------------------------------------------------


class TestResultModule:
    def test_module_is_importable(self):
        importlib.import_module("scaphoid.result")

    def test_all_is_complete(self):
        from scaphoid import result

        assert set(result.__all__) == {
            "ExitCodeRange",
            "ExitCodeType",
            "ExecutionResult",
        }

    def test_execution_result_importable_directly(self):
        from scaphoid.result import ExecutionResult

        assert inspect.isclass(ExecutionResult)

    def test_output_stream_not_re_exported(self):
        # OutputStream lives in streams, not result — keep the boundary clear
        from scaphoid import result

        assert "OutputStream" not in result.__all__


# ---------------------------------------------------------------------------
# scaphoid.binary_execution
# ---------------------------------------------------------------------------


class TestBinaryExecutionModule:
    def test_module_is_importable(self):
        importlib.import_module("scaphoid.binary_execution")

    def test_all_is_complete(self):
        from scaphoid import binary_execution

        assert set(binary_execution.__all__) == {"RunConfig", "BinaryExecution", "run"}

    def test_all_names_are_importable(self):
        from scaphoid import binary_execution

        for name in binary_execution.__all__:
            assert hasattr(binary_execution, name)

    def test_run_config_importable_directly(self):
        from scaphoid.binary_execution import RunConfig

        assert inspect.isclass(RunConfig)

    def test_shell_importable_directly(self):
        from scaphoid.binary_execution import BinaryExecution

        assert inspect.isclass(BinaryExecution)

    def test_run_importable_directly(self):
        from scaphoid.binary_execution import run

        assert callable(run)


# ---------------------------------------------------------------------------
# Identity — top-level re-exports are the same objects as the originals
# ---------------------------------------------------------------------------


class TestReExportIdentity:
    def test_execution_result_is_same_object(self):
        from scaphoid import ExecutionResult
        from scaphoid.result import ExecutionResult as _Direct

        assert ExecutionResult is _Direct

    def test_output_stream_is_same_object(self):
        from scaphoid import OutputStream
        from scaphoid.streams import OutputStream as _Direct

        assert OutputStream is _Direct

    def test_stream_search_is_same_object(self):
        from scaphoid import StreamSearch
        from scaphoid.streams import StreamSearch as _Direct

        assert StreamSearch is _Direct

    def test_shell_is_same_object(self):
        from scaphoid import BinaryExecution
        from scaphoid.binary_execution import BinaryExecution as _Direct

        assert BinaryExecution is _Direct

    def test_run_is_same_object(self):
        from scaphoid import run
        from scaphoid.binary_execution import run as _Direct

        assert run is _Direct
