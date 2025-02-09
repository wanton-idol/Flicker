package elasticSeach

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SuperMatch/model"
	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/zapLogger"
	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"go.uber.org/zap"
	"io"
	"log"
	"strings"
	"time"
)

const (
	userStoriesIndex = "user_stories"
	IsoDateFormat    = "2006-01-02"
)

type UserStoriesIndexer interface {
	CreateIndex() error
	IndexUserStories(userStories elasticsearchPkg.UserStories, doc []byte) error
	GetUserStoriesByProfileID(userProfileID int) ([]elasticsearchPkg.UserStories, error)
	GetUserStoriesByLocation(location model.UserLocation) ([]elasticsearchPkg.UserStories, error)
}

type UserStoriesIndexerImpl struct {
	IndexName string
	esClient  *opensearch.Client
	logger    *zap.Logger
}

func NewUserStoriesIndexerImpl() *UserStoriesIndexerImpl {
	return &UserStoriesIndexerImpl{
		IndexName: userStoriesIndex,
		esClient:  EsClient,
		logger:    zapLogger.Logger,
	}
}

func (e *UserStoriesIndexerImpl) CreateIndex() error {
	index := e.IndexName
	mapping := `{
		"settings": {
	   	"number_of_shards": 2,
	   	"number_of_replicas": 1
		},
	  "mappings": {
	      "properties": {
	        "id": {
	           "type": "text"
	        },
	        "user_profile_id": {
	           "type": "integer"
	        },
	        "text": {
	           "type": "text"
	        },
	        "media_url": {
	           "type": "text"
	        },
			 "media_type": {
				"type": "text"
			 },
	        "location": {
	          "type": "geo_point"
	        },
	        "expires_at": {
	           "type": "date"
	        },
	        "created_at": {
	           "type": "date"
	        }
	    }
	  }
	}`
	exists, err := e.esClient.Indices.Exists([]string{index})
	if err != nil {
		e.logger.Error("error in checking index exists or not in elasticSearch", zap.Error(err))
		return err
	}

	if exists.StatusCode == 200 {
		e.logger.Error("index already exists", zap.Error(err))
		return fmt.Errorf("index already exists")
	} else if exists.StatusCode != 404 {
		e.logger.Error("error in index existence response", zap.Error(err))
		return fmt.Errorf("error in index existence response")
	}

	res, err := e.esClient.Indices.Create(index, e.esClient.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		e.logger.Error("error in creating index in elasticSearch", zap.Error(err))
		return err
	}
	if res.StatusCode != 200 {
		e.logger.Error("error in creating index in elasticSearch", zap.Error(err))
		return err
	}

	return nil
}

func (e *UserStoriesIndexerImpl) IndexUserStories(userStories elasticsearchPkg.UserStories, doc []byte) error {
	log.Println(e.esClient)

	res, err := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: userStories.ID,
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)

	if err != nil {
		e.logger.Error("error in indexing user stories in elasticSearch", zap.Error(err))
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			e.logger.Error("error in closing response body", zap.Error(err))
		}
	}(res.Body)

	log.Println(res)

	return nil
}

func QueryStoriesByProfileID(userProfileID int) map[string]interface{} {

	var mustMap []map[string]interface{}

	x := map[string]interface{}{
		"match": map[string]interface{}{
			"user_profile_id": userProfileID,
		},
	}
	mustMap = append(mustMap, x)

	y := map[string]interface{}{
		"range": map[string]interface{}{
			"expires_at": map[string]interface{}{
				"gte": time.Now().Format(IsoDateFormat),
			},
		},
	}

	mustMap = append(mustMap, y)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustMap,
			},
		},
		"sort": []map[string]interface{}{
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
	return query
}

func (e *UserStoriesIndexerImpl) GetUserStoriesByProfileID(userProfileID int) ([]elasticsearchPkg.UserStories, error) {
	query := QueryStoriesByProfileID(userProfileID)
	qq, _ := json.Marshal(query)
	res, err := opensearchapi.SearchRequest{
		Index: []string{e.IndexName},
		Body:  strings.NewReader(string(qq)),
	}.Do(context.Background(), e.esClient)
	if err != nil {
		e.logger.Error("error in searching user stories in elasticSearch", zap.Error(err))
		return nil, err
	}

	var userStories []elasticsearchPkg.UserStories
	var elasticResponse elasticsearchPkg.StoriesResponse

	if res.IsError() {
		var mp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&mp); err != nil {
			e.logger.Error("Error parsing the response body: ", zap.Error(err))
			return nil, err
		} else {
			// Print the response status and error information.
			log.Println("",
				res.Status(),
				mp["error"].(map[string]interface{})["type"],
				mp["error"].(map[string]interface{})["reason"],
			)
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		e.logger.Error("error: ", zap.Error(err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err := json.Unmarshal(bodyBytes, &elasticResponse); err != nil {
		e.logger.Error("Error parsing the response body to profile object: ", zap.Error(err))
		return nil, err
	}

	for _, val := range elasticResponse.HitsObject.Hits {
		userStories = append(userStories, val.Source)
	}

	return userStories, nil

}

func QueryStoriesByLocation(location model.UserLocation) map[string]interface{} {

	var mustMap []map[string]interface{}

	x := map[string]interface{}{
		"geo_distance": map[string]interface{}{
			"distance": "200km",
			"location": map[string]interface{}{
				"lat": location.Latitude,
				"lon": location.Longitude,
			},
		},
	}
	mustMap = append(mustMap, x)

	y := map[string]interface{}{
		"range": map[string]interface{}{
			"expires_at": map[string]interface{}{
				"gte": time.Now().Format(IsoDateFormat),
			},
		},
	}

	mustMap = append(mustMap, y)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustMap,
			},
		},
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": location.Latitude,
						"lon": location.Longitude,
					},
					"order":         "asc",
					"unit":          "km",
					"distance_type": "plane",
				},
			},
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
	}
	return query
}

func (e *UserStoriesIndexerImpl) GetUserStoriesByLocation(location model.UserLocation) ([]elasticsearchPkg.UserStories, error) {
	query := QueryStoriesByLocation(location)
	qq, _ := json.Marshal(query)
	res, err := opensearchapi.SearchRequest{
		Index: []string{e.IndexName},
		Body:  strings.NewReader(string(qq)),
	}.Do(context.Background(), e.esClient)
	if err != nil {
		e.logger.Error("error in searching user stories in elasticSearch", zap.Error(err))
		return nil, err
	}

	var userStories []elasticsearchPkg.UserStories
	var elasticResponse elasticsearchPkg.StoriesResponse

	if res.IsError() {
		var mp map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&mp); err != nil {
			e.logger.Error("Error parsing the response body: ", zap.Error(err))
			return nil, err
		} else {
			// Print the response status and error information.
			log.Println("",
				res.Status(),
				mp["error"].(map[string]interface{})["type"],
				mp["error"].(map[string]interface{})["reason"],
			)
		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		e.logger.Error("error: ", zap.Error(err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	if err := json.Unmarshal(bodyBytes, &elasticResponse); err != nil {
		e.logger.Error("Error parsing the response body to profile object: ", zap.Error(err))
		return nil, err
	}

	for _, val := range elasticResponse.HitsObject.Hits {
		userStories = append(userStories, val.Source)
	}

	return userStories, nil
}
