@scholar
Feature: Scholar tier
  The full object model: topt (formerly 2fa), card, rotate, stats, and the
  destructive file subcommands move/del.

  Scenario: topt help
    When I run sphinx "topt --help"
    Then the exit code is 0
    And stdout contains "two-factor authentication"

  Scenario: the 2fa alias still resolves to topt
    When I run sphinx "2fa --help"
    Then the exit code is 0
    And stdout contains "two-factor authentication"

  Scenario: card group help
    When I run sphinx "card --help"
    Then the exit code is 0
    And stdout contains "Card operations"

  Scenario: rotate help
    When I run sphinx "rotate --help"
    Then the exit code is 0
    And stdout contains "Rotate an entry"

  Scenario: stats help
    When I run sphinx "stats --help"
    Then the exit code is 0
    And stdout contains "database statistics"

  Scenario: file gains the scholar-tier move and del subcommands
    When I run sphinx "file --help"
    Then the exit code is 0
    And stdout contains "Move a file or directory"
    And stdout contains "Remove files or directories"
