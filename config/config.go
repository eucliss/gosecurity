package config

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Struct to represent the Database section
type DatabaseConfig struct {
	Type         string `yaml:"type"`
	CertLocation string `yaml:"cert_location"`
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	NewIndex     string `yaml:"new_index_on_launch"`
}

// Struct to represent the full configuration
type MonitoringConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// func LoadConfig() (MonitoringConfig, error) {
// 	data, err := os.ReadFile("config.yaml")
// 	if err != nil {
// 		return MonitoringConfig{}, fmt.Errorf("error reading YAML file: %w", err)
// 	}

// 	var config MonitoringConfig

// 	// Unmarshal the YAML data into the Config struct
// 	err = yaml.Unmarshal(data, &config)
// 	return config, err
// }
