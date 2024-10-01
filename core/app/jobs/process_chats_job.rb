class ProcessChatsJob
  include Sidekiq::Worker

  def perform(data)
    puts "New Chat #{data}"
  end
end
