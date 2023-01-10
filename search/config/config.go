package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"jam3.com/common/logs"
	"log"
	"os"
)

var C = InitConfig()

type Config struct {
	viper      *viper.Viper
	SC         *ServerConfig
	GC         *GrpcConfig
	DB         *DbConfig
	ES         *EsConfig
	EtcdConfig *EtcdConfig
}
type ServerConfig struct {
	Name string
	Addr string
}
type EsConfig struct {
	Addr string
}
type EtcdConfig struct {
	Addrs []string
}
type GrpcConfig struct {
	Name    string
	Addr    string
	Version string
	Weight  int64
}
type DbConfig struct {
	Dsn string
}

func InitConfig() *Config {
	conf := &Config{viper: viper.New()}
	workDir, _ := os.Getwd()
	conf.viper.SetConfigName("app")
	conf.viper.SetConfigType("yaml")
	conf.viper.AddConfigPath(workDir + "/config")
	err := conf.viper.ReadInConfig()

	if err != nil {
		log.Fatalln(err)
	}
	conf.ReadDbConfig()
	conf.ReadServerConfig()
	conf.InitZapLog()
	conf.ReadGrpcConfig()
	conf.ReadEtcdConfig()
	conf.ReadEsConfig()
	return conf
}
func (c *Config) InitZapLog() {
	lc := &logs.LogConfig{
		DebugFileName: c.viper.GetString("zap.DebugFileName"),
		InfoFileName:  c.viper.GetString("zap.InfoFileName"),
		WarnFileName:  c.viper.GetString("zap.WarnFileName"),
		MaxSize:       c.viper.GetInt("zap.MaxSize"),
		MaxAge:        c.viper.GetInt("zap.MaxAge"),
		MaxBackups:    c.viper.GetInt("zap.MaxBackups"),
	}
	err := logs.InitLogger(lc)
	if err != nil {
		log.Fatalln(err)
	}
}
func (c *Config) ReadEsConfig() {
	ec := &EsConfig{}
	ec.Addr = c.viper.GetString("es.addr")
	c.ES = ec
}
func (c *Config) ReadDbConfig() {
	dc := &DbConfig{}
	dc.Dsn = c.viper.GetString("db.dsn")
	c.DB = dc
}

func (c *Config) ReadServerConfig() {
	sc := &ServerConfig{}
	sc.Name = c.viper.GetString("server.name")
	sc.Addr = c.viper.GetString("server.addr")
	c.SC = sc
}

func (c *Config) RedisConfig() *redis.Options {
	return &redis.Options{
		Addr:     c.viper.GetString("redis.host") + ":" + c.viper.GetString("redis.port"),
		Password: c.viper.GetString("redis.password"),
		DB:       c.viper.GetInt("redis.db"),
	}
}

func (c *Config) ReadGrpcConfig() {
	gc := &GrpcConfig{}
	gc.Name = c.viper.GetString("grpc.name")
	gc.Addr = c.viper.GetString("grpc.addr")
	gc.Version = c.viper.GetString("grpc.version")
	gc.Weight = c.viper.GetInt64("grpc.weight")
	c.GC = gc
}

func (c *Config) ReadEtcdConfig() {
	ec := &EtcdConfig{}
	var addrs []string
	err := c.viper.UnmarshalKey("etcd.addrs", &addrs)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Etcd init :", addrs)
	ec.Addrs = addrs
	c.EtcdConfig = ec
}
