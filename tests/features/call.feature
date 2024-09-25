Feature: Call

  @long
  Scenario: Create call
    When I make a request to create a call
      | audio_url | https://github.com/ggerganov/whisper.cpp/raw/refs/heads/master/samples/jfk.wav |
    Then I should receive call created success response
    And  I wait till the call is processed using long poll
    And  get call should return success response

  Scenario: Create call with invalid audio document
    When I make a request to create a call
      | audio_url | https://example.com/not-audio-file.txt |
    Then I should receive call created success response
    And  I wait till the call is processed
    And  get call should return unprocessable entity

  Scenario: Get call
    When I make a request to get sample processed call
    And  I should receive call processed response

  Scenario: Get non-existing call
    When I make a request to get non-existing call id
    Then I should receive not found error
