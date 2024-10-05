class AddIndeicesToApplications < ActiveRecord::Migration[7.2]
  def change
    add_index :applications, :token, unique: true 
  end
end
