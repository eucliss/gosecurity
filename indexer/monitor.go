package indexer

import (
	"bufio"
	"fmt"
	"gosecurity/sys"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type LogFile struct {
	Path         string
	LastPosition int64
	TargetIndex  string
	Format       string
}

type Monitor struct {
	LogFiles []LogFile
	Watcher  *fsnotify.Watcher
}

func MonitorSources() {
	monitor, err := NewMonitor()
	if err != nil {
		log.Fatalf("Failed to initialize monitor: %v", err)
	}
	monitor.Start()
}

func BuildLogFiles() []LogFile {
	var logFiles []LogFile
	for _, source := range sys.System.Sources.FileSources {
		logFiles = append(logFiles, LogFile{
			Path:         source.Path,
			LastPosition: 0,
			TargetIndex:  source.TargetIndex,
			Format:       source.Format,
		})
	}
	return logFiles
}

func NewMonitor() (*Monitor, error) {
	logFiles := BuildLogFiles()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("error creating watcher: %v", err)
	}

	return &Monitor{
		LogFiles: logFiles,
		Watcher:  watcher,
	}, nil
}

func (m *Monitor) Start() {
	for _, lf := range m.LogFiles {
		err := m.Watcher.Add(filepath.Dir(lf.Path))
		if err != nil {
			log.Printf("Error watching %s: %v", lf.Path, err)
			continue
		}
		go m.tailFile(lf)
	}

	go func() {
		for {
			select {
			case event, ok := <-m.Watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					for i, lf := range m.LogFiles {
						if event.Name == lf.Path {
							go m.tailFile(m.LogFiles[i])
						}
					}
				}
			case err, ok := <-m.Watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()
}

func (m *Monitor) tailFile(lf LogFile) {
	file, err := os.Open(lf.Path)
	if err != nil {
		log.Printf("Error opening file %s: %v", lf.Path, err)
		return
	}
	defer file.Close()

	_, err = file.Seek(lf.LastPosition, io.SeekStart)
	if err != nil {
		log.Printf("Error seeking in file %s: %v", lf.Path, err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		m.forwardLog(line, lf.TargetIndex)
		lf.LastPosition, _ = file.Seek(0, io.SeekCurrent)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", lf.Path, err)
	}
}

func (m *Monitor) forwardLog(logEntry string, targetIndex string) {
	// Here we're using the existing Elasticsearch setup to forward logs
	// You might want to adjust this based on your specific needs
	doc := map[string]interface{}{
		"timestamp": time.Now().UTC(),
		"log_entry": logEntry,
	}

	// TODO: Add the log entry to the event stream

	sys.System.Db.InsertDocument(targetIndex, doc)
	fmt.Printf("Forwarded log entry to index %s: %s\n", targetIndex, logEntry)
}

func (m *Monitor) Stop() {
	m.Watcher.Close()
}
