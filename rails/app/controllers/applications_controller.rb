class ApplicationsController < ApplicationController
  skip_before_action :verify_authenticity_token

  def index
    # Limit only the last 10 apps (later might add pagination)
    @applications = Application.take(10)
    render json: @applications.as_json(only: [:token, :name])
  end
  
  def show
    @application = Application.find_by(token: params[:token])
    if @application
      render json: @application.as_json(only: [:token, :name])
    else
      render json: { error: 'Application not found' }, status: :not_found
    end
  end

  def create
    name = app_params['name']
    
    # Validate presence of the name attribute
    if name.blank?
      return render json: { error: 'Application name is required' }, status: :bad_request
    end

    token = SecureRandom.hex(10)
    CreateApplicationJob.perform_async({"name" => name, "token" => token})
    render json: { name: name, token: token }, status: :created
  end

  # Strong parameters to allow only certain attributes
  private

    def app_params
      params.require(:application).permit(:name)
    end
end
