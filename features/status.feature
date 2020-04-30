Feature: get status
  In order to know mackerel-api status
  As an API user
  I need to be able to request status

  Scenario: should get status code
    When I send "GET" request to "/"
    Then the response code should be 200