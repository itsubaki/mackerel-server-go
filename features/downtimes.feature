Feature:
  In order to notify with downtime
  As an API user
  I need to be able to request downtimes

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register downtime
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "example",
        "memo": "my first downtime",
        "start": 1351700100,
        "duration": 60,
        "recurrence": {
          "type": "weekly",
          "interval": 60,
          "weekdays": [ "Sunday", "Monday" ],
          "until": 1351700100
        },
        "serviceScopes": [ "production" ],
        "serviceExcludeScopes": [ "develop" ],
        "roleScopes": [ "db" ],
        "roleExcludeScopes": [ "application" ],
        "monitorScopes": [ "12345678901" ],
        "monitorExcludeScopes": [ "12345678902" ]
      }
      """
    When I send "POST" request to "/api/v0/downtimes"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "name": "example",
        "memo": "my first downtime",
        "start": 1351700100,
        "duration": 60,
        "recurrence": {
          "type": "weekly",
          "interval": 60,
          "weekdays": [ "Sunday", "Monday" ],
          "until": 1351700100
        },
        "serviceScopes": [ "Hatena-Bookmark" ],
        "serviceExcludeScopes": [ "Hatena-Blog" ],
        "roleScopes": [ "Hatena-Bookmark:db-master" ],
        "roleExcludeScopes": [ "Hatena-Blog:db-master" ],
        "monitorScopes": [ "12345678901" ],
        "monitorExcludeScopes": [ "12345678902" ]
      }
      """
    Then I keep the JSON response at "id" as "$DOWNTIME_ID"

  Scenario: should get downtime
    When I send "GET" request to "/api/v0/downtimes"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "downtimes": [
          {
            "id": "@string@",
            "name": "example",
            "memo": "my first downtime",
            "start": 1351700100,
            "duration": 60,
            "recurrence": {
              "type": "weekly",
              "interval": 60,
              "weekdays": [ "Sunday", "Monday" ],
              "until": 1351700100
            },
            "serviceScopes": [ "Hatena-Bookmark" ],
            "serviceExcludeScopes": [ "Hatena-Blog" ],
            "roleScopes": [ "Hatena-Bookmark:db-master" ],
            "roleExcludeScopes": [ "Hatena-Blog:db-master" ],
            "monitorScopes": [ "12345678901" ],
            "monitorExcludeScopes": [ "12345678902" ]
          }
        ]
      }
      """

  Scenario: should update downtime
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "example",
        "memo": "my first downtime",
        "start": 1351700100,
        "duration": 60,
      }
      """
    When I send "PUT" request to "/api/v0/downtimes/$DOWNTIME_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$DOWNTIME_ID"
      }
      """

  Scenario: should get updated downtime
    When I send "GET" request to "/api/v0/downtimes"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "downtimes": [
          {
            "id": "@string@",
            "name": "example",
            "memo": "my first downtime",
            "start": 1351700100,
            "duration": 60
          }
        ]
      }
      """

  Scenario: should delete downtime
    When I send "DELETE" request to "/api/v0/downtimes/$DOWNTIME_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "name": "example",
        "memo": "my first downtime",
        "start": 1351700100,
        "duration": 60
      }
      """

  Scenario: should get empty downtime
    When I send "GET" request to "/api/v0/downtimes"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "downtimes": []
      }
      """