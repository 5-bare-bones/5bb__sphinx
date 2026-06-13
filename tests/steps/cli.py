"""
Step definitions for the sphinx CLI behave suite.

The binary under test is taken from the SPHINX_BIN environment variable so the
same scenarios can run against any tier build (bin/sphinx_apprentice ..
bin/sphinx_tutor). It defaults to "sphinx" on PATH.
"""

import os
import shlex

from behave import given, then, when

from scaphoid import EnvironmentConfig, RunConfig, BinaryExecution


def _binary():
    """Resolve the binary under test to an absolute path (cwd is per-scenario)."""
    b = os.environ.get("SPHINX_BIN", "sphinx")
    if os.sep in b or (os.altsep and os.altsep in b):
        return os.path.abspath(b)
    return b  # bare name → resolved on PATH by the subprocess


BINARY = _binary()


def _run(context, args, stdin=None):
    env = EnvironmentConfig.merge(context.extra_env) if context.extra_env else None
    cfg = RunConfig(path=context.tmpdir, environment=env)
    result = BinaryExecution(cfg).run([BINARY, *args], stdin=stdin)
    context.returncode = result.exit_code
    context.stdout = result.stdout.content
    context.stderr = result.stderr.content


# ---------------------------------------------------------------------------
# Given
# ---------------------------------------------------------------------------


@given('a file "{name}" with content')
def step_file_with_content(context, name):
    path = os.path.join(context.tmpdir, name)
    os.makedirs(os.path.dirname(path) or context.tmpdir, exist_ok=True)
    with open(path, "w") as f:
        f.write(context.text)


@given('the environment variable "{name}" is set to "{value}"')
def step_set_env(context, name, value):
    context.extra_env[name] = value


# ---------------------------------------------------------------------------
# When
# ---------------------------------------------------------------------------


@when('I run sphinx "{command}"')
def step_run(context, command):
    _run(context, shlex.split(command) if command else [])


@when("I run sphinx with no arguments")
def step_run_noargs(context):
    _run(context, [])


@when('I run sphinx "{command}" with input')
def step_run_with_input(context, command):
    stdin = context.text if context.text.endswith("\n") else context.text + "\n"
    _run(context, shlex.split(command) if command else [], stdin=stdin)


# ---------------------------------------------------------------------------
# Then
# ---------------------------------------------------------------------------


@then("the exit code is {code:d}")
def step_exit_code(context, code):
    assert context.returncode == code, (
        f"Expected exit {code}, got {context.returncode}\n"
        f"stdout: {context.stdout!r}\nstderr: {context.stderr!r}"
    )


@then("the exit code is not {code:d}")
def step_exit_code_not(context, code):
    assert context.returncode != code, (
        f"Expected exit code other than {code}\n"
        f"stdout: {context.stdout!r}\nstderr: {context.stderr!r}"
    )


@then('stdout contains "{text}"')
def step_stdout_contains(context, text):
    assert text in context.stdout, (
        f"Expected stdout to contain {text!r}\n"
        f"stdout: {context.stdout!r}\nstderr: {context.stderr!r}"
    )


@then('stdout does not contain "{text}"')
def step_stdout_not_contains(context, text):
    assert text not in context.stdout, (
        f"Expected stdout NOT to contain {text!r}\nstdout: {context.stdout!r}"
    )


@then('stderr contains "{text}"')
def step_stderr_contains(context, text):
    assert text in context.stderr, (
        f"Expected stderr to contain {text!r}\nstderr: {context.stderr!r}"
    )


@then('the output contains "{text}"')
def step_output_contains(context, text):
    combined = context.stdout + context.stderr
    assert text in combined, (
        f"Expected output to contain {text!r}\n"
        f"stdout: {context.stdout!r}\nstderr: {context.stderr!r}"
    )


@then('a file "{name}" exists')
def step_file_exists(context, name):
    path = os.path.join(context.tmpdir, name)
    assert os.path.exists(path), f"Expected file {name!r} to exist under the scenario dir"
