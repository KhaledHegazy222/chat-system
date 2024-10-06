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
    @message = Message.joins(chat: :application)
                      .where(applications: { token: params[:application_token] })
                      .where(chats: { number: params[:chat_number] })
                      .where(number: params[:messages_number])
                      .first

    if not @message
      render json: { error: 'Message not found' }, status: :not_found
      return
    end

    UpdateMessageJob.perform_sync({"id" => @message.id, "content" => message_params['content']})
    render json: {status: "success"}
  
  end

  # Strong parameters to allow only certain attributes
  private
  def message_params
    params.require(:message).permit(:content)
  end

end
