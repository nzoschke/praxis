class CreateUsersTable < ActiveRecord::Migration
  def change
    create_table :users_tables do |t|
      t.uuid :id
      t.text :email
    end
    add_index :users_tables, :id, unique: true
    add_index :users_tables, :email, unique: true
  end
end
