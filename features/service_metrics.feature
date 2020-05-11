Feature:
  In order to monitor the service metric values
  As an API user
  I need to be able to request services

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register services
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
       """
      {
        "name": "ExampleService",
        "memo": "This is an example"
      }
      """
    When I send "POST" request to "/api/v0/services"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "name": "ExampleService",
        "memo": "This is an example",
        "roles": []
      }
      """

  Scenario: should register service metric values
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      [
        {
          "name":"cpu",
          "time": 1351700030,
          "value": 1.234
        },
        {
          "name":"memory",
          "time": 1351700050,
          "value": 5.678
        }
      ]
      """
    When I send "POST" request to "/api/v0/services/ExampleService/tsdb"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get service metric values
    When I send "GET" request to "/api/v0/services/ExampleService/metrics?name=cpu&from=1351700000&to=1351700100"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metrics": [
          {
            "name":"cpu",
            "time": 1351700030,
            "value":1.234
          }
        ]
      }
      """

  Scenario: should get service metric names
    When I send "GET" request to "/api/v0/services/ExampleService/metric-names"
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