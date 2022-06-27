Feature:
  In order to monitor the service
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

  Scenario: should get services
    When I send "GET" request to "/api/v0/services"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "services": [
          {
            "name": "ExampleService",
            "memo": "This is an example",
            "roles": []
          }
        ]
      }
      """

  Scenario: should delete services
    When I send "DELETE" request to "/api/v0/services/ExampleService"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "name": "ExampleService",
        "memo": "This is an example",
        "roles": []
      }
      """

  Scenario: should get empty services
    When I send "GET" request to "/api/v0/services"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "services": []
      }
      """
