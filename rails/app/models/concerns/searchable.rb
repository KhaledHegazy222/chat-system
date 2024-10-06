module Searchable
  extend ActiveSupport::Concern

  included do
    include Elasticsearch::Model
    include Elasticsearch::Model::Callbacks

    mapping do
      indexes :content, type: :text
    end
    def self.search(query)
      params = {
        query: {
          match: {
              content: query
          }
        }
      }

      self.__elasticsearch__.search(params).records.to_a
    end
  end
end