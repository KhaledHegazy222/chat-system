namespace :chat do
  desc "Run chat system seeding"
  task seed: :environment do
    
      APPLICATIONS_NUMBER = 10
      MIN_CHATS_NUMBER = 50
      MAX_CHATS_NUMBER = 100
      MIN_MESSAGES_NUMBER = 100
      MAX_MESSAGES_NUMBER = 300

      APPLICATIONS_NUMBER.times do |i|
        chats_count = rand(MIN_CHATS_NUMBER..MAX_CHATS_NUMBER)
        
        # Create the Application record
        app = Application.create!(
          token: SecureRandom.hex(10),
          name: "Application_#{i + 1}",
          chats_count: chats_count
        )

        # Prepare bulk inserts for chats
        chats = []
        chats_count.times do |j|
          messages_count = rand(MIN_MESSAGES_NUMBER..MAX_MESSAGES_NUMBER)
          chats << {
            number: j + 1,
            application_id: app.id,
            title: "Chat_#{j + 1} for #{app.name}",
            messages_count: messages_count,
            created_at: Time.current,
            updated_at: Time.current
          }
        end

        
        # Bulk create Chats
        Chat.insert_all(chats)

        chats = Chat.where(application_id: app.id)
        messages = []
        
        # Now we need to seed messages for each chat
        chats.each do |chat|
          chat_id = chat.id
          messages_count = chat.messages_count

          messages_count.times do |k|
            messages << {
              number: k + 1,
              chat_id: chat_id,
              content: "Message #{k + 1} in Chat_#{chat[:number]} for #{app.name}",
              created_at: Time.current,
              updated_at: Time.current
            }
          end
        end
        # Bulk create Messages
        Message.insert_all(messages)

        puts "Created #{chats_count} chats for #{app.name}."
      end

      # Populating Redis Server Data
      Application.all.each do |app|
        application_name_in_hashset = "app##{app.token}"
        REDIS.hset('applications_chats_count', application_name_in_hashset , app.chats_count)
      end

      Chat.all.includes(:application).each do |chat|
        chat_name_in_hashset = "chat##{chat.application.token}-#{chat.number}"
        REDIS.hset("chats_messages_count", chat_name_in_hashset, chat.messages_count)
      end

      puts "Seeded Applications, Chats, and Messages successfully!"
  end
end