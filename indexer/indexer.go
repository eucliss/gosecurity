package indexer

import (
	"fmt"
	"gosecurity/sys"
)

// var Database db.Config

type NormalizedEvent struct {
	Event map[string]interface{}
	Index string
}

var NormalizedChannel = make(chan NormalizedEvent)

func Store() {
	fmt.Println("Opening store channel (StoreChannel)")
	for {
		e := <-NormalizedChannel
		fmt.Println("Storing event.", e.Event)
		sys.System.Db.InsertDocument(e.Index, e.Event)
	}
}
