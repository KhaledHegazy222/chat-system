class UpdateMessageJob
  include Sidekiq::Job

  def perform(data)
    id = data["id"]
    content = data["content"]
    message = Message.find(id)  # Find the message by its ID
    message.update(content: content)  # Update the content value
    
  end
end
