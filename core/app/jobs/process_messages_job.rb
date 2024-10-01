class ProcessMessagesJob
  include Sidekiq::Worker

  def perform(json_data)
    data = JSON.parse(json_data)
    puts "New Message #{data}"
  end
end
