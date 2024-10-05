class AddTitleToChats < ActiveRecord::Migration[7.2]
  def change
    add_column :chats, :title, :string, null: false
  end
end
