require 'sinatra'
require 'haml' # templating engine

set :bind, "0.0.0.0"
port = ENV["PORT"] || "9000"
set :port, port

#this method calls haml template defined in views folder
get '/' do

  #hello world haml template
  haml :helloWorld

end                                                                                                                            
