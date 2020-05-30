Feature:
  In order to annotate graph
  As an API user
  I need to be able to request graph-annotations

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register graph annotation
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "title": "deploy",
        "description": "link: https://example.com/",
        "from": 1484000000,
        "to": 1484000030,
        "service": "ExampleService",
        "roles": [
          "ExampleRole1",
          "ExampleRole2"
        ]
      }
      """
    When I send "POST" request to "/api/v0/graph-annotations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "title": "deploy",
        "description": "link: https://example.com/",
        "from": 1484000000,
        "to": 1484000030,
        "service": "ExampleService",
        "roles": [
          "ExampleRole1",
          "ExampleRole2"
        ]
      }
      """
    Then I keep the JSON response at "id" as "$ANNOTATION_ID"

  Scenario: should update graph annotation
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "id": "$ANNOTATION_ID",
        "title": "updated deploy",
        "description": "link: https://example.com/",
        "from": 1484000000,
        "to": 1484000030,
        "service": "ExampleService",
        "roles": [
          "ExampleRole1",
          "ExampleRole2"
        ]
      }
      """
    When I send "PUT" request to "/api/v0/graph-annotations/$ANNOTATION_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$ANNOTATION_ID",
        "title": "updated deploy",
        "description": "link: https://example.com/",
        "from": 1484000000,
        "to": 1484000030,
        "service": "ExampleService",
        "roles": [
          "ExampleRole1",
          "ExampleRole2"
        ]
      }
      """

  Scenario: should get graph annotations
    When I send "GET" request to "/api/v0/graph-annotations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "graphAnnotations": [
          {
            "id": "$ANNOTATION_ID",
            "title": "updated deploy",
            "description": "link: https://example.com/",
            "from": 1484000000,
            "to": 1484000030,
            "service": "ExampleService",
            "roles": [
              "ExampleRole1",
              "ExampleRole2"
            ]
          }
        ]
      }
      """

  Scenario: should delete graph annotation
    When I send "DELETE" request to "/api/v0/graph-annotations/$ANNOTATION_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$ANNOTATION_ID",
        "title": "updated deploy",
        "description": "link: https://example.com/",
        "from": 1484000000,
        "to": 1484000030,
        "service": "ExampleService",
        "roles": [
          "ExampleRole1",
          "ExampleRole2"
        ]
      }
      """

  Scenario: should get empty graph annotations
    When I send "GET" request to "/api/v0/graph-annotations"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "graphAnnotations": []
      }
      """
