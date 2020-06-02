Feature:
  In order to manage user
  As an API user
  I need to be able to request users

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should get empty users
    When I send "GET" request to "/api/v0/users"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "users": []
      }
      """

  Scenario: should get users
    Given the following users exist:
      | org_id      | id      | screen_name | email               |
      | 4b825dc642c | 1234568 | example     | example@example.com |
    When I send "GET" request to "/api/v0/users"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "users": [
          {
            "id": "1234568",
            "screenName": "example",
            "email": "example@example.com",
            "authority": "owner",
            "isInRegistrationProcess": true,
            "isMFAEnabled": true,
            "authenticationMethods": [
              "google"
             ],
            "joinedAt": "@number@"
          }
        ]
      }
      """

  Scenario: should delete user
    Given the following users exist:
      | org_id      | id      | screen_name | email               |
      | 4b825dc642c | 9876543 | example     | example@example.com |
    When I send "DELETE" request to "/api/v0/users/9876543"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "9876543",
        "screenName": "example",
        "email": "example@example.com",
        "authority": "owner",
        "isInRegistrationProcess": true,
        "isMFAEnabled": true,
        "authenticationMethods": [
          "google"
        ],
        "joinedAt": "@number@"
      }
      """