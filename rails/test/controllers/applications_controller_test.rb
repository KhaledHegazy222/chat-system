require 'test_helper'

class ApplicationsControllerTest < ActionDispatch::IntegrationTest
  # Setup a couple of applications for testing
  setup do
    Message.delete_all
    Chat.delete_all
    Application.delete_all
    @application = Application.create!(name: 'App 1', token: SecureRandom.hex(10))
    @application2 = Application.create!(name: 'App 2', token: SecureRandom.hex(10))
  end

  # Test the index action
  test "should get index" do
    get applications_url
    assert_response :success
    json_response = JSON.parse(@response.body)
    
    # Test if response contains both applications
    assert_equal 2, json_response.length
    assert_equal @application.name, json_response[0]["name"]
    assert_equal @application.token, json_response[0]["token"]
  end

  # Test the show action for a valid token
  test "should show application" do
    get "/applications/#{@application.token}"
    assert_response :success
    json_response = JSON.parse(@response.body)

    # Test if correct application is returned
    assert_equal @application.name, json_response["name"]
    assert_equal @application.token, json_response["token"]
  end
end
