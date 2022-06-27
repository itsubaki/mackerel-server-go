Feature:
  In order to monitor the service metadata
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


  Scenario: should register service metadata
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "message": "this is service metadata"
      }
      """
    When I send "PUT" request to "/api/v0/services/ExampleService/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get service metadata
    When I send "GET" request to "/api/v0/services/ExampleService/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "message": "this is service metadata"
      }
      """

  Scenario: should get service metadata list
    When I send "GET" request to "/api/v0/services/ExampleService/metadata"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metadata": [
          {
            "namespace": "example"
          }
        ]
      }
      """

  Scenario: should delete service metadata
    When I send "DELETE" request to "/api/v0/services/ExampleService/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get empty host metadata list
    When I send "GET" request to "/api/v0/services/ExampleService/metadata"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metadata": []
      }
      """
