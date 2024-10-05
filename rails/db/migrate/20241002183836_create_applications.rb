class CreateApplications < ActiveRecord::Migration[7.2]
  def change
    create_table :applications do |t|
      t.string :token, null: false
      t.string :name, null: false

      t.timestamps
    end
  end
end
