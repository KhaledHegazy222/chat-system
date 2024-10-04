class ProcessMessagesJob
  include Sidekiq::Worker
  BATCH_SIZE = 3
  def perform()
    redis = Redis.new(host:"redis",port: 6379,db: 0)
    
    puts "Proecessing Messages......"  
    created_messages = []
    updated_messages = []
    BATCH_SIZE.times do
      # This Queue contains newly created messages
      data = redis.blpop('messages_queue',timeout: 1)
      if data
        message = JSON.parse(data[1])
        puts "Message: #{message}"
        puts "Message Operation: #{message['operation']}"
        if message['operation'] == 'create'
            created_messages << message['data']
        elsif message['operation'] == 'update'
            updated_messages << message['data']
        end
      end
    end

    return if created_messages.empty? and updated_messages.empty?

    created_messages.each do |message_date|
      # TODO: this code is creating only one instance, modify it to bulk insert the created messages here
      number = message_date['number']
      chat_number = message_date['chat_number']
      application_token = message_date['application_token']
      content = message_date['content']
      
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

    # TODO: Bulk Updates....
    
  end
end
