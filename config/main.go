package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"time"
	"webCrawler3/infra/log"
)

// Client Hints

type Config struct {
	MaxConcurrency int          `mapstructure:"MAX_CONCURRENCY"`
	MaxDepth       int          `mapstructure:"MAX_DEPTH"`
	PostgresURI    string       `mapstructure:"POSTGRES_URI"`
	AppName        string       `mapstructure:"APP_NAME"`
	TimeFormat     string       `mapstructure:"TIME_FORMAT"`
	TimeZone       string       `mapstructure:"TIME_ZONE"`
	Cache          *CacheConfig `mapstructure:"CACHE"`
	Proxy          *Proxy       `mapstructure:"PROXY"`
	Filter         *Filter      `mapstructure:"FILTER"`
	UserAgent      string       `mapstructure:"USER_AGENT"`
	ClusterMode    bool         `mapstructure:"CLUSTER_MODE"`
	ClusterAddr    uint32       `mapstructure:"CLUSTER_ADDR"`
	SeedDb         bool         `mapstructure:"SEED_DB"`
	Port           uint32       `mapstructure:"PORT"`
	CoreConfig     *CoreConfig  `mapstructure:"CORE_CONFIG"`
}
type CacheConfig struct {
	DBDir string `mapstructure:"DB_DIR"`
	Mode  string `mapstructure:"MODE"` // "mem" or "disc'
}
type Proxy struct {
	Enabled  bool   `mapstructure:"ENABLED"`
	ProxyURL string `mapstructure:"PROXY_URL"`
}
type Filter struct {
	Tlds        []string `mapstructure:"TLDS"`
	IgnoreLocal bool     `mapstructure:"IGNORE_LOCAL"`
}

var Conf *Config

func loadByFlag(cfg *Config) error {
	// Atualiza a variável global Conf
	Conf = cfg
	return nil
}

// GetAcceptableMimeType returns the acceptable mime types
func (c *Config) GetAcceptableMimeType() string {
	var queryBuilder strings.Builder
	for _, acceptMime := range c.CoreConfig.AcceptableMimeTypes.MimeTypes {
		queryBuilder.WriteString(acceptMime)
		queryBuilder.WriteString(",")
	}
	return queryBuilder.String()
}

// GetDBMaxConns returns the maximum number of connections to the database
func (c *Config) GetDBMaxConns() int32 {
	// +2 para conexões de controle
	return int32(c.MaxConcurrency) + 2
}

// GetDBMinConns returns the minimum number of connections to the database
func (c *Config) GetDBMinConns() int32 {
	// 50% do máximo
	return int32(c.MaxConcurrency / 2)
}

// GetReqTimeout returns the request timeout
func (c *Config) GetReqTimeout() time.Duration {
	if c.Proxy.Enabled {
		return 30 * time.Second
	}
	return 5 * time.Second
}

// GetCrlIdleTime returns the time to sleep when the queue is empty
func (c *Config) GetCrlIdleTime() time.Duration {
	return 5 * time.Second
}

// GetAddrManager returns the address to listen
func (c *Config) GetAddrManager() string {
	return fmt.Sprintf(":%d", c.Port)
}
func (c *Config) GetHints() map[string]string {
	return c.CoreConfig.ClientHints.Hints
}
func (c *Config) HintsEnabled() bool {
	return c.CoreConfig.ClientHints.Enabled
}

