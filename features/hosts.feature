Feature:
  In order to monitor the host
  As an API user
  I need to be able to request hosts

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should get hosts
    When I send "GET" request to "/api/v0/hosts"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "hosts": []
      }
      """

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

  Scenario: should get host information
    When I send "GET" request to "/api/v0/hosts/$HOST_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "host": {
          "id": "$HOST_ID",
          "name": "host01",
          "status": "working",
          "memo": "",
          "createdAt": "@number@",
          "isRetired": false,
          "roles": {},
          "meta": {
            "agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)",
            "agent-revision": "2f531c6",
            "agent-version": "0.4.2"
          }
        }
      }
      """

  Scenario: should update host information
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "cucumber-host01",
        "meta": {
          "agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)",
          "agent-revision": "2f531c6",
          "agent-version": "0.4.2"
        }
      }
      """
    When I send "PUT" request to "/api/v0/hosts/$HOST_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$HOST_ID"
      }
      """

  Scenario: should update host status
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "status": "poweroff"
      }
      """
    When I send "POST" request to "/api/v0/hosts/$HOST_ID/status"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """

  Scenario: should get updated host information
    When I send "GET" request to "/api/v0/hosts/$HOST_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "host": {
          "id": "$HOST_ID",
          "name": "cucumber-host01",
          "status": "poweroff",
          "memo": "",
          "createdAt": "@number@",
          "isRetired": false,
          "roles": {},
          "meta": {
            "agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)",
            "agent-revision": "2f531c6",
            "agent-version": "0.4.2"
          }
        }
      }
      """

  Scenario: should retire host
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {}
      """
    When I send "POST" request to "/api/v0/hosts/$HOST_ID/retire"
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
        "tsdbLatest": {}
      }
      """
