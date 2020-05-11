Feature:
  In order to monitor the role
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

  Scenario: should register role
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "ExampleRole",
        "memo": "This is an example"
      }
      """
    When I send "POST" request to "/api/v0/services/ExampleService/roles"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "name": "ExampleRole",
        "memo": "This is an example"
      }
      """

  Scenario: should register role metadata
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "message": "this is role metadata"
      }
      """
    When I send "PUT" request to "/api/v0/services/ExampleService/roles/ExampleRole/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get role metadata
    When I send "GET" request to "/api/v0/services/ExampleService/roles/ExampleRole/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "message": "this is role metadata"
      }
      """

  Scenario: should get role metadata list
    When I send "GET" request to "/api/v0/services/ExampleService/roles/ExampleRole/metadata"
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

  Scenario: should delete role metadata
    When I send "DELETE" request to "/api/v0/services/ExampleService/roles/ExampleRole/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get empty role metadata list
    When I send "GET" request to "/api/v0/services/ExampleService/roles/ExampleRole/metadata"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metadata": []
      }
      """
