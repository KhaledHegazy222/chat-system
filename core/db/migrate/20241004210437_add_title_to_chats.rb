class AddTitleToChats < ActiveRecord::Migration[7.2]
  def change
    add_column :chats, :title, :string
  end
end
