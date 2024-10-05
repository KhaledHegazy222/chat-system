class AddIndeicesToMessages < ActiveRecord::Migration[7.2]
  def change
    add_index :messages, [:chat_id, :number], unique: true
  end
end
