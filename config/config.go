package config

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
)

// Config - settings of application.
type Config struct {
	Logger Logger  `json:"logger"`
	DB     DB      `json:"db"`
	Grid   Grid    `json:"grid"`
	Statsd *Statsd `json:"statsd,omitempty"`
}

// Grid general settings
type Grid struct {
	ClientType       string     `json:"client_type"`
	Port             int        `json:"port"`
	StrategyList     []Strategy `json:"strategy_list"`
	BusyNodeDuration string     `json:"busy_node_duration"` // duration string format ex. 12m, see time.ParseDuration()
	// todo: выпилить и сделать равным дедлайну http запроса
	ReservedDuration string `json:"reserved_node_duration"` // duration string format ex. 12m, see time.ParseDuration()
}

// Strategy - Describes the algorithm of node selection.
type Strategy struct {
	Params   json.RawMessage `json:"params"` // ex. docker config, kubernetes config, etc.
	Type     string          `json:"type"`
	Limit    int             `json:"limit"`
	NodeList []Node          `json:"node_list"`
}

// Node - Describes node properties and capabilities. Applicable only for on-demand strategies.
type Node struct {
	Params           json.RawMessage          `json:"params"` // ex. image_name, etc.
	CapabilitiesList []map[string]interface{} `json:"capabilities_list"`
}

// Logger - Configuration of logger.
type Logger struct {
	Level string `json:"level"`
}

// DB - Configuration of storage.
type DB struct {
	Implementation string `json:"implementation"`
	Connection     string `json:"connection"`
}

// Statsd - Settings of metrics sender.
type Statsd struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Prefix   string `json:"prefix"`
	Enable   bool   `json:"enable"`
}

// New - Constructor of config.
func New() *Config {
	return &Config{}
}

// LoadFromFile - config loader from json file.
func (c *Config) LoadFromFile(path string) error {
	log.Printf(path)
	if path == "" {
		return errors.New("empty configuration file path")
	}

	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)

	return jsonParser.Decode(&c)
}
