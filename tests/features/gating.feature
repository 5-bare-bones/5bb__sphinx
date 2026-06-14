@gating
Feature: Tier gating
  Lower tiers ship less code: a higher-tier command is not merely hidden, it is
  absent from the binary, so the CLI reports it as an unknown command.

  This feature must be run against the apprentice (default) binary.

  Scenario Outline: top-level "<command>" is unavailable at the apprentice tier
    When I run sphinx "<command>"
    Then the exit code is not 0
    And stderr contains "unknown command"

    Examples:
      | command |
      | file    |
      | card    |
      | topt    |
      | vault   |
      | session |
      | config  |
      | debug   |
      | forge   |
      | riddle  |

  Scenario: higher-tier entry verbs are absent from the apprentice entry group
    When I run sphinx "entry --help"
    Then the exit code is 0
    And stdout contains "Add an entry"
    And stdout does not contain "Copy entry credentials"
    And stdout does not contain "Edit an entry"
    And stdout does not contain "Rotate an entry"
