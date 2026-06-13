"""
file:
    name: environment.py
    id: 019e9850-0000-7000-0000-000000000001
    defines: __all__
---
EnvironmentMode and EnvironmentConfig — subprocess environment management.
"""

from __future__ import annotations

import os
from dataclasses import dataclass
from enum import Enum
from typing import Mapping

__all__ = ["EnvironmentMode", "EnvironmentConfig"]


class EnvironmentMode(Enum):
    """How the subprocess environment is constructed relative to os.environ."""

    CLEAN   = "clean"    # completely empty — no variables at all
    REPLACE = "replace"  # exactly the given variables; os.environ ignored
    MERGE   = "merge"    # os.environ + variables (variables win on conflict)
    ALLOW   = "allow"    # selected keys from os.environ + variables
    DENY    = "deny"     # os.environ minus blocked keys + variables


@dataclass(frozen=True)
class EnvironmentConfig:
    """Immutable environment specification for a subprocess call.

    Five modes — choose the factory that matches your intent::

        EnvironmentConfig.clean()
        EnvironmentConfig.replace({"PATH": "/usr/bin"})
        EnvironmentConfig.merge({"DEBUG": "1"})
        EnvironmentConfig.allow({"PATH", "HOME"}, {"PORT": "8080"})
        EnvironmentConfig.deny({"SECRET_KEY"}, {"DEBUG": "1"})
    """

    mode: EnvironmentMode = EnvironmentMode.MERGE
    variables: dict[str, str] | None = None
    allow_list: frozenset[str] | None = None
    deny_list:  frozenset[str] | None = None

    def __post_init__(self) -> None:
        if self.mode is EnvironmentMode.REPLACE and not self.variables:
            raise ValueError("REPLACE mode requires a non-empty variables dict")
        if self.mode is EnvironmentMode.ALLOW and not self.allow_list:
            raise ValueError("ALLOW mode requires allow_list")
        if self.mode is EnvironmentMode.DENY and not self.deny_list:
            raise ValueError("DENY mode requires deny_list")

    # --- factories -----------------------------------------------------------

    @classmethod
    def clean(cls) -> "EnvironmentConfig":
        """Completely empty environment — no variables whatsoever."""
        return cls(mode=EnvironmentMode.CLEAN)

    @classmethod
    def replace(cls, variables: Mapping[str, str]) -> "EnvironmentConfig":
        """Use *variables* as the entire environment; nothing from os.environ."""
        return cls(mode=EnvironmentMode.REPLACE, variables=dict(variables))

    @classmethod
    def merge(cls, variables: Mapping[str, str] | None = None) -> "EnvironmentConfig":
        """os.environ + *variables* (variables win on conflict)."""
        return cls(mode=EnvironmentMode.MERGE, variables=dict(variables or {}))

    @classmethod
    def allow(
        cls,
        keys: set[str],
        variables: Mapping[str, str] | None = None,
    ) -> "EnvironmentConfig":
        """Inherit only *keys* from os.environ, then overlay *variables*."""
        return cls(
            mode=EnvironmentMode.ALLOW,
            variables=dict(variables or {}),
            allow_list=frozenset(keys),
        )

    @classmethod
    def deny(
        cls,
        keys: set[str],
        variables: Mapping[str, str] | None = None,
    ) -> "EnvironmentConfig":
        """Inherit all of os.environ except *keys*, then overlay *variables*."""
        return cls(
            mode=EnvironmentMode.DENY,
            variables=dict(variables or {}),
            deny_list=frozenset(keys),
        )

    # --- mutation (immutable update) -----------------------------------------

    def with_var(self, key: str, value: str) -> "EnvironmentConfig":
        """Return a new EnvironmentConfig with *key* set to *value*."""
        new_vars = dict(self.variables or {})
        new_vars[key] = value
        return EnvironmentConfig(
            mode=self.mode,
            variables=new_vars,
            allow_list=self.allow_list,
            deny_list=self.deny_list,
        )

    # --- resolution ----------------------------------------------------------

    def resolve(self) -> dict[str, str]:
        """Return the final environment dict to pass to subprocess."""
        match self.mode:
            case EnvironmentMode.CLEAN:
                return {}
            case EnvironmentMode.REPLACE:
                return dict(self.variables or {})
            case EnvironmentMode.MERGE:
                result = dict(os.environ)
                result.update(self.variables or {})
                return result
            case EnvironmentMode.ALLOW:
                result = {k: os.environ[k] for k in self.allow_list if k in os.environ}
                result.update(self.variables or {})
                return result
            case EnvironmentMode.DENY:
                result = {k: v for k, v in os.environ.items() if k not in self.deny_list}
                result.update(self.variables or {})
                return result
