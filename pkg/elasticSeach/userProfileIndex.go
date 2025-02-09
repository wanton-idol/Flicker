package elasticSeach

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/opensearch-project/opensearch-go/opensearchapi"

	elasticsearchPkg "github.com/SuperMatch/model/elasticSearch"
	"github.com/SuperMatch/zapLogger"
	"github.com/opensearch-project/opensearch-go"
)

var userProfileIndex = "user_profile"

//go:generate mockgen -package mocks -destination mocks/userProfileIndex_mock.go github.com/SuperMatch/pkg/elasticSeach ElasticSearchIndexer
type ElasticSearchIndexer interface {
	IndexUserProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error
	UpdateUserProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error
	GetUserProfile(userProfileId int) (elasticsearchPkg.UserProfile, error)
	SearchProfile(query map[string]interface{}) ([]elasticsearchPkg.UserProfile, error)
	UpdateSearchProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error
	CreateIndex() error
}

type ElasticSearchIndexerImpl struct {
	IndexName string
	esClient  *opensearch.Client
}

func NewElasticSearchIndexerImpl() *ElasticSearchIndexerImpl {
	return &ElasticSearchIndexerImpl{
		IndexName: userProfileIndex,
		esClient:  EsClient,
	}
}

func (e *ElasticSearchIndexerImpl) IndexUserProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error {

	res, err := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(userProfile.Id),
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}

func (e *ElasticSearchIndexerImpl) UpdateUserProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error {
	res, err := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(userProfile.Id),
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}.Do(context.Background(), e.esClient)

	if err != nil {
		return err
	}

	defer res.Body.Close()
	return nil
}

func (e *ElasticSearchIndexerImpl) GetUserProfile(userProfileId int) (elasticsearchPkg.UserProfile, error) {

	userProfile := elasticsearchPkg.UserProfile{}

	res, err := e.esClient.Get(
		e.IndexName,
		strconv.Itoa(userProfileId),
		e.esClient.Get.WithContext(context.Background()),
	)

	if err != nil {
		zapLogger.Logger.Error("Error while getting user profile: ", zap.Error(err))
		return userProfile, err
	}

	userProfile, err = parseEsApiResponse(res)

	if err != nil {
		zapLogger.Logger.Error("Error while parsing user profile: ", zap.Error(err))
		return userProfile, err
	}

	defer res.Body.Close()
	return userProfile, nil
}

func (e *ElasticSearchIndexerImpl) SearchProfile(query map[string]interface{}) ([]elasticsearchPkg.UserProfile, error) {

	userProfiles := []elasticsearchPkg.UserProfile{}

	qq, _ := json.Marshal(query)

	zapLogger.Logger.Debug(string(qq))

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		zapLogger.Logger.Error("Error encoding query:", zap.Error(err))
	}

	res, err := e.esClient.Search(
		e.esClient.Search.WithContext(context.Background()),
		e.esClient.Search.WithIndex(userProfileIndex),
		e.esClient.Search.WithBody(&buf),
		e.esClient.Search.WithTrackTotalHits(true),
		e.esClient.Search.WithPretty(),
	)

	if err != nil {
		zapLogger.Logger.Error("Error getting response: ", zap.Error(err))
	}

	defer res.Body.Close()

	var elasticResponse elasticsearchPkg.ElasticResponse

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			zapLogger.Logger.Error("Error parsing the response body: ", zap.Error(err))
			return nil, err
		} else {
			// Print the response status and error information.
			zapLogger.Logger.Error(fmt.Sprintf(res.Status(), e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"]))

		}
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	zapLogger.Logger.Debug(bodyString)

	if err := json.Unmarshal(bodyBytes, &elasticResponse); err != nil {
		zapLogger.Logger.Error("Error parsing the response body to profile object: ", zap.Error(err))
		return nil, err
	}

	var result QueryResults
	result.TotalCount = elasticResponse.HitsObject.Total.Value
	for _, val := range elasticResponse.HitsObject.Hits {
		userProfiles = append(userProfiles, val.Source)
	}
	result.Hits = userProfiles
	return userProfiles, nil
}

func (e *ElasticSearchIndexerImpl) UpdateSearchProfile(userProfile elasticsearchPkg.UserProfile, doc []byte) error {

	zapLogger.Logger.Info(string(doc))

	req := opensearchapi.IndexRequest{
		Index:      e.IndexName,
		DocumentID: strconv.Itoa(userProfile.Id),
		Body:       strings.NewReader(string(doc)),
		Refresh:    "true",
	}

	res, err := req.Do(context.Background(), e.esClient)

	if err != nil {
		zapLogger.Logger.Error("error in updating SearchProfile", zap.Error(err))
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		zapLogger.Logger.Error(res.String())
		return errors.New(res.String())
	}
	return nil
}

//
//################################
//.
//.
//..
//################################

// private methods for Elasticsearch
type esapiResponse struct {
	Source elasticsearchPkg.UserProfile `json:"_source"`
}

func parseEsApiResponse(response *opensearchapi.Response) (elasticsearchPkg.UserProfile, error) {

	elasticResponse := esapiResponse{}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	zapLogger.Logger.Debug(bodyString)

	if err := json.Unmarshal(bodyBytes, &elasticResponse); err != nil {
		zapLogger.Logger.Error("Error parsing the response body to profile object: ", zap.Error(err))
		return elasticResponse.Source, err
	}
	return elasticResponse.Source, err
}

type QueryResults struct {
	Hits       []elasticsearchPkg.UserProfile
	TotalCount int
}

func (e *ElasticSearchIndexerImpl) CreateIndex() error {
	index := e.IndexName
	file, err := os.ReadFile("elasticSearchIndex.es")
	if err != nil {
		log.Println("error in reading the file")
		return err
	}

	mapping := string(file)

	exists, err := e.esClient.Indices.Exists([]string{index})
	if err != nil {
		log.Println("error in checking index exists or not in elasticSearch")
		return err
	}

	if exists.StatusCode == 200 {
		log.Println("index already exists")
		return fmt.Errorf("index already exists")
	} else if exists.StatusCode != 404 {
		log.Println("error in index existence response")
		return fmt.Errorf("error in index existence response")
	}

	res, err := e.esClient.Indices.Create(index, e.esClient.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		log.Println("error in creating index in elasticSearch")
		return err
	}
	if res.StatusCode != 200 {
		log.Println("error in creating index in elasticSearch")
		return err
	}

	return nil
}
