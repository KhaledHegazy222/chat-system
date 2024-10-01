class ProcessMessagesJob
  include Sidekiq::Worker

  def perform(data)
    puts "New Message #{data}"
  end
end
