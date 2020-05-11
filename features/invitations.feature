Feature:
  In order to invite user
  As an API user
  I need to be able to request invitations

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should invite user
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "email": "example@example.com",
        "authority": "viewer"
      }
      """
    When I send "POST" request to "/api/v0/invitations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "email": "example@example.com",
        "authority": "viewer"
      }
      """

  Scenario: should get invitations
    When I send "GET" request to "/api/v0/invitations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "invitations": [
          {
            "email": "example@example.com",
            "authority": "viewer",
            "expiresAt": "@number@"
          }
        ]
      }
      """

  Scenario: should revoke invitations
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "email": "example@example.com"
      }
      """
    When I send "POST" request to "/api/v0/invitations/revoke"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get empty invitations
    When I send "GET" request to "/api/v0/invitations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "invitations": []
      }
      """