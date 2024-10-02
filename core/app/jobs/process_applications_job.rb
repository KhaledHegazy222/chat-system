class ProcessApplicationsJob
  include Sidekiq::Worker

  def perform(data)
    name = data['name']
    token = data['token']
    @application =  Application.new(name:name,token:token)
    @application.save
    puts "Application Created : (#{name}, #{token})"
  end
end
