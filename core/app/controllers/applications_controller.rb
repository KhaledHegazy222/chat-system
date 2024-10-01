require 'securerandom'

class ApplicationsController < ApplicationController
  skip_before_action :verify_authenticity_token

  def index
    @applications = Application.all
    render json: @applications
  end
  def show
    @application = Application.find(params[:id])
    render json: @application
  end

  def create
    name = app_params['name']
    uuid = SecureRandom.uuid
    ProcessApplicationsJob.perform_async(name,uuid)
    render json: {status: "success"}

  end

  private
    # Only allow a list of trusted parameters through.
    def app_params
      params.require(:application).permit(:name)
    end
  

end
