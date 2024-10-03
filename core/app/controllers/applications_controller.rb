require 'securerandom'

class ApplicationsController < ApplicationController
  skip_before_action :verify_authenticity_token

  def index
    @applications = Application.all
    render json: @applications.as_json(only: [:token, :name])
  end
  def show
    @application = Application.where(token: params[:id])
    render json: @application.as_json(only: [:token, :name])
  end

  def create
    name = app_params['name']
    uuid = SecureRandom.uuid
    ProcessApplicationsJob.perform_async({"name" => name,"token" => uuid})
    render json: {name: name,uuid:uuid}

  end

  private
    # Only allow a list of trusted parameters through.
    def app_params
      params.require(:application).permit(:name)
    end
  

end
