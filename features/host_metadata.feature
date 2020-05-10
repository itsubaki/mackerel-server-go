Feature:
  In order to monitor the host metadata
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

  Scenario: should register host metadata
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "message": "this is host metadata"
      }
      """
    When I send "PUT" request to "/api/v0/hosts/$HOST_ID/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get host metadata
    When I send "GET" request to "/api/v0/hosts/$HOST_ID/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "message": "this is host metadata"
      }
      """

  Scenario: should get host metadata list
    When I send "GET" request to "/api/v0/hosts/$HOST_ID/metadata"
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

  Scenario: should delete host metadata
    When I send "DELETE" request to "/api/v0/hosts/$HOST_ID/metadata/example"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get empty host metadata list
    When I send "GET" request to "/api/v0/hosts/$HOST_ID/metadata"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "metadata": []
      }
      """
