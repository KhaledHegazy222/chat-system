# config/initializers/redis_listener.rb
Thread.new do
  redis = Redis.new(host:"redis",port: 6379,db: 0)

  loop do
    data = redis.brpop("chats_queue","messages_queue", timeout: 5) # Blocking pop operation
    if data
      queue_name = data[0]  
      payload = data[1]     

      # Assuming that payload is always json object pased in string format
      data = JSON.parse(payload)
      
      # Dispatch data to different workers based on the queue name
      case queue_name
      when "chats_queue"
        ProcessChatsJob.perform_async(data) 
      when "messages_queue"
        ProcessMessagesJob.perform_async(data) 
      end 
    end
  end
end