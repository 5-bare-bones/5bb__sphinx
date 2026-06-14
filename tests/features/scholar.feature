@scholar
Feature: Scholar tier
  topt (formerly 2fa), card, entry rotate, the vault group (stats), and the
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

  Scenario: entry rotate help
    When I run sphinx "entry rotate --help"
    Then the exit code is 0
    And stdout contains "Rotate an entry"

  Scenario: vault group help lists the scholar subcommand
    When I run sphinx "vault --help"
    Then the exit code is 0
    And stdout contains "Vault administration"
    And stdout contains "Show vault statistics"

  Scenario: vault stats help
    When I run sphinx "vault stats --help"
    Then the exit code is 0
    And stdout contains "vault statistics"

  Scenario: file gains the scholar-tier move and del subcommands
    When I run sphinx "file --help"
    Then the exit code is 0
    And stdout contains "Move a file or directory"
    And stdout contains "Remove files or directories"
