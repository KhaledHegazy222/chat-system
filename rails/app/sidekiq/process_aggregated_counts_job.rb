class ProcessAggregatedCountsJob
  include Sidekiq::Job


  def perform

    puts "Process Aggregations..."
    
    # Locking the key to make sure only one process is executing this job 
    # (might run on multiple rails servers)
    lock_key = 'aggregated_count_lock'
    if acquire_lock(lock_key)
      begin
        # Wrap It in transaction
        ActiveRecord::Base.transaction do 
          
          # Update Applications Count Fields
          ActiveRecord::Base.connection.execute(
            "update applications a join ( 
              select a.id as id, count(a.id) as cnt 
              from applications a join chats c 
              on a.id = c.application_id 
              group by a.id
            ) aggregated_counts
            on a.id = aggregated_counts.id 
            set a.chats_count = aggregated_counts.cnt;")

          # Update Chats Count Fields
          ActiveRecord::Base.connection.execute(
            "update chats c join ( 
              select c.id as id, count(c.id) as cnt 
              from chats c join messages m 
              on c.id = m.chat_id 
              group by c.id 
            ) aggregated_counts
            on c.id = aggregated_counts.id 
            set c.messages_count = aggregated_counts.cnt;")

        end
      ensure
        release_lock(lock_key)
      end
    else
      puts "Failed To Aquire The lock, skipping the job..."
    end
  end


  private
    

    def acquire_lock(lock_key, timeout: 10, interval: 1)
      start_time = Time.now

      loop do
        expiry_seconds = 1800 # 30 minutes = 30 * 60 seconds
        acquired = REDIS.set(lock_key, 'locked', nx: true, ex: expiry_seconds)

        return true if acquired

        if Time.now - start_time > timeout
          return false
        end

        sleep(interval)
      end
    end

    def release_lock(lock_key)
      REDIS.del(lock_key)
    end
end
