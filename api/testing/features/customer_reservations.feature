# file: version.feature
Feature: make a reservation

  Background: 
    Given the server address is "http://localhost:9999/api"

  @essential
  Scenario: allow POST to create a reservation
    Given the request body is "createReservation"
    When I send "POST" request to "/reservations/"
    Then the response code should be 201
    And the response body should match "createReservation"