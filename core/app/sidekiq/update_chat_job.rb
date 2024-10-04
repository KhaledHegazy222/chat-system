class UpdateChatJob
  include Sidekiq::Job

  def perform(data)
    id = data['id']
    title = data['title']

    chat = Chat.find(id)  
    chat.update(title: title)  
  end
end
