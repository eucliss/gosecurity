package sys

// Full System Data structures

import (
	"gosecurity/config"
	"gosecurity/db"
)

type SystemStruct struct {
	Monitor  config.MonitoringConfig
	Sources  config.SourcesConfig
	Db       db.Database
	DbConfig db.Config
}

var System SystemStruct
