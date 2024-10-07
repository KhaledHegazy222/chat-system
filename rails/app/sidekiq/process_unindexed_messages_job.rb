class ProcessUnindexedMessagesJob
  include Sidekiq::Worker

  def perform()
    
    puts "Proecessing Unindexed Messages......"  
    ActiveRecord::Base.transaction do
        @messages = Message.all.where(es_indexed: false)
        # Bulk Index for all of these messages
        bulk_body = @messages.map do |message|
          [
            { index: { _id: message.id, _index: Message.index_name } }, 
            message.as_json
          ]
        end.flatten
        Message.__elasticsearch__.client.bulk({
          index: Message.index_name,
          body: bulk_body
        })
    
        # Update es_indexed to true for all processed messages
        Message.where(es_indexed: false).update_all(es_indexed: true)
    end
    puts "Completed indexing messages!"
  end
end
