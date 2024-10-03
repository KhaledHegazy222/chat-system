class ProcessMessagesJob
  include Sidekiq::Worker

  def perform(data)
    number = data['data']['number']
    chat_number = data['data']['chat_number']
    application_token = data['data']['application_token']
    content = data['data']['content']
    
    @applications = Application.where(token:application_token)
    if @applications.count == 0
      puts "Application Not Found"
      return
    end

    @application = @applications[0]
  
    @chats = @application.chats.where(number:chat_number)
    if @chats.count == 0
      puts "Chat Not Found"
      return
    end

    @chat = @chats[0]
    @message =  @chat.messages.create(number:number,content:content)
    puts "New Message Created #{@message.id} #{@message.number} #{@message.content}"
  end
end