func loadByConfigFile() error {
	var cfg Config
	vip := viper.New()

	// Definindo valores padrão
	vip.SetDefault("APP_NAME", "WebCrawler")
	vip.SetDefault("PORT", 9090)
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")
	vip.SetDefault("DB_DIR", "/tmp/WebCrawler")
	vip.SetDefault("POSTGRES_URI", "postgres://postgres:Strong@P4ssword@localhost/crawler")
	vip.SetDefault("MAX_CONCURRENCY", 10)
	vip.SetDefault("MAX_DEPTH", 2)
	vip.SetDefault("USER_AGENT", "Go-http-client/1.1")

	vip.SetDefault("SEED_DB", true)

	vip.SetDefault("CLUSTER_MODE", false)
	vip.SetDefault("CLUSTER_ADDR", 6379)

	vip.SetDefault("PROXY.ENABLED", false)
	vip.SetDefault("PROXY.PROXY_URL", "http://localhost:4444")

	vip.SetDefault("FILTER.TLDS", []string{})
	vip.SetDefault("FILTER.IGNORE_LOCAL", false)

	// Definindo valores padrão para Core
	coreDefaultConfig := genDefaultConfig()
	vip.SetDefault("CORE_CONFIG.ACCEPTABLE_MIME_TYPES.MIME_TYPES", coreDefaultConfig.AcceptableMimeTypes.MimeTypes)
	vip.SetDefault("CORE_CONFIG.ACCEPTABLE_SCHEMA.SCHEMAS", coreDefaultConfig.AcceptableSchema.Schemas)
	vip.SetDefault("CORE_CONFIG.DENY_SUFFIXES.SUFFIXES", coreDefaultConfig.DenySuffixes.Suffixes)
	vip.SetDefault("CORE_CONFIG.COMMON_STOP_WORDS.LANGUAGE_STOP_WORDS", coreDefaultConfig.CommonStopWords.LanguageStopWords)
	vip.SetDefault("CORE_CONFIG.MAP_WORDS_TYPE", coreDefaultConfig.MapWordsType)
	vip.SetDefault("CORE_CONFIG.CLIENT_HINTS", coreDefaultConfig.ClientHints)

	// Lendo o arquivo de configuração conf.yml
	vip.SetConfigName("conf")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("/opt/crw")
	vip.AddConfigPath("/etc/crw")

	// Lendo as configurações do arquivo conf.yml
	if err := vip.ReadInConfig(); err != nil {
		// Se o arquivo conf.yml não for encontrado, continue sem erro
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			err := vip.WriteConfigAs("conf.yml")
			if err != nil {
				return err
			}
			log.Logger.Info("Arquivo de configuração não encontrado, criando arquivo conf.yml no diretório atual")
		}
	}
	// Se APP_NAME não estiver definido no padrão, retorne um erro
	if vip.IsSet("APP_NAME") {
		regex := regexp.MustCompile("^[A-Za-z0-9]+$")
		if !regex.MatchString(vip.GetString("APP_NAME")) {
			return errors.New("APP_NAME só pode conter letras e números")
		}
	}
	if vip.IsSet("FILTER.TLDS") {
		vip.Set("FILTER.TLDS", splitComma(vip.GetString("FILTER.TLDS")))
	}

	// Atribua as configurações ao cfg
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	// Atualiza a variável global Conf
	Conf = &cfg

	return nil
}
func LoadConfig() error {
	return loadByConfigFile()
	//var (
	//	maxConcurrency = flag.Int("maxConcurrency", 10, "Max number of concurrent requests")
	//	maxDepth       = flag.Int("maxDepth", 2, "Max depth to crawl")
	//	enabledConfig  = flag.Bool("config", false, "Enable config file")
	//	enableProxy    = flag.Bool("proxy", false, "Enable Proxy")
	//	proxyURL       = flag.String("proxyURL", "http://localhost:4444", "Proxy URL")
	//	cacheMode      = flag.Bool("mem", false, "Cache mode")
	//	// tlds list of Top-Level Domains
	//	tlds        = flag.String("tlds", "", "TLDs to filter EX: com,br,org")
	//	postgresURI = flag.String("postgresURI", "postgres://postgres:Strong@P4ssword@localhost/crawler", "Postgres URI")
	//	userAgent   = flag.String("userAgent", "Go-http-client/1.1", "User-Agent")
	//	clusterMode = flag.Bool("cluster", false, "Cluster mode")
	//	seedDb      = flag.Bool("seedDb", false, "Seed Database")
	//)
	//flag.Parse()
	//if *enabledConfig {
	//	log.Logger.Info("Carregando configurações do arquivo")
	//	return loadByConfigFile()
	//} else {
	//	log.Logger.Info("Carregando configurações por flag")
	//	return loadByFlag(&Config{
	//		AppName:        "WebCrawler",
	//		Port:           9090,
	//		TimeFormat:     "02-Jan-2006",
	//		TimeZone:       "America/Sao_Paulo",
	//		MaxConcurrency: *maxConcurrency,
	//		MaxDepth:       *maxDepth,
	//		PostgresURI:    *postgresURI,
	//		UserAgent:      *userAgent,
	//		ClusterMode:    *clusterMode,
	//		ClusterAddr:    6379,
	//		SeedDb:         *seedDb,
	//		Cache: &CacheConfig{
	//			DBDir: "/tmp/WebCrawler",
	//			Mode:  Ternary(*cacheMode, "mem", "disc"),
	//		},
	//		Proxy: &Proxy{
	//			Enabled:  *enableProxy,
	//			ProxyURL: *proxyURL,
	//		},
	//		Filter: &Filter{
	//			Tlds:        splitComma(*tlds),
	//			IgnoreLocal: false,
	//		},
	//	})
	//
	//}
}
