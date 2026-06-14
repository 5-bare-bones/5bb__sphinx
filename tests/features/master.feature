@master
Feature: Master tier
  Vault administration: config plus the vault import/export/restore subcommands.

  Scenario: config help
    When I run sphinx "config --help"
    Then the exit code is 0
    And stdout contains "configuration file"

  Scenario: vault import help
    When I run sphinx "vault import --help"
    Then the exit code is 0
    And stdout contains "Import entries"

  Scenario: vault export help
    When I run sphinx "vault export --help"
    Then the exit code is 0
    And stdout contains "Export entries"

  Scenario: vault restore help
    When I run sphinx "vault restore --help"
    Then the exit code is 0
    And stdout contains "Restore the vault"
