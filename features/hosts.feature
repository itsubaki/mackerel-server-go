Feature: get hosts
  In order to know hosts information
  As an API user
  I need to be able to request hosts

  Scenario: should get hosts
    Given I fill the X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/hosts"
    Then the response code should be 200
    And the response should match json:
      """
      {"hosts":[]}
      """

  Scenario: should get latest host metric values
    Given I fill the X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/tsdb/latest"
    Then the response code should be 200
    And the response should match json:
      """
      {"tsdbLatest":{}}
      """