# This file should ensure the existence of records required to run the application in every environment (production,
# development, test). The code here should be idempotent so that it can be executed at any point in every environment.
# The data can then be loaded with the bin/rails db:seed command (or created alongside the database with db:setup).
#
# Example:
#
#   ["Action", "Comedy", "Drama", "Horror"].each do |genre_name|
#     MovieGenre.find_or_create_by!(name: genre_name)
#   end

# Seed Applications

APPLICATIONS_NUMBER = 10
MIN_CHATS_NUMBER = 20
MAX_CHATS_NUMBER = 30
MIN_MESSAGES_NUMBER = 50
MAX_MESSAGES_NUMBER = 100

APPLICATIONS_NUMBER.times do |i|
  chats_count = rand(MIN_CHATS_NUMBER..MAX_CHATS_NUMBER)
  messages_count = rand(MIN_MESSAGES_NUMBER..MAX_MESSAGES_NUMBER)
  # Create the Application record
  app = Application.create!(
    token: SecureRandom.hex(10),
    name: "Application_#{i + 1}",
    chats_count: chats_count
  )

  # Prepare bulk inserts for chats
  chats = []
  chats_count.times do |j|
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

puts "Seeded Applications, Chats, and Messages successfully!"