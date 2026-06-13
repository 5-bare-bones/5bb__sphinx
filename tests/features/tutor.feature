@tutor
Feature: Tutor tier
  Vault-minting (forge) and knowledge-derived-key riddles. Compiled only with
  the tutor build tag; these run end to end without a vault.

  Scenario: forge mints a vault from a world manifest
    Given a file "world.yaml" with content
      """
      students:
        - name: alice
          passphrase: correct horse battery staple
          entries:
            - name: github
              username: alice
              password: hunter2
      """
    When I run sphinx "forge world.yaml -o vaults"
    Then the exit code is 0
    And stdout contains "minted 1 vault"
    And a file "vaults/alice.db" exists

  Scenario: forge refuses to overwrite an existing vault
    Given a file "world.yaml" with content
      """
      students:
        - name: alice
          passphrase: pw
      """
    When I run sphinx "forge world.yaml -o vaults"
    Then the exit code is 0
    When I run sphinx "forge world.yaml -o vaults"
    Then the exit code is not 0
    And stderr contains "already exists"

  Scenario: riddle help
    When I run sphinx "riddle --help"
    Then the exit code is 0
    And stdout contains "knowledge-derived-key riddles"

  Scenario: a sealed riddle is revealed only by the correct answer
    When I run sphinx "riddle seal --prompt I-speak-without-a-mouth --answer an-echo --secret tier-token -o echo.riddle"
    Then the exit code is 0
    And a file "echo.riddle" exists
    When I run sphinx "riddle solve echo.riddle" with input
      """
      an-echo
      """
    Then the exit code is 0
    And stdout contains "tier-token"

  Scenario: a wrong answer does not reveal the secret
    When I run sphinx "riddle seal --prompt q --answer the-moon --secret hidden -o moon.riddle"
    Then the exit code is 0
    When I run sphinx "riddle solve moon.riddle" with input
      """
      the-sun
      """
    Then the exit code is not 0
    And stderr contains "wrong answer"
