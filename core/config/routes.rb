
Rails.application.routes.draw do
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html
  require 'sidekiq/web'
  
  mount Sidekiq::Web => '/sidekiq'

  get '/applications', to: "applications#index"
  get '/applications/:token', to: "applications#show"
  post '/applications', to: "applications#create"

  get '/chats', to: "chats#index"
  get '/chats/:application_token/:chat_number', to: "chats#show"

  get '/messages', to: "messages#index"
  
  

end
