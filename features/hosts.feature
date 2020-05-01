Feature: get hosts
  In order to know hosts information
  As an API user
  I need to be able to request hosts

  Scenario: should get hosts
    Given I set X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/hosts"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "hosts":[]
      }
      """

  Scenario: should get latest host metric values
    Given I set X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/tsdb/latest"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "tsdbLatest":{}
      }
      """

  Scenario: should register host
    Given I set X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    Given I set Content-Type header with "application/json"
    Given I set request body with:
      """
      {
        "name":"host01",
        "meta": {
          "agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)",
          "agent-revision": "2f531c6",
          "agent-version": "0.4.2"
        }
      }
      """
    When I send "POST" request to "/api/v0/hosts"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "408a5aaf1a2"
      }
      """

  Scenario: should get host status
    Given I set X-Api-Key header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"
    When I send "GET" request to "/api/v0/hosts/408a5aaf1a2"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "host":{
          "id":"408a5aaf1a2",
          "name":"host01",
          "status":"working",
          "memo":"",
          "createdAt":1588311448,
          "isRetired":false,
          "roles":{},
          "meta":{
            "agent-name":"mackerel-agent/0.27.0 (Revision dfbccea)",
            "agent-revision":"2f531c6",
            "agent-version":"0.4.2"
          }
        }
      }
      """