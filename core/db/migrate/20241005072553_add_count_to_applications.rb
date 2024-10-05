class AddCountToApplications < ActiveRecord::Migration[7.2]
  def change
    add_column :applications, :chats_count, :integer, null: false, default: 0
  end
end
