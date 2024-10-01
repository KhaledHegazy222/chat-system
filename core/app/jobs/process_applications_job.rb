class ProcessApplicationsJob
  include Sidekiq::Worker

  def perform(*args)
    name = args[0]
    token = args[1]
    @application =  Application.new(name:name,token:token)
    @application.save
    puts "Application Created : (#{name}, #{token})"
  end
end
