package main

import (
	"fmt"
	"gosecurity/config"
	"gosecurity/db"
	"gosecurity/indexer"
	"gosecurity/sys"
)

func InitializeElastic(def config.MonitoringConfig) (c db.Config, res *db.Config) {
	fmt.Println("Initializing Elasticsearch...")

	c = db.Config{
		Location: def.Database.CertLocation,
		Address:  fmt.Sprintf("https://%s:%s", def.Database.Host, def.Database.Port),
	}
	c.SetUsername(def.Database.User)
	c.SetPassword(def.Database.Password)
	c.Initialize()

	// newIndex := db.Index{
	// 	Name: def.Database.NewIndex,
	// 	Mapping: `
	// 	{
	// 	  "settings": {
	// 		"number_of_shards": 1
	// 	  },
	// 	  "mappings": {
	// 		"properties": {
	// 		  "Name": {
	// 			"type": "text"
	// 		  },
	// 		  "Description": {
	// 			"type": "text"
	// 		  },
	// 		  "Hostname": {
	// 			"type": "text"
	// 		  },
	// 		  "Time": {
	// 			"type": "text"
	// 		  }
	// 		}
	// 	  }
	// 	}`,
	// }
	// c.CreateIndices(newIndex)
	res = &c
	return
}

func DeleteIndex(c db.Config, index string) {
	c.DeleteIndex(index)
}

func main() {

	// Load the base config
	base_config, err := config.LoadConfig()
	sys.System.Monitor = base_config
	if err != nil {
		fmt.Println("Error loading base config:", err)
		return
	}

	// Load the sources config
	var source_config config.SourcesConfig
	source_config, s_err := config.LoadSources()
	if s_err != nil {
		fmt.Println("Error loading Sources config:", s_err)
		return
	}
	sys.System.Sources = source_config

	// Initialize the Elastic DB
	db_config, db_addr := InitializeElastic(base_config)
	// indexer.Database = c
	// var d db.Database = c_addr

	sys.System.DbConfig = db_config
	sys.System.Db = db_addr

	fmt.Println("DB Configured, available indices:")
	fmt.Println(db_config.GetIndices())
	// Delete the index after the function returns
	// Garbage collection
	// defer DeleteIndex(c, cfg.Database.NewIndex)

	// Start the ingest engine
	go indexer.Open()
	go indexer.Store()

	fmt.Println("Ingest engine started, waiting for events...")

	// Emit 1 event a second
	// go indexer.EmitEvents(20)

	// Get the Alert
	// alert.SetAlertPath("alert/alerts/")

	// a, _ := alert.Load("machine_fuel_consumption.yaml")
	// a, _ := alert.Load("network_status.yaml")

	// fmt.Println("Running the alert.")
	// results, triggered := a.Run()
	// fmt.Println(results)
	// fmt.Println(triggered)
	// fmt.Println("Done running the alert")
	// sys.System.Db.Query(cfg.Database.NewIndex, alertConfig.Query)

	// Query the doc to test
	fmt.Println("Querying the DB")
	r := sys.System.Db.Query("network_data", `
		{
			"query": {
				"match_all": {}
			}
		}
	`)

	fmt.Println(sys.System.Db.GetResults(r))

	// 	d.Query(cfg.Database.NewIndex, `
	// 	{
	// 		"query": {
	// 		  "match": {
	// 			"machine_id": "Irrigation_System_1"
	// 		  }
	// 		},
	// 		"size": 1
	// 	  }
	// `)

	// q := `{"query": "SELECT * FROM test-index WHERE status = ALLOW" }`
	// d.Query(cfg.Database.NewIndex,
	// q,
	// )

	// Let it bang
	select {}
	return
}

// func DoElasticStuff() {
// 	// c := db.Config{
// 	// 	Location: "db/http_ca.crt",
// 	// 	Address:  "https://localhost:9200",
// 	// }
// 	// c.SetUsername("elastic")
// 	// c.SetPassword("MoEO249xSwsNW3oWEc5F")
// 	// fmt.Println(c.Username())
// 	// c.Initialize()

// 	// newIndex := db.Index{
// 	// 	Name: "test-index2",
// 	// 	Mapping: `
// 	// 	{
// 	// 	  "settings": {
// 	// 		"number_of_shards": 1
// 	// 	  },
// 	// 	  "mappings": {
// 	// 		"properties": {
// 	// 		  "title": {
// 	// 			"type": "text"
// 	// 		  },
// 	// 		  "content": {
// 	// 			"type": "text"
// 	// 		  }
// 	// 		}
// 	// 	  }
// 	// 	}`,
// 	// }
// 	// c.CreateIndices(newIndex)
// 	// fmt.Println("Index created")

// 	c.InsertDocument(
// 		"test-index",
// 		map[string]interface{}{
// 			"title":   "New Title!!",
// 			"content": "This is a document indexed by the Go Elasticsearch client.",
// 		},
// 	)
// 	fmt.Println("Document inserted")

// 	c.Query("test-index", `{
// 		"query": {
// 			"match": {
// 				"title": "New"
// 			}
// 		}
// 	}`)

// 	c.DeleteIndex("test-index")

// }
