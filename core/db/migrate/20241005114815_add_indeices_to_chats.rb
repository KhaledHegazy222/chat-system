class AddIndeicesToChats < ActiveRecord::Migration[7.2]
  def change
    add_index :chats, [:application_id, :number], unique: true
  end
end
