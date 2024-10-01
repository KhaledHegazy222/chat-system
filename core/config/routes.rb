
Rails.application.routes.draw do
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html
  require 'sidekiq/web'
  
  mount Sidekiq::Web => '/sidekiq'

  resources :applications
  

end
