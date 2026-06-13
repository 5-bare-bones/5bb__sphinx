@keeper
Feature: Keeper tier
  Persistent state and stewardship: backup and session.

  Scenario: backup help
    When I run sphinx "backup --help"
    Then the exit code is 0
    And stdout contains "backup"

  Scenario: session help
    When I run sphinx "session --help"
    Then the exit code is 0
    And stdout contains "session"
