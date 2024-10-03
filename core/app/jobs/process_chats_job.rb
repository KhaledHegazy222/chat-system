class ProcessChatsJob
  include Sidekiq::Worker

  def perform(data)
    
    token = data['data']['application_token']
    number = data['data']['number']
    puts "#{token} ------- #{number}"
    @applications = Application.where(token: token)
    puts "#{@applications.count}"
    if @applications.count == 0
      puts "No Application found"
      return
    end
    puts "#{@applications[0]}"
    puts "#{@applications[0].id}"

    @application = @applications[0]
    @chat = @application.chats.create(number:number)
    puts "New Chat Created (#{@application.token},#{@chat.number})"
  end
end
