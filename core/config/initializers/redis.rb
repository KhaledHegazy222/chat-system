require 'redis'

REDIS = Redis.new(url: ENV.fetch("REDIS_URL"))
