# file: version.feature
Feature: smoke test all HTTP endpoints

  Background:
    Given the server address is "http://localhost:9999/api"

  @essential
  Scenario Outline: "<method>" to "<path>" with "<payload>" returns <statusCode>
    Given the request body is "<payload>"
    When I send "<method>" request to "<path>"
    Then the response code should be <statusCode>
    And the response body should match "<responseBody>"
    Examples:
      | method | path           | payload           | statusCode | responseBody      |
      | POST   | /reservations/ | createReservation | 201        | createReservation |

