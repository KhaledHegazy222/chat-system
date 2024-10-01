# config/initializers/redis_listener.rb
Thread.new do
  redis = Redis.new

  loop do
    data = redis.brpop("chats_queue","messages_queue", timeout: 5) # Blocking pop operation
    if data
      queue_name = data[0]  
      payload = data[1]     
      
      # Dispatch data to different workers based on the queue name
      case queue_name
      when "chats_queue"
        ProcessChatsJob.perform_async(payload) 
      when "messages_queue"
        ProcessMessagesJob.perform_async(payload) 
      end 
    end
  end
end