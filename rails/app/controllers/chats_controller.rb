class ChatsController < ApplicationController
  skip_before_action :verify_authenticity_token  

  def index
    # Limit only the last 10 chats (later might add pagination)
    @chats = Chat.includes(:application).take(10)
    render json: @chats.as_json(
      only: [:number,:title],
      include: {application: {only: [:token,:name]}}
    )
  end

  def show
    if params[:application_token].blank? || params[:chat_number].blank?
      return render json: { error: 'Missing required parameters' }, status: :bad_request
    end

    @chat = Chat
            .includes(:application)
            .find_by(
              number: params[:chat_number],
              application: { token: params[:application_token] }
            )

    if @chat
      render json: @chat.as_json(
        only: [:number,:title],
        include: [
          {application: {only: [:token,:name]}},
        ]
        )
    else
      render json: { error: 'Chat not found' }, status: :not_found
    end
  end

  def edit
    if params[:application_token].blank? || params[:chat_number].blank? || chat_params[:title].blank?
      return render json: { error: 'Missing required parameters' }, status: :bad_request
    end

    @chat = Chat
            .joins(:application)
            .find_by(
              number: params[:chat_number],
              application: { token: params[:application_token] }
            )

    if not @chat
      return render json: { error: 'Chat not found' }, status: :not_found
    end
    
    UpdateChatJob.perform_async({ "id" => @chat.id, "title" => chat_params['title'] })
    render json: { status: "success" }
  end
  
  private 

  def chat_params
    params.require(:chat).permit(:title)
  end
end
