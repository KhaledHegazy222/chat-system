class MessagesController < ApplicationController
  skip_before_action :verify_authenticity_token

  def index
    @messages = Message.joins(chat: :application).all
    
    if params[:q].present?
      @messages = Message.search(params[:q])
    end

    @messages = @messages.take(10)
   
    render json: @messages.as_json(
      only: [:number,:content],
      # include: {
      #   chat: {
      #     only: [:number],
      #     include: {
      #       application: {
      #         only: [:name, :token]
      #       }
      #     }
      #   }
      # }
    )
  end

  def show
    if params[:application_token].blank? || params[:chat_number].blank? || params[:messages_number].blank?
      return render json: { error: 'Missing required parameters' }, status: :bad_request
    end

    @message = Message.joins(chat: :application)
                      .where(applications: { token: params[:application_token] })
                      .where(chats: { number: params[:chat_number] })
                      .where(number: params[:messages_number])
                      .first

    if @message
      render json: @message.as_json(
        only: [:number,:content],
        # include: {
        #   chat: {
        #     only: [:number],
        #     include: {
        #       application: {
        #         only: [:name, :token]
        #       }
        #     }
        #   }
        # }
      )
    else
      render json: { error: 'Message not found' }, status: :not_found
    end
  end

  def search
    if params[:application_token].blank? || params[:chat_number].blank?
      return render json: { error: 'Missing required parameters' }, status: :bad_request
    end

    @messages = Message.joins(chat: :application).all
                       .where(applications: { token: params[:application_token] })
                       .where(chats: { number: params[:chat_number] })  
                  
    if params[:q].present?
      @messages = @messages.search(params[:q])
    end

    @messages = @messages.take(10)
    render json: @messages.as_json(
      only: [:number,:content],
      # include: {
      #   chat: {
      #     only: [:number],
      #     include: {
      #       application: {
      #         only: [:name, :token]
      #       }
      #     }
      #   }
      # }
    )  
  end

  def edit
    if params[:application_token].blank? || params[:chat_number].blank? || params[:messages_number].blank? || message_params[:content].blank?
      return render json: { error: 'Missing required parameters' }, status: :bad_request
    end

    @message = Message.joins(chat: :application)
                      .where(applications: { token: params[:application_token] })
                      .where(chats: { number: params[:chat_number] })
                      .where(number: params[:messages_number])
                      .first

    if not @message
      return render json: { error: 'Message not found' }, status: :not_found
    end

    UpdateMessageJob.perform_sync({"id" => @message.id, "content" => message_params['content']})
    render json: { status: "success" }
  end

  # Strong parameters to allow only certain attributes
  private

  def message_params
    params.require(:message).permit(:content)
  end
end
