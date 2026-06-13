@master
Feature: Master tier
  Vault administration: config, export, import, restore.

  Scenario: config help
    When I run sphinx "config --help"
    Then the exit code is 0
    And stdout contains "configuration file"

  Scenario: export help
    When I run sphinx "export --help"
    Then the exit code is 0
    And stdout contains "Export entries"

  Scenario: import help
    When I run sphinx "import --help"
    Then the exit code is 0
    And stdout contains "Import entries"

  Scenario: restore help
    When I run sphinx "restore --help"
    Then the exit code is 0
    And stdout contains "Restore the database"
