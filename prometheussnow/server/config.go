package server

import (
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// the configuration for the http backend service
type Config struct {
	LogLevel     string
	Id           string
	AuthMode     string
	Path         string
	Port         string
	Username     string
	Password     string
	Metrics      bool
	snowurl      string
	snowusername string
	snowpassword string
}

func NewConfig() (*Config, error) {
	log.Info().Msg("Loading configuration.")

	// use viper to load configuration data
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	// reads the configuration file
	err := v.ReadInConfig()
	if err != nil { // handle errors reading the config file
		log.Error().Msgf("Fatal error config file: %s \n", err)
		return nil, err
	}

	// binds all environment variables to make it container friendly
	v.AutomaticEnv()
	v.SetEnvPrefix("SNOW") // prefixes all env vars
	// replace character to support environment variable format
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	_ = v.BindEnv("Id")
	_ = v.BindEnv("LogLevel")
	_ = v.BindEnv("Path")
	_ = v.BindEnv("Port")
	_ = v.BindEnv("AuthMode")
	_ = v.BindEnv("Username")
	_ = v.BindEnv("Password")
	_ = v.BindEnv("Metrics")
	_ = v.BindEnv("snowurl")
	_ = v.BindEnv("snowusername")
	_ = v.BindEnv("snowpassword")

	// creates a config struct and populate it with values
	c := new(Config)

	// general configuration
	c.Id = v.GetString("Id")
	c.LogLevel = v.GetString("LogLevel")
	c.AuthMode = v.GetString("AuthMode")
	c.Metrics = v.GetBool("Metrics")
	c.Username = v.GetString("Username")
	c.Password = v.GetString("Password")
	c.Port = v.GetString("Port")
	c.Path = v.GetString("Path")
	c.snowurl = v.GetString("snowurl")
	c.snowusername = v.GetString("snowusername")
	c.snowpassword = v.GetString("snowpassword")

	// set the log level
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logLevel, err := zerolog.ParseLevel(strings.ToLower(c.LogLevel))
	if err != nil {
		log.Warn().Msg(err.Error())
		log.Info().Msg("defaulting log level to INFO")
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.SetGlobalLevel(logLevel)

	return c, nil
}

func (c *Config) debugLevel() bool {
	return strings.ToLower(c.LogLevel) == "debug"
}
