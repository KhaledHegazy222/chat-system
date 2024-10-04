class ProcessChatsJob
  include Sidekiq::Worker
  BATCH_SIZE = 3
  def perform()
    redis = Redis.new(host:"redis",port: 6379,db: 0)
    
    puts "Proecessing Chats......"  
    created_chats = []
    BATCH_SIZE.times do
      # This Queue contains newly created chats
      data = redis.blpop('chats_queue',timeout: 1)
      if data
        chat = JSON.parse(data[1])
        created_chats << chat
      end
    end
    
    return if created_chats.empty?

    puts "New Chats: #{created_chats}"

    created_chats.each do |chat_data|
      # TODO: this code is creating only one instance, modify it to bulk insert the created chats here
      token = chat_data['data']['application_token']
      number = chat_data['data']['number']

      @applications = Application.where(token: token)
      if @applications.count == 0
        puts "No Application found"
        return
      end
  
      @application = @applications[0]
      @chat = @application.chats.create(number:number)
      puts "New Chat Created (#{@application.token},#{@chat.number})"


      # Add Chat to Redis Chats Hashset (number of messsages = 0)
      chat_name_in_hashset = "chat##{token}-#{number}"
      redis.hset("chats_messages_count",chat_name_in_hash,0)
      
          
    end

    # TODO: Bulk Updates....
    
  end
end
