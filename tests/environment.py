"""
Behave environment hooks for sphinx.

Each scenario gets a fresh temporary working directory that also doubles as the
process HOME, so sphinx creates its config and vault under <tmpdir>/.sphinx and
never touches the developer's real ~/.sphinx.
"""

import os
import shutil
import tempfile


def before_scenario(context, scenario):
    context.returncode = None
    context.stdout = ""
    context.stderr = ""
    context.tmpdir = tempfile.mkdtemp(prefix="sphinx_test_")
    # Isolate every sphinx invocation: HOME points at the scratch dir, so the
    # config (.sphinx/sphinx.yaml) and vault (.sphinx/sphinx.db) are created there.
    context.extra_env = {"HOME": context.tmpdir}


def after_scenario(context, scenario):
    if hasattr(context, "tmpdir") and os.path.isdir(context.tmpdir):
        shutil.rmtree(context.tmpdir)
