Feature:
  In order to monitor the check reports
  As an API user
  I need to be able to request monitoring

  Background:
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb"

  Scenario: should register checks report
    Given I set "Content-Type" header with "application/json"
    Given I set request body:
      """
      {
        "reports": [
          {
            "source": {
              "type": "host",
              "hostId": "12345678901"
            },
            "name": "check-tcp_mysql",
            "status": "WARNING",
            "message": "TCP CRITICAL: dial tcp [::1]:3306",
            "occurredAt": 1351700100
          }
         ]
      }
      """
    When I send "POST" request to "/api/v0/monitoring/checks/report"
    Then the response code should be 200
    Then the response should match json:
      """
      {
        "success": true
      }
      """