class ProcessApplicationsJob
  include Sidekiq::Worker

  def perform(data)
    name = data['name']
    token = data['token']
    @application =  Application.new(name:name,token:token)
    @application.save
    puts "Application Created : (#{name}, #{token})"
    
    redis = Redis.new(host:"redis",port: 6379,db: 0)
    redis.hset('applications_chats_count', "app##{token}" , 0)
    
  end
end
