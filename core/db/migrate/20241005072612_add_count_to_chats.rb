class AddCountToChats < ActiveRecord::Migration[7.2]
  def change
    add_column :chats, :messages_count, :integer, null: false, default: 0
  end
end
