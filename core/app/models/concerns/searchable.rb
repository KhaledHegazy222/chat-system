module Searchable
  extend ActiveSupport::Concern

  included do
    include Elasticsearch::Model
    include Elasticsearch::Model::Callbacks

    mapping do
      # mapping definition goes here
      indexes :content, type: 'text'
    end
    def self.search(query)
      params = {
        query: {
          multi_match: {
            query: query,
            fields: ['content'],
          }
        }
      }

      self.__elasticsearch__.search(params).records.to_a
    end
  end
end