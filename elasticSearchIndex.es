{
  "settings": {
    "number_of_shards": 3,
    "number_of_replicas": 2
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "first_name": {
        "type": "text"
      },
      "last_name": {
        "type": "text"
      },
      "is_verified": {
        "type": "boolean"
      },
      "date_of_birth": {
        "type": "date"
      },
      "gender": {
        "type": "keyword"
      },
      "sexual_orientation": {
        "type": "keyword"
      },
      "location": {
        "type": "geo_point"
      },
      "education": {
        "type": "nested",
        "properties": {
          "college": {
            "type": "text"
          },
          "education_level": {
            "type": "keyword"
          }
        }
      },
      "occupation": {
        "type": "keyword"
      },
      "marital_status": {
        "type": "keyword"
      },
      "religion": {
        "type": "keyword"
      },
      "height": {
        "type": "integer"
      },
      "weight": {
        "type": "integer"
      },
      "looking_for": {
        "type": "keyword"
      }
    }
  }
}