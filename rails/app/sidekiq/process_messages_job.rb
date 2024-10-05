class ProcessMessagesJob
  include Sidekiq::Worker

  # This Worker Batches Create at most 10K New Messages per execution
  BATCH_SIZE = 10000
  def perform()
    
    
    puts "Proecessing Messages......"  
    created_messages = []
    application_tokens = []
    chat_numbers = []
    BATCH_SIZE.times do
      data = REDIS.lpop('messages_queue')
      if data
        message = JSON.parse(data)
        created_messages << message
        application_tokens << message['application_token']
        chat_numbers << message['chat_number']
      else
        break
      end
    end

    # Here we're selecting all chats that in the selected applications and with any number in the selected numbers
    @chats = Chat.joins(:application)
                .where(application: { token: application_tokens.uniq })
                .where(number: chat_numbers.uniq)


    chats_map = @chats.each_with_object({}) do |chat, hash|
       hash[[chat.application.token, chat.number]] = chat.id
    end

    bulk_insert_data = []

    created_messages.each do |message_date|
      token = message_date['application_token']
      chat_number = message_date['chat_number']
      number = message_date['number']
      content = message_date['content']
      chat_id = chats_map[[token,chat_number]]
      
      if chat_id
        bulk_insert_data << {
          chat_id: chat_id,
          number: number,
          content: content
        }

      else
        puts "No Chat found for (token,number): (#{token},#{chat_number})"
      end
    end

    return if created_messages.empty?

    Message.insert_all(bulk_insert_data)
    puts "Batch inserted #{bulk_insert_data.size} messages"
    
  end
end
