class ProcessChatsJob
  include Sidekiq::Worker

  # This Worker Batches Create at most 10K New Chats per execution
  BATCH_SIZE = 10_000

  def perform
    puts "Processing Chats..."
    created_chats = []
    application_tokens = []

    BATCH_SIZE.times do
      # This Queue contains newly created chats
      data = REDIS.lpop('chats_queue')
      if data
        chat = JSON.parse(data)
        created_chats << chat
        application_tokens << chat['application_token']
      end
    end

    return if created_chats.empty?

    ActiveRecord::Base.transaction do
      @applications = Application.where(token: application_tokens.uniq)

      # Create a hash map for quick lookups: { token => application_id }
      application_map = @applications.each_with_object({}) do |app, hash|
        hash[app.token] = app.id
      end

      bulk_insert_data = []
      redis_insert_data = []

      created_chats.each do |chat_data|
        token = chat_data['application_token']
        number = chat_data['number']
        title = chat_data['title']
        application_id = application_map[token]

        if application_id
          bulk_insert_data << {
            application_id: application_id,
            number: number,
            title: title,
          }

          # Prepare data for Redis insertions
          chat_name_in_hashset = "chat##{token}-#{number}"
          redis_insert_data << {
            key: chat_name_in_hashset,
            value: 0  # Number of messages = 0
          }
        else
          puts "No Application found for token: #{token}"
        end
      end

      # Bulk insert all chats
      Chat.insert_all(bulk_insert_data)
      puts "Batch inserted #{bulk_insert_data.size} chats"

      # Add Chat to Redis Chats Hashset
      redis_insert_data.each do |chat_data|
        REDIS.hset("chats_messages_count", chat_data[:key], chat_data[:value])
      end
    end
  rescue StandardError => e
    puts "Failed to process chats: #{e.message}"
  end
end
