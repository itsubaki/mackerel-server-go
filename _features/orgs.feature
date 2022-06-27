Feature:
  In order to know org name
  As an API user
  I need to be able to request org

  Scenario: should get org
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/org"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "name": "fixture"
      }
      """

  Scenario: forbidden
    When I send "GET" request to "/api/v0/org"
    Then the response code should be 403
