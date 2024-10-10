package indexer

import (
	"fmt"
	"gosecurity/config"
	"time"
)

// type Events struct {
// 	ids chan int
// }

// func (e Events) String() string {
// 	return fmt.Sprintf("Events{ids: %v}", e.ids)
// }

// func (e Events) increment(curr int) {
// 	e.ids <- curr + 1
// }

// func (e Events) event() (res Event) {
// 	id := <-e.ids
// 	defer e.increment(id)
// 	res = Event{
// 		Name:        "Event" + fmt.Sprint(id),
// 		Description: "New Security Event",
// 		Hostname:    "hostname" + fmt.Sprint(id),
// 		Time:        "2020-01-01T00:00:00Z",
// 	}
// 	return
// }

// func initialize() (e Events) {
// 	e = Events{
// 		ids: make(chan int, 1),
// 	}
// 	e.ids <- 0
// 	return
// }

func EmitEvents(interval int) {
	// e := initialize()

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()
	fmt.Println("Emitting events every", interval, "seconds")

	// Event loop
	for range ticker.C {
		go func() {
			FileChannel <- config.FileSource{
				Name:        "Network",
				Format:      "json",
				Description: "Test file",
				Path:        "indexer/data/network.json",
				Source:      "network",
				TargetIndex: "network_data",
			}
		}()
	}

}
