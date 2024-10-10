package indexer

import (
	"fmt"
	"gosecurity/config"
	"gosecurity/db"
	"gosecurity/sys"
)

type IngestEvent struct {
	Index string
	Data  []byte
}

var FileChannel = make(chan config.FileSource)
var EventStream = make(chan IngestEvent)

func prepareIndex(index string) {
	exists, _ := sys.System.DbConfig.IndexExists(index)
	if exists {
		return
	}
	fmt.Println("Index not found, creating: " + index)
	newIndex := db.Index{
		Name: index,
		Mapping: `
		{
		  "settings": {
			"number_of_shards": 1
		  },
		  "mappings": {
			"properties": {
			  "Name": {
				"type": "text"
			  },
			  "Description": {
				"type": "text"
			  },
			}
		  }
		}`,
	}
	sys.System.Db.CreateIndices(newIndex)
}

func Ingest() {
	for {
		e := <-EventStream
		go prepareIndex(e.Index)
		sys.System.Db.InsertDocument(e.Index, e.Data)
	}
}

func Open() {
	fmt.Println("Opening ingestion channel (FileChannel)")
	for {
		f := <-FileChannel
		prepareIndex(f.TargetIndex)

		fmt.Printf("Ingesting file: %s:%s\n", f.Name, f.Format)
		switch f.Format {
		case "csv":
			go ProcessCSV(f.Path, f.TargetIndex)
		case "json":
			data, _ := LoadFile(f)
			go ProcessJSON(data, f.TargetIndex)
		case "xml":
			data, _ := LoadFile(f)
			go ProcessXML(data, f.TargetIndex)
			fmt.Println("Parsing XML file")
		case "html":
			fmt.Println("Parsing HTML file")
		default:
			fmt.Println("Unknown file type")

		}
	}
}
