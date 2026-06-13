@gating
Feature: Tier gating
  Lower tiers ship less code: a higher-tier command is not merely hidden, it is
  absent from the binary, so the CLI reports it as an unknown command.

  This feature must be run against the apprentice (default) binary.

  Scenario Outline: "<command>" is unavailable at the apprentice tier
    When I run sphinx "<command>"
    Then the exit code is not 0
    And stderr contains "unknown command"

    Examples:
      | command |
      | copy    |
      | edit    |
      | card    |
      | topt    |
      | rotate  |
      | stats   |
      | backup  |
      | session |
      | config  |
      | export  |
      | import  |
      | restore |
      | debug   |
      | forge   |
      | riddle  |
