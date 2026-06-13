@adept
Feature: Adept tier
  Modifying entries and the file object: copy, edit, del, file.

  Scenario: copy help
    When I run sphinx "copy --help"
    Then the exit code is 0
    And stdout contains "Copy entry credentials"

  Scenario: edit help
    When I run sphinx "edit --help"
    Then the exit code is 0
    And stdout contains "Edit an entry"

  Scenario: del help
    When I run sphinx "del --help"
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
