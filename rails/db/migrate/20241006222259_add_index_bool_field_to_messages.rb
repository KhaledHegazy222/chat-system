class AddIndexBoolFieldToMessages < ActiveRecord::Migration[7.2]
  def change
    add_column :messages, :es_indexed, :boolean, default: false, null: false
    
    # Add an index on the es_indexed column
    add_index :messages, :es_indexed
  end
end
