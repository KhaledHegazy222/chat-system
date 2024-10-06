class ProcessUnindexedMessagesJob
  include Sidekiq::Worker

  def perform()
    
    puts "Proecessing Unindexed Messages......"  
    ActiveRecord::Base.transaction do
        @messages = Message.all.where(es_indexed: false)
        @messages.each do |message|
          # Indexing Each Document Individually
          # TODO: Bulk Index for all of these messages
          message.__elasticsearch__.index_document
        end
    
        # Update es_indexed to true for all processed messages
        Message.where(es_indexed: false).update_all(es_indexed: true)
    end
    puts "Completed indexing messages!"
  end
end
