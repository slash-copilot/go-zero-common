package config

type PostgresConfig struct {
	DSN         string
	Debug       bool `json:",default=false"`
	AutoMigrate bool `json:",default=false"`
}
