Feature:
  In order to manage user
  As an API user
  I need to be able to request users

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should get users
    When I send "GET" request to "/api/v0/users"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "users": []
      }
      """

