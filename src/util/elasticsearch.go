package util

//var es *elasticsearch.Client
//
//func initEs() {
//	var err error
//	es, err = elasticsearch.NewDefaultClient()
//	if err != nil {
//		Logger.Error("Error creating the client: %s", err)
//	}
//	res, err := es.Info()
//	defer res.Body.Close()
//	if err != nil {
//		Logger.Error("Error getting response: %s", err)
//	}
//	Logger.Info("elasticsearch started \nversion:  " + elasticsearch.Version)
//}
//
//func buf(title string, content string) bytes.Buffer {
//	var buf bytes.Buffer
//	doc := map[string]interface{}{
//		"title":   title,
//		"content": content,
//	}
//	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
//		Logger.Error(err, "Error encoding doc")
//	}
//	return buf
//}
//
//func InsertES(pid int, title string, content string) {
//	buf := buf(title, content)
//	req := esapi.IndexRequest{
//		Index:      "post_test",       // Index name
//		Body:       &buf,              // Document body
//		DocumentID: strconv.Itoa(pid), // Document ID
//		Refresh:    "true",            // Refresh
//	}
//
//	res, err := req.Do(context.Background(), es)
//	if err != nil {
//		log.Fatalf("Error getting response: %s", err)
//	}
//	defer res.Body.Close()
//
//	log.Println(res)
//
//	res, err = es.Get("post_test", strconv.Itoa(pid))
//	if err != nil {
//		Logger.Error(err, "Error get response")
//	}
//	defer res.Body.Close()
//	fmt.Println(res.String())
//}
//
//func SearchES() {
//	var r map[string]interface{}
//	var buf bytes.Buffer
//	query := map[string]interface{}{
//		"query": map[string]interface{}{
//			"match": map[string]interface{}{
//				"title": "te33",
//			},
//		},
//	}
//	if err := json.NewEncoder(&buf).Encode(query); err != nil {
//		Logger.Error("Error encoding query: %s", err)
//	}
//
//	// Perform the search request.
//	res, err := es.Search(
//		es.Search.WithContext(context.Background()),
//		es.Search.WithIndex("post_test"),
//		es.Search.WithBody(&buf),
//		es.Search.WithTrackTotalHits(true),
//		es.Search.WithPretty(),
//	)
//	if err != nil {
//		Logger.Error("Error getting response: %s", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		var e map[string]interface{}
//		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
//			Logger.Error("Error parsing the response body: %s", err)
//		} else {
//			// Print the response status and error information.
//			Logger.Error("[%s] %s: %s",
//				res.Status(),
//				e["error"].(map[string]interface{})["type"],
//				e["error"].(map[string]interface{})["reason"],
//			)
//		}
//	}
//
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		Logger.Error("Error parsing the response body: %s", err)
//	}
//	// Print the response status, number of results, and request duration.
//	fmt.Printf(
//		"[%s] %d hits; took: %dms",
//		res.Status(),
//		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
//		int(r["took"].(float64)),
//	)
//	// Print the ID and document source for each hit.
//	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
//		fmt.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
//	}
//
//}
//
//func Search() {
//	var buf bytes.Buffer
//	query := map[string]interface{}{
//		"query": map[string]interface{}{
//			"match": map[string]interface{}{
//				"title": "te33",
//			},
//		},
//		//"highlight": map[string]interface{}{
//		//	"pre_tags" : []string{"<font color='red'>"},
//		//	"post_tags" : []string{"</font>"},
//		//	"fields" : map[string]interface{}{
//		//		"title" : map[string]interface{}{},
//		//	},
//		//},
//	}
//	if err := json.NewEncoder(&buf).Encode(query); err != nil {
//		Logger.Error(err, "Error encoding query")
//	}
//	// Perform the search request.
//	res, err := es.Search(
//		es.Search.WithContext(context.Background()),
//		es.Search.WithIndex("post_test"),
//		es.Search.WithBody(&buf),
//		es.Search.WithTrackTotalHits(true),
//		es.Search.WithPretty(),
//	)
//	if err != nil {
//		Logger.Error(err, "Error getting response")
//	}
//	defer res.Body.Close()
//	fmt.Println(res.String())
//}
