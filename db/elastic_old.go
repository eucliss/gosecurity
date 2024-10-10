package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

func Start() {

	caCert, err := os.ReadFile("db/http_ca.crt") // Replace with the path to your CA certificate
	if err != nil {
		log.Fatalf("Error reading CA certificate: %s", err)
	}

	// Create a CA certificate pool and add your CA cert
	// caCertPool := x509.NewCertPool()
	// caCertPool.AppendCertsFromPEM(caCert)

	// Create a custom HTTP transport with TLS configuration
	// tlsConfig := &tls.Config{
	// 	RootCAs: caCertPool, // Use the CA pool with your cert
	// }

	// transport := &http.Transport{
	// 	TLSClientConfig: tlsConfig, // Add your TLS config
	// }

	// Elasticsearch configuration
	cfg := elasticsearch8.Config{
		Addresses: []string{
			"https://localhost:9200", // If running Go app locally
			// "http://elasticsearch:9200", // If running Go app inside Docker
		},
		// If you set up authentication, include Username and Password
		Username: "elastic",              // Default user for Elasticsearch Docker image
		Password: "MoEO249xSwsNW3oWEc5F", // Replace with actual password if set
		CACert:   caCert,
	}

	// Create Elasticsearch client
	es, err := elasticsearch8.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	// Check Elasticsearch version
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting Elasticsearch info: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Connected to Elasticsearch")

	index := "test-index"
	mapping := `
    {
      "settings": {
        "number_of_shards": 1
      },
      "mappings": {
        "properties": {
          "title": {
            "type": "text"
          },
		  "content": {
			"type": "text"
		  }
        }
      }
    }`

	res, err = es.Indices.Create(
		index,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
	// Index a new document (JSON body)
	doc := map[string]interface{}{
		"title":   "Go and Elasticsearch",
		"content": "This is a document indexed by the Go Elasticsearch client.",
	}
	docJSON, _ := json.Marshal(doc)
	res, err = es.Index("test-index", bytes.NewReader(docJSON))
	if err != nil {
		log.Fatalf("Error indexing document: %s", err)
	}
	defer res.Body.Close()
	fmt.Println("Document indexed successfully")

	// Refresh the index to make the document searchable
	_, err = es.Indices.Refresh(es.Indices.Refresh.WithIndex("test-index"))
	if err != nil {
		log.Fatalf("Error refreshing index: %s", err)
	}

	// Perform a search query
	query := `{
		"query": {
			"match": {
				"title": "Go"
			}
		}
	}`

	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test-index"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error executing search query: %s", err)
	}
	defer res.Body.Close()

	// Parse and print search results
	var searchResult map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		log.Fatalf("Error parsing the search results: %s", err)
	}
	fmt.Printf("Search result: %+v\n", searchResult)

	// Iterate the document "hits" returned by API call
	for _, hit := range searchResult["hits"].(map[string]interface{})["hits"].([]interface{}) {

		// Parse the attributes/fields of the document
		doc := hit.(map[string]interface{})

		// The "_source" data is another map interface nested inside of doc
		source := doc["_source"]
		fmt.Println("doc _source:", reflect.TypeOf(source))

		// Get the document's _id and print it out along with _source data
		docID := doc["_id"]
		fmt.Println("docID:", docID)
		fmt.Println("_source:", source)
	} // end of response iteration

}
