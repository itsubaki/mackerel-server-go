Feature:
  In order to monitor the host metric values
  As an API user
  I need to be able to request hosts

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register host
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "host01",
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
        "id": "@string@"
      }
      """
    Then I keep the JSON response at "id" as "$HOST_ID"

  Scenario: should register host metric values
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      [
        {
          "hostId": "$HOST_ID",
          "name":"cpu",
          "time": 1351700030,
          "value": 1.234
        },
        {
          "hostId": "$HOST_ID",
          "name":"memory",
          "time": 1351700050,
          "value": 5.678
        }
      ]
      """
    When I send "POST" request to "/api/v0/tsdb"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get latest host metric values
    When I send "GET" request to "/api/v0/tsdb/latest"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "tsdbLatest": {
          "$HOST_ID": {
            "cpu": {
              "name":"cpu",
              "value":1.234
            },
            "memory": {
              "name": "memory",
              "value": 5.678
            }
          }
        }
      }
      """

  Scenario: should get host metric values
    When I send "GET" request to "/api/v0/hosts/$HOST_ID/metrics?name=cpu&from=1351700000&to=1351700100"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metrics": [
          {
            "hostId": "$HOST_ID",
            "name": "cpu",
            "time": 1351700030,
            "value": 1.234
          }
        ]
      }
      """

  Scenario: should get host metric names
    When I send "GET" request to "/api/v0/hosts/$HOST_ID/metric-names"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "names": [
          "cpu",
          "memory"
        ]
      }
      """