class CreateApplicationJob
  include Sidekiq::Job

  def perform(data)
    name = data['name']
    token = data['token']
    @application =  Application.new(name:name,token:token)
    @application.save
    
    redis = Redis.new(host:"redis",port: 6379,db: 0)

    # Add Application to Redis Applications Hashset (number of chats = 0)
    application_name_in_hashset = "app##{token}"
    redis.hset('applications_chats_count', application_name_in_hashset , 0)
  end
end