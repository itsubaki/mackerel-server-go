Feature:
  In order to notify to groups
  As an API user
  I need to be able to request notification-groups

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register notification group
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "slack",
        "notificationLevel": "all",
        "childNotificationGroupIds": [],
        "childChannelIds":["2vh7AZ21abc"],
        "monitors": [
          {
            "id": "2qtozU21abc",
            "skipDefault": false
          }
        ],
        "services": [
          {
            "name": "Example-Service-1"
          },
          {
            "name": "Example-Service-2"
          }
        ]
      }
      """
    When I send "POST" request to "/api/v0/notification-groups"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "name": "slack",
        "notificationLevel": "all",
        "childNotificationGroupIds": [],
        "childChannelIds":["2vh7AZ21abc"],
        "monitors": [
          {
            "id": "2qtozU21abc",
            "skipDefault": false
          }
        ],
        "services": [
          {
            "name": "Example-Service-1"
          },
          {
            "name": "Example-Service-2"
          }
        ]
      }
      """
    Then I keep the JSON response at "id" as "$GROUP_ID"

  Scenario: should update notification group
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "name": "slack",
        "notificationLevel": "critical",
        "childNotificationGroupIds": [],
        "childChannelIds":["2vh7AZ21abc"],
        "monitors": [
          {
            "id": "2qtozU21abc",
            "skipDefault": false
          }
        ],
        "services": [
          {
            "name": "Example-Service-1"
          },
          {
            "name": "Example-Service-2"
          }
        ]
      }
      """
    When I send "PUT" request to "/api/v0/notification-groups/$GROUP_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$GROUP_ID",
        "name": "slack",
        "notificationLevel": "critical",
        "childNotificationGroupIds": [],
        "childChannelIds":["2vh7AZ21abc"],
        "monitors": [
          {
            "id": "2qtozU21abc",
            "skipDefault": false
          }
        ],
        "services": [
          {
            "name": "Example-Service-1"
          },
          {
            "name": "Example-Service-2"
          }
        ]
      }
      """

  Scenario: should get notification group
    When I send "GET" request to "/api/v0/notification-groups"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "notificationGroups":[
          {
            "id": "$GROUP_ID",
            "name": "slack",
            "notificationLevel": "critical",
            "childNotificationGroupIds": [],
            "childChannelIds": ["2vh7AZ21abc"],
            "monitors": [
              {
                "id": "2qtozU21abc",
                "skipDefault":false
              }
            ],
            "services": [
              {
                "name": "Example-Service-1"
              },
              {
                "name": "Example-Service-2"
              }
            ]
          }
        ]
      }
      """

  Scenario: should delete notification group
    When I send "DELETE" request to "/api/v0/notification-groups/$GROUP_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$GROUP_ID",
        "name": "slack",
        "notificationLevel": "critical",
        "childNotificationGroupIds": [],
        "childChannelIds":["2vh7AZ21abc"],
        "monitors": [
          {
            "id": "2qtozU21abc",
            "skipDefault": false
          }
        ],
        "services": [
          {
            "name": "Example-Service-1"
          },
          {
            "name": "Example-Service-2"
          }
        ]
      }
      """

  Scenario: should get empty notification group
    When I send "GET" request to "/api/v0/notification-groups"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "notificationGroups": []
      }
      """