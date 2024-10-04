class MessagesController < ApplicationController
  
  def index
    @messages = Message.includes(chat: :application).all
    if params[:q].present?
      @messages = @messages.search(params[:q])
    end
    render json: @messages.as_json(only: [:content],include: {
      chat: {
        only: [:number],
        include: {
          application: {
            only: [:name, :token]
          }
        }
      }
    })
  end
end
