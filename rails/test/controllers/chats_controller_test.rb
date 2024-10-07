require 'test_helper'

class ChatsControllerTest < ActionDispatch::IntegrationTest
  # Setup applications and chats for testing
  setup do
    Message.delete_all
    Chat.delete_all
    Application.delete_all
    @application = Application.create!(name: 'App 1', token: SecureRandom.hex(10))
    
    @chat1 = Chat.create!(number: 1, title: 'Chat 1', application: @application)
    @chat2 = Chat.create!(number: 2, title: 'Chat 2', application: @application)
  end

  # Test the show action for a valid chat
  test "should show chat" do
    get "/applications/#{@application.token}/chats/#{@chat1.number}"
    assert_response :success
    json_response = JSON.parse(@response.body)

    # Test if correct chat is returned
    assert_equal @chat1.title, json_response["title"]
    assert_equal @chat1.number, json_response["number"]
  end

  # Test the show action for an invalid chat
  test "should not show non-existing chat" do
    get "/applications/#{@application.token}/chats/999"
    assert_response :not_found
    json_response = JSON.parse(@response.body)

    # Test if error message is returned
    assert_equal 'Chat not found', json_response["error"]
  end

  # Test the edit action for a non-existing chat
  test "should not edit non-existing chat" do
    patch "/applications/#{@application.token}/chats/999", params: { chat: { title: 'Updated Chat' } }
    assert_response :not_found
    json_response = JSON.parse(@response.body)

    # Test if error message is returned
    assert_equal 'Chat not found', json_response["error"]
  end
end
