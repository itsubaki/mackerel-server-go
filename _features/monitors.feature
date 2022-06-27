Feature:
  In order to monitor the hosts, services
  As an API user
  I need to be able to request monitors

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register monitor
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "type": "host",
        "name": "disk.aa-00.writes.delta",
        "memo": "This monitor is for Hatena Blog.",
        "duration": 3,
        "metric": "disk.aa-00.writes.delta",
        "operator": ">", "warning": 20000.0,
        "critical": 400000.0, "maxCheckAttempts": 3,
        "notificationInterval": 60,
        "scopes": [ "Hatena-Blog" ],
        "excludeScopes": [ "Hatena-Bookmark: db-master" ]
      }
      """
    When I send "POST" request to "/api/v0/monitors"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "type": "host",
        "name": "disk.aa-00.writes.delta",
        "memo": "This monitor is for Hatena Blog.",
        "duration": 3,
        "metric": "disk.aa-00.writes.delta",
        "operator": ">", "warning": 20000.0,
        "critical": 400000.0, "maxCheckAttempts": 3,
        "notificationInterval": 60,
        "scopes": [ "Hatena-Blog" ],
        "excludeScopes": [ "Hatena-Bookmark: db-master" ]
      }
      """
    Then I keep the JSON response at "id" as "$MONITOR_ID"

  Scenario: should get monitors
    When I send "GET" request to "/api/v0/monitors"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "monitors": [
          {
            "id": "$MONITOR_ID",
            "type": "host",
            "name": "disk.aa-00.writes.delta",
            "memo": "This monitor is for Hatena Blog.",
            "duration": 3,
            "metric": "disk.aa-00.writes.delta",
            "operator": ">", "warning": 20000.0,
            "critical": 400000.0, "maxCheckAttempts": 3,
            "notificationInterval": 60,
            "scopes": [ "Hatena-Blog" ],
            "excludeScopes": [ "Hatena-Bookmark: db-master" ]
          }
        ]
      }
      """

  Scenario: should delete monitor
    When I send "DELETE" request to "/api/v0/monitors/$MONITOR_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$MONITOR_ID",
        "type": "host",
        "name": "disk.aa-00.writes.delta",
        "memo": "This monitor is for Hatena Blog.",
        "duration": 3,
        "metric": "disk.aa-00.writes.delta",
        "operator": ">", "warning": 20000.0,
        "critical": 400000.0, "maxCheckAttempts": 3,
        "notificationInterval": 60,
        "scopes": [ "Hatena-Blog" ],
        "excludeScopes": [ "Hatena-Bookmark: db-master" ]
      }
      """

  Scenario: should get empty monitors
    When I send "GET" request to "/api/v0/monitors"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "monitors": []
      }
      """
