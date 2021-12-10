package config

import "time"

type Duration time.Duration
type Config struct {
	PgConfig PgConfig `json:"pg_config"`
	ListenPort int `json:"listen_port"`
	ShutdownTimeout   Duration `json:"shutdown_timeout" split_words:"true"`
}

type PgConfig struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	UserName       string `json:"user_name"`
	Password       string `json:"password"`
	Database       string `json:"database"`
	Sslmode        string `json:"sslmode"`
	MaxOpenConns   int    `json:"max_open_conns"`
	MigrationsDir  string `json:"migrations_dir"`
	MigrateVersion int    `json:"version"`
	JsonDir         string `json:"json_dir"`
}
func (d Duration) WithDefault(def time.Duration) time.Duration {
	if d == 0 {
		return def
	}
	return time.Duration(d)
}
