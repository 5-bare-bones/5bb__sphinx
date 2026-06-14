@adept
Feature: Adept tier
  Entry modification (copy, edit, del) and the file object.

  Scenario: entry copy help
    When I run sphinx "entry copy --help"
    Then the exit code is 0
    And stdout contains "Copy entry credentials"

  Scenario: entry edit help
    When I run sphinx "entry edit --help"
    Then the exit code is 0
    And stdout contains "Edit an entry"

  Scenario: entry del help
    When I run sphinx "entry del --help"
    Then the exit code is 0
    And stdout contains "Remove entries"

  Scenario: file group help lists the adept subcommands
    When I run sphinx "file --help"
    Then the exit code is 0
    And stdout contains "File operations"
    And stdout contains "Add files"
    And stdout contains "Read file"
    And stdout contains "List files"
    And stdout contains "Create stored files"

  Scenario: file add help
    When I run sphinx "file add --help"
    Then the exit code is 0
