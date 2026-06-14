@keeper
Feature: Keeper tier
  Persistent state and stewardship: vault backup and session.

  Scenario: vault backup help
    When I run sphinx "vault backup --help"
    Then the exit code is 0
    And stdout contains "backup"

  Scenario: session help
    When I run sphinx "session --help"
    Then the exit code is 0
    And stdout contains "session"
