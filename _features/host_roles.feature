Feature:
  In order to monitor the host
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

  Scenario: should update host roles
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "roleFullnames": [
          "Hatena-Bookmark:db-master"
        ]
      }
      """
    When I send "PUT" request to "/api/v0/hosts/$HOST_ID/role-fullnames"
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
          "name": "host01",
          "status": "working",
          "memo": "",
          "createdAt": "@number@",
          "isRetired": false,
          "roles": {
            "Hatena-Bookmark": [
              "db-master"
            ]
          },
          "roleFullnames": [
              "Hatena-Bookmark:db-master"
          ],
          "meta": {
            "agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)",
            "agent-revision": "2f531c6",
            "agent-version": "0.4.2"
          }
        }
      }
      """

  Scenario: should delete services
    When I send "DELETE" request to "/api/v0/services/Hatena-Bookmark"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "name": "Hatena-Bookmark",
        "memo": "",
        "roles": [
          "db-master"
        ]
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
