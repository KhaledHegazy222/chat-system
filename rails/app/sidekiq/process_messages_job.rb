class ProcessMessagesJob
  include Sidekiq::Worker
  
  # This Worker Batches Create at most 10K New Messages per execution
  BATCH_SIZE = 10_000

  def perform
    puts "Processing Messages..."
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

    return if created_messages.empty?

    ActiveRecord::Base.transaction do
      # Here we're selecting all chats that are in the selected applications and with any number in the selected numbers
      @chats = Chat.joins(:application)
                   .where(application: { token: application_tokens.uniq })
                   .where(number: chat_numbers.uniq)

      # Create a hash map for quick lookups: { [application_token, chat_number] => chat_id }
      chats_map = @chats.each_with_object({}) do |chat, hash|
        hash[[chat.application.token, chat.number]] = chat.id
      end

      bulk_insert_data = []

      created_messages.each do |message_data|
        token = message_data['application_token']
        chat_number = message_data['chat_number']
        number = message_data['number']
        content = message_data['content']
        chat_id = chats_map[[token, chat_number]]

        if chat_id
          bulk_insert_data << {
            chat_id: chat_id,
            number: number,
            content: content
          }
        else
          puts "No Chat found for (token, number): (#{token}, #{chat_number})"
        end
      end

      # Bulk insert messages into the database
      Message.insert_all(bulk_insert_data)
      puts "Batch inserted #{bulk_insert_data.size} messages"

      # Trigger the job to index the messages in Elasticsearch
      ProcessUnindexedMessagesJob.perform_sync
    end
  rescue StandardError => e
    puts "Failed to process messages: #{e.message}"
  end
end
