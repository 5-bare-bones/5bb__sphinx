import os
import pytest
from scaphoid.environment import EnvironmentConfig, EnvironmentMode


# ---------------------------------------------------------------------------
# EnvironmentMode enum
# ---------------------------------------------------------------------------

class TestEnvironmentMode:
    def test_clean(self):   assert EnvironmentMode.CLEAN.value   == "clean"
    def test_replace(self): assert EnvironmentMode.REPLACE.value == "replace"
    def test_merge(self):   assert EnvironmentMode.MERGE.value   == "merge"
    def test_allow(self):   assert EnvironmentMode.ALLOW.value   == "allow"
    def test_deny(self):    assert EnvironmentMode.DENY.value    == "deny"


# ---------------------------------------------------------------------------
# CLEAN
# ---------------------------------------------------------------------------

class TestClean:
    def test_mode(self):
        assert EnvironmentConfig.clean().mode is EnvironmentMode.CLEAN

    def test_resolve_returns_empty(self):
        assert EnvironmentConfig.clean().resolve() == {}

    def test_os_environ_not_included(self):
        assert "PATH" not in EnvironmentConfig.clean().resolve()


# ---------------------------------------------------------------------------
# REPLACE
# ---------------------------------------------------------------------------

class TestReplace:
    def test_mode(self):
        assert EnvironmentConfig.replace({"A": "1"}).mode is EnvironmentMode.REPLACE

    def test_variables_stored(self):
        assert EnvironmentConfig.replace({"A": "1"}).variables == {"A": "1"}

    def test_requires_non_empty_variables(self):
        with pytest.raises(ValueError, match="REPLACE mode requires"):
            EnvironmentConfig(mode=EnvironmentMode.REPLACE, variables=None)

    def test_empty_dict_also_raises(self):
        with pytest.raises(ValueError, match="REPLACE mode requires"):
            EnvironmentConfig(mode=EnvironmentMode.REPLACE, variables={})

    def test_resolve_returns_only_given_vars(self):
        assert EnvironmentConfig.replace({"A": "1"}).resolve() == {"A": "1"}

    def test_os_environ_not_included(self):
        assert "PATH" not in EnvironmentConfig.replace({"A": "1"}).resolve()

    def test_resolve_returns_copy(self):
        cfg = EnvironmentConfig.replace({"A": "1"})
        cfg.resolve()["NEW"] = "x"
        assert "NEW" not in cfg.resolve()


# ---------------------------------------------------------------------------
# MERGE
# ---------------------------------------------------------------------------

class TestMerge:
    def test_mode(self):
        assert EnvironmentConfig.merge().mode is EnvironmentMode.MERGE

    def test_no_args_stores_empty_variables(self):
        assert EnvironmentConfig.merge().variables == {}

    def test_variables_stored(self):
        assert EnvironmentConfig.merge({"K": "v"}).variables == {"K": "v"}

    def test_resolve_includes_os_environ(self):
        assert "PATH" in EnvironmentConfig.merge().resolve()

    def test_variables_overlay_os_environ(self):
        assert EnvironmentConfig.merge({"PATH": "x"}).resolve()["PATH"] == "x"

    def test_extra_vars_added(self):
        result = EnvironmentConfig.merge({"SCAPHOID_MERGE": "yes"}).resolve()
        assert result["SCAPHOID_MERGE"] == "yes"
        assert "PATH" in result


# ---------------------------------------------------------------------------
# ALLOW
# ---------------------------------------------------------------------------

