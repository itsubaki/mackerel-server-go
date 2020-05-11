Feature:
  In order to notify to channels
  As an API user
  I need to be able to request channels

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register channel
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "type": "email",
        "name": "My Channel",
        "emails": [
          "myaddress@example.com"
        ],
        "userIds": [
          "userId"
        ],
        "events": [
          "alert"
        ]
      }
      """
    When I send "POST" request to "/api/v0/channels"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "@string@",
        "type": "email",
        "name": "My Channel",
        "emails": [
          "myaddress@example.com"
        ],
        "userIds": [
          "userId"
        ],
        "events": [
          "alert"
        ]
      }
      """
    Then I keep the JSON response at "id" as "$CHANNEL_ID"

  Scenario: should get channel
    When I send "GET" request to "/api/v0/channels"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "channels": [
          {
            "id":"$CHANNEL_ID",
            "name": "My Channel",
            "type": "email",
            "events": [
              "alert"
            ],
            "emails": [
              "myaddress@example.com"
            ],
            "userIds": [
              "userId"
            ]
          }
        ]
      }
      """

  Scenario: should delete channel
    When I send "DELETE" request to "/api/v0/channels/$CHANNEL_ID"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "$CHANNEL_ID",
        "type": "email",
        "name": "My Channel",
        "emails": [
          "myaddress@example.com"
        ],
        "userIds": [
          "userId"
        ],
        "events": [
          "alert"
        ]
      }
      """

  Scenario: should get empty channels
    When I send "GET" request to "/api/v0/channels"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "channels": []
      }
      """