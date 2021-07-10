Feature:
  In order to know service status
  As an API user
  I need to be able to request status

  Scenario: should get status
    When I send "GET" request to "/status"
    Then the response code should be 200

  Scenario: should get status at root
    When I send "GET" request to "/"
    Then the response code should be 200
