@apprentice
Feature: Apprentice tier (default build)
  The commands every sphinx binary ships with: add, list, gen, clear and the
  interactive walker. Available in every tier.

  Scenario: Root help shows usage and the apprentice commands
    When I run sphinx "--help"
    Then the exit code is 0
    And stdout contains "Usage:"
    And stdout contains "Available Commands:"
    And stdout contains "Add an entry"
    And stdout contains "List entries"
    And stdout contains "Generate a random password"
    And stdout contains "Clear clipboard"

  Scenario: Running with no arguments shows help
    When I run sphinx with no arguments
    Then the exit code is 0
    And the output contains "Usage:"

  Scenario: Version reports the build tier
    When I run sphinx "--version"
    Then the exit code is 0
    And stdout contains "tier:"

  Scenario: add help
    When I run sphinx "add --help"
    Then the exit code is 0
    And stdout contains "Add an entry"

  Scenario: list help
    When I run sphinx "list --help"
    Then the exit code is 0
    And stdout contains "List entries"

  Scenario: gen help
    When I run sphinx "gen --help"
    Then the exit code is 0
    And stdout contains "Generate a random password"

  Scenario: gen produces a password without a vault
    When I run sphinx "gen -l 12"
    Then the exit code is 0
    And stdout contains "Password:"

  Scenario: clear help
    When I run sphinx "clear --help"
    Then the exit code is 0
    And stdout contains "Clear clipboard"
