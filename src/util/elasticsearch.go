package util

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"ngb/config"
	"strconv"
)

var es *elasticsearch.Client

func initEs() {
	var err error
	es, err = elasticsearch.NewDefaultClient()
	if err != nil {
		Logger.Error("Error creating the client: %s", err)
	}
	Logger.Info("elasticsearch started")
}

func buf(title string, content string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	doc := map[string]interface{}{
		"title":   title,
		"content": content,
	}
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return bytes.Buffer{}, err
	}
	return buf, nil
}

func InsertES(pid int, title string, content string) error {
	buf, err := buf(title, content)
	if err != nil {
		return err
	}
	req := esapi.IndexRequest{
		Index:      config.C.Elasticsearch.Index, // Index name
		Body:       &buf,                         // Document body
		DocumentID: strconv.Itoa(pid),            // Document ID
		Refresh:    "true",                       // Refresh
	}

	res, err := req.Do(context.Background(), es)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func Search(keyword string) ([]int, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": keyword,
						},
					},
					{
						"match": map[string]interface{}{
							"content": keyword,
						},
					},
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	// Perform the search request.
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex(config.C.Elasticsearch.Index),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var (
		r    map[string]interface{}
		pids []int
	)
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		pid, err := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))
		if err != nil {
			return nil, err
		}
		pids = append(pids, pid)
	}

	return pids, nil
}
