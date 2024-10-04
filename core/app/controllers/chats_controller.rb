class ChatsController < ApplicationController
  skip_before_action :verify_authenticity_token  
  def index
    # Limit only the last 10 chats (later might add pagination)
    @chats = Chat.includes(:application).take(10)
    render json: @chats.as_json(only: [:number,:title],include: {application: {only: [:token,:name]}})
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
          only: [:number,:title],
          include: [
            {application: {only: [:token,:name]}},
          ])
    elsif 
      render json: { error: 'Chat not found' }, status: :not_found
    end
  
  end


  def edit
    @chat = Chat
    .includes(:application)
    .find_by(
      number:params[:chat_number],
      application: {token: params[:application_token]}
    )

    if not @chat 
      render json: { error: 'Chat not found' }, status: :not_found
      return
    end
    
    UpdateChatJob.perform_async({"id" => @chat.id,"title" => chat_params['title']})
    render json: {status: "success"}
    
  end
  
  private 

  def chat_params
    params.require(:chat).permit(:title)
  end
end