class TestAllow:
    def test_mode(self):
        assert EnvironmentConfig.allow({"PATH"}).mode is EnvironmentMode.ALLOW

    def test_allow_list_stored_as_frozenset(self):
        cfg = EnvironmentConfig.allow({"PATH", "HOME"})
        assert isinstance(cfg.allow_list, frozenset)
        assert cfg.allow_list == frozenset({"PATH", "HOME"})

    def test_requires_allow_list(self):
        with pytest.raises(ValueError, match="ALLOW mode requires"):
            EnvironmentConfig(mode=EnvironmentMode.ALLOW, allow_list=None)

    def test_resolve_includes_only_allowed_keys(self):
        result = EnvironmentConfig.allow({"PATH"}).resolve()
        assert "PATH" in result

    def test_resolve_excludes_non_allowed_keys(self):
        result = EnvironmentConfig.allow({"PATH"}).resolve()
        for key in result:
            assert key == "PATH" or key in (EnvironmentConfig.allow({"PATH"}).variables or {})

    def test_missing_allowed_key_silently_omitted(self):
        result = EnvironmentConfig.allow({"SCAPHOID_NONEXISTENT"}).resolve()
        assert "SCAPHOID_NONEXISTENT" not in result

    def test_variables_overlaid_on_allowed_base(self):
        result = EnvironmentConfig.allow({"PATH"}, {"EXTRA": "yes"}).resolve()
        assert result["EXTRA"] == "yes"
        assert "PATH" in result

    def test_variables_override_allowed_key(self):
        result = EnvironmentConfig.allow({"PATH"}, {"PATH": "overridden"}).resolve()
        assert result["PATH"] == "overridden"


# ---------------------------------------------------------------------------
# DENY
# ---------------------------------------------------------------------------

class TestDeny:
    def test_mode(self):
        assert EnvironmentConfig.deny({"SECRET"}).mode is EnvironmentMode.DENY

    def test_deny_list_stored_as_frozenset(self):
        cfg = EnvironmentConfig.deny({"SECRET", "TOKEN"})
        assert isinstance(cfg.deny_list, frozenset)
        assert cfg.deny_list == frozenset({"SECRET", "TOKEN"})

    def test_requires_deny_list(self):
        with pytest.raises(ValueError, match="DENY mode requires"):
            EnvironmentConfig(mode=EnvironmentMode.DENY, deny_list=None)

    def test_resolve_excludes_denied_keys(self):
        result = EnvironmentConfig.deny({"PATH"}).resolve()
        assert "PATH" not in result

    def test_resolve_includes_non_denied_keys(self):
        result = EnvironmentConfig.deny({"SCAPHOID_NONEXISTENT"}).resolve()
        assert "PATH" in result

    def test_variables_overlaid_on_deny_base(self):
        result = EnvironmentConfig.deny({"SCAPHOID_NONEXISTENT"}, {"EXTRA": "yes"}).resolve()
        assert result["EXTRA"] == "yes"

    def test_variables_can_restore_denied_key(self):
        result = EnvironmentConfig.deny({"PATH"}, {"PATH": "restored"}).resolve()
        assert result["PATH"] == "restored"


# ---------------------------------------------------------------------------
# with_var
# ---------------------------------------------------------------------------

class TestWithVar:
    def test_returns_new_instance(self):
        cfg = EnvironmentConfig.merge({"A": "1"})
        assert cfg.with_var("B", "2") is not cfg

    def test_new_key_added(self):
        cfg = EnvironmentConfig.merge().with_var("A", "1")
        assert cfg.variables["A"] == "1"

    def test_existing_key_overridden(self):
        cfg = EnvironmentConfig.merge({"A": "old"}).with_var("A", "new")
        assert cfg.variables["A"] == "new"

    def test_original_unchanged(self):
        cfg = EnvironmentConfig.merge({"A": "1"})
        cfg.with_var("B", "2")
        assert "B" not in cfg.variables

    def test_preserves_mode(self):
        assert EnvironmentConfig.deny({"X"}).with_var("A", "1").mode is EnvironmentMode.DENY

    def test_preserves_allow_list(self):
        cfg = EnvironmentConfig.allow({"PATH"}).with_var("X", "1")
        assert cfg.allow_list == frozenset({"PATH"})

    def test_preserves_deny_list(self):
        cfg = EnvironmentConfig.deny({"SECRET"}).with_var("X", "1")
        assert cfg.deny_list == frozenset({"SECRET"})

    def test_chaining(self):
        cfg = EnvironmentConfig.merge().with_var("A", "1").with_var("B", "2")
        assert cfg.variables == {"A": "1", "B": "2"}


# ---------------------------------------------------------------------------
# immutability
# ---------------------------------------------------------------------------

class TestImmutability:
    def test_mode_cannot_be_reassigned(self):
        with pytest.raises(Exception):
            EnvironmentConfig.merge().mode = EnvironmentMode.REPLACE

    def test_variables_cannot_be_reassigned(self):
        with pytest.raises(Exception):
            EnvironmentConfig.merge({"A": "1"}).variables = {}
