class HomesController < ApplicationController
  def home
    @greeting = "Hello, #{Rails.env} Rails!"
  end
end
