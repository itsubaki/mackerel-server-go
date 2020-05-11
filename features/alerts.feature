Feature:
  In order to get alerts
  As an API user
  I need to be able to request alerts

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should get alerts
    When I send "GET" request to "/api/v0/alerts"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "alerts": []
      }
      """

  Scenario: should close alert
    Given the following alerts exist:
      | org_id      | id      | status   | monitor_id | type |
      | 4b825dc642c | 1234568 | CRITICAL | 1234567890 | host |
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "reason": "closed manually"
      }
      """
    When I send "POST" request to "/api/v0/alerts/1234568/close"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "id": "1234568",
        "status": "OK",
        "monitorId": "1234567890",
        "type": "host",
        "reason": "closed manually",
        "openedAt": "@number@",
        "closedAt": "@number@"
      }
      """