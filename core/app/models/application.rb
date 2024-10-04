class Application < ApplicationRecord
  has_many :chats,  dependent: :delete_all
end
