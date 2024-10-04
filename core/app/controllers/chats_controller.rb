class ChatsController < ApplicationController
  
  def index
    # Limit only the last 10 chats (later might add pagination)
    @chats = Chat.includes(:application).limit(10)
    render json: @chats.as_json(only: [:number],include: {application: {only: [:token,:name]}})
  end

  def show
    @chat = Chat
    .includes(:application)
    .find_by(
      number:params[:chat_number],
      application: {token: params[:application_token]}
    )


    if @chat
        render json: @chat.as_json(
          only: [:number],
          include: [
            {application: {only: [:token,:name]}},
          ])
    elsif 
      render json: { error: 'Chat not found' }, status: :not_found
    end
  
  end
end
