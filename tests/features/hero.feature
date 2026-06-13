@hero
Feature: Hero tier (dev build)
  The contributor's dev surface: debug. Compiled only with the dev build tag.

  Scenario: debug help
    When I run sphinx "debug --help"
    Then the exit code is 0
    And stdout contains "decrypt-dump"
