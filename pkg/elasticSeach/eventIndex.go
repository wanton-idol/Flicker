package elasticSeach

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/SuperMatch/model"
	"github.com/SuperMatch/model/dto"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/zapLogger"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
)

var eventIndex = "events"

type EventIndexer interface {
	CreateIndex() error
	IndexUserEvent(event elasticsearchPkg.Event, doc []byte) error
	GetEventsById(eventId int) ([]elasticsearchPkg.Event, error)
	SearchEvents(user model.UserProfile, page model.Pagination, filters dto.EventFilterDTO) ([]elasticsearchPkg.Event, error)
	UpdateUserEvent(event elasticsearchPkg.Event, doc []byte) error
	DeleteUserEvent(eventID int) error
}

type EventIndexerImpl struct {
	IndexName string
	esClient  *opensearch.Client
}

func NewEventIndexerImpl() *EventIndexerImpl {
	return &EventIndexerImpl{
		IndexName: eventIndex,
		esClient:  EsClient,
	}
}

func (e *EventIndexerImpl) CreateIndex() error {
	index := e.IndexName
	mapping := `{
  "settings": {
    "number_of_shards": 2,
    "number_of_replicas": 1
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "integer"
      },
      "user_id": {
        "type": "integer"
      },
      "type": {
        "type": "text"
      },
      "event_time": {
        "type": "date"
      },
      "description": {
        "type": "text"
      },
      "attendees": {
        "type": "text"
      },
      "expires_at": {
        "type": "date"
      },
      "address1": {
        "type": "text"
      },
      "address2": {
        "type": "text"
      },
      "city": {
        "type": "text"
      },
      "state": {
        "type": "text"
      },
      "pincode": {
        "type": "text"
      },
      "location": {
        "type": "geo_point"
      }
    }
  }
}`

	exists, err := e.esClient.Indices.Exists([]string{index})
	if err != nil {
		zapLogger.Logger.Error("error in checking index exists or not in elasticSearch", zap.Error(err))
		return err
	}

	if exists.StatusCode == http.StatusOK {
		zapLogger.Logger.Error("index already exists", zap.Error(err))
		return fmt.Errorf("index already exists")
	} else if exists.StatusCode != http.StatusNotFound {
		zapLogger.Logger.Error("error in index existence response", zap.Error(err))
		return fmt.Errorf("error in index existence response")
	}

	res, err := e.esClient.Indices.Create(index, e.esClient.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		zapLogger.Logger.Error("error in creating index in elasticSearch", zap.Error(err))
		return err
	}
	if res.StatusCode != http.StatusOK {
		zapLogger.Logger.Error("error in creating index in elasticSearch", zap.Error(err))
		return errors.New("error in creating index in elasticSearch")
	}
	return nil
}

func (e *EventIndexerImpl) IndexUserEvent(event elasticsearchPkg.Event, doc []byte) error {
	res, err := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(int(event.ID)),
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)

	if err != nil {
		zapLogger.Logger.Error("error in indexing user event in elasticSearch", zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zapLogger.Logger.Error("error in closing response body", zap.Error(err))
		}
	}(res.Body)

	log.Println(res)

	return nil
}

func (e *EventIndexerImpl) GetEventsById(eventId int) ([]elasticsearchPkg.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (e *EventIndexerImpl) SearchEvents(profile model.UserProfile, page model.Pagination, filters dto.EventFilterDTO) ([]elasticsearchPkg.Event, error) {

	query := createEventsSearchQuery(profile, page, filters)
	qq, _ := json.Marshal(query)
	zapLogger.Logger.Debug("search query for events:", zap.String("query", string(qq)))

	res, err := opensearchapi.SearchRequest{
		Index: []string{e.IndexName},
		Body:  strings.NewReader(string(qq)),
	}.Do(context.Background(), e.esClient)
	if err != nil {
		zapLogger.Logger.Error("error in searching events in elasticSearch", zap.Error(err))
		return nil, err
	}

	var events []elasticsearchPkg.Event
	var elasticResponse EventElasticResponse

	if res.IsError() {
		var mp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&mp); err != nil {
			zapLogger.Logger.Error("Error parsing the response body: ", zap.Error(err))
			return nil, err
		} else {
			// Print the response status and error information.
			zapLogger.Logger.Debug(fmt.Sprintf(
				res.Status(),
				mp["error"].(map[string]interface{})["type"],
				mp["error"].(map[string]interface{})["reason"],
			))
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		zapLogger.Logger.Error("error: ", zap.Error(err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err := json.Unmarshal(bodyBytes, &elasticResponse); err != nil {
		zapLogger.Logger.Error("Error parsing the response body to profile object: ", zap.Error(err))
		return nil, err
	}

	for _, val := range elasticResponse.HitsObject.Hits {
		events = append(events, val.Source)
	}

	return events, nil
}

func createEventsSearchQuery(profile model.UserProfile, page model.Pagination, filters dto.EventFilterDTO) map[string]interface{} {

	var mustMap []map[string]interface{}

	distance := strconv.Itoa(filters.Distance) + "km"
	log.Println(distance)
	geo := map[string]interface{}{
		"geo_distance": map[string]interface{}{
			"distance": strconv.Itoa(filters.Distance) + "km",
			"location": map[string]interface{}{
				"lat": profile.Latitude,
				"lon": profile.Longitude,
			},
		},
	}
	mustMap = append(mustMap, geo)

	eventTime := map[string]interface{}{
		"range": map[string]interface{}{
			"event_time": map[string]interface{}{
				"gte": filters.StartDate.Format(IsoDateFormat),
				"lte": filters.EndDate.Format(IsoDateFormat),
			},
		},
	}
	mustMap = append(mustMap, eventTime)

	if filters.Type != "" {
		eventType := map[string]interface{}{
			"term": map[string]interface{}{
				"type": map[string]interface{}{
					"value": filters.Type,
				},
			},
		}
		mustMap = append(mustMap, eventType)
	}

	query := map[string]interface{}{
		"from": (page.PageNumber - 1) * page.PageSize, // Calculate the starting index
		"size": page.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustMap,
			},
		},
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": profile.Latitude,
						"lon": profile.Longitude,
					},
					"order":         "asc",
					"unit":          "km",
					"distance_type": "plane",
				},
			},
			{
				"expires_at": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}
	return query
}

type EventElasticResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	HitsObject struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore *float64 `json:"max_score"`
		Hits     []Hits   `json:"hits"`
	} `json:"hits"`
}

type Hits struct {
	Id     string                 `json:"id"`
	Index  string                 `json:"_index"`
	Score  string                 `json:"_score"`
	Source elasticsearchPkg.Event `json:"_source"`
	// Sort   []float64   `json:"sort"`
}

func (e *EventIndexerImpl) UpdateUserEvent(event elasticsearchPkg.Event, doc []byte) error {
	res, err := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(int(event.ID)),
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)

	if err != nil {
		zapLogger.Logger.Error("error in updating user event in elasticSearch", zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zapLogger.Logger.Error("error in closing response body", zap.Error(err))
		}
	}(res.Body)

	log.Println(res)

	return nil

}

func (e *EventIndexerImpl) DeleteUserEvent(eventID int) error {
	res, err := opensearchapi.DeleteRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(eventID),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)
	if err != nil {
		zapLogger.Logger.Error("error in deleting user event in elasticSearch", zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			zapLogger.Logger.Error("error in closing response body", zap.Error(err))
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		zapLogger.Logger.Error("error in deleting index in elasticSearch", zap.Error(err))
		return errors.New("error in deleting user event in elasticSearch")
	}
	return nil

}
