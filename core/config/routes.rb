
Rails.application.routes.draw do
  # Define your application routes per the DSL in https://guides.rubyonrails.org/routing.html
  require 'sidekiq/web'
  
  mount Sidekiq::Web => '/sidekiq'

  get '/applications', to: "applications#index"
  get '/applications/:token', to: "applications#show"
  post '/applications', to: "applications#create"

  get '/chats', to: "chats#index"
  get '/applications/:application_token/chats/:chat_number', to: "chats#show"
  patch '/applications/:application_token/chats/:chat_number', to: "chats#edit"

  get '/messages', to: "messages#index"
  get '/applications/:application_token/chats/:chat_number/messages', to: "messages#search"
  get '/applications/:application_token/chats/:chat_number/messages/:messages_number', to: "messages#show"
  patch '/applications/:application_token/chats/:chat_number/messages/:messages_number', to: "messages#edit"
  
  
  

end
