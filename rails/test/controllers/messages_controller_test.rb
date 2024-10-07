require 'test_helper'
require 'sidekiq/testing'

class MessagesControllerTest < ActionDispatch::IntegrationTest


  setup do
    Message.delete_all
    Chat.delete_all
    Application.delete_all
    @application = Application.create!(name: 'App 1', token: SecureRandom.hex(10))
    @chat = Chat.create!(number: 1, title: 'Chat 1', application: @application)
    @message1 = Message.create!(number: 1, content: 'Hello, World!', chat: @chat)
    @message2 = Message.create!(number: 2, content: 'Goodbye, World!', chat: @chat)
  end

  # Test the index action
  test "should get index" do
    get "/applications/#{ @application.token}/chats/#{ @chat.number}/messages"

    assert_response :success
    json_response = JSON.parse(@response.body)

    assert_equal 2, json_response.length
    assert_equal @message1.content, json_response[0]["content"]
    assert_equal @message2.content, json_response[1]["content"]
  end

  # Test the show action for a valid message
  test "should show message" do
    get "/applications/#{ @application.token}/chats/#{ @chat.number}/messages/#{@message1.number}"

    assert_response :success
    json_response = JSON.parse(@response.body)

    assert_equal @message1.content, json_response["content"]
    assert_equal @message1.number, json_response["number"]
  end

  # Test the show action for a non-existent message
  test "should not show non-existent message" do
    get "/applications/#{ @application.token}/chats/#{ @chat.number}/messages/#{999}"
    assert_response :not_found
    json_response = JSON.parse(@response.body)

    assert_equal 'Message not found', json_response["error"]
  end

end
