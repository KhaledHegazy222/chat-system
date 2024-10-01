class ProcessChatsJob
  include Sidekiq::Worker

  def perform(json_data)
    data = JSON.parse(json_data)
    token = data['token']
    number = data['number']
    application = Application.where(token: token)[0]
    if application
        chat = application.chats.create(number:number)
        puts "New Chat Created (#{application.token},#{chat.number})"
    else
        puts "No Application found"
    end
  end
end
