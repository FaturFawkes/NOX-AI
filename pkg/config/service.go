package config

import (
	"nox-ai/pkg/utils"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	ServicePort string
	Logger      *zap.Logger
	DB          *gorm.DB
	Redis       *redis.Client
}

type IServiceConfig interface {
	GetServicePort() string
	GetLogger() *zap.Logger
	GetDatabase() *gorm.DB
	GetRedis() *redis.Client
}

func InitService() IServiceConfig {

	logger := utils.Logger()

	servicePort := os.Getenv("SERVICE_PORT")
	if servicePort == "" {
		panic("service port not set")
	}

	mysqlHost := os.Getenv("MYSQL_HOST")
	if mysqlHost == "" {
		panic("mysql host cannot be nil")
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	if mysqlUser == "MYSQL_USER" {
		panic("mysql username cannot be nil")
	}

	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	if mysqlPassword == "" {
		panic("mysql password cannot be nil")
	}

	mysqlPort := os.Getenv("MYSQL_PORT")
	if mysqlPort == "" {
		panic("mysql port cannot be nil")
	}

	mysqlDb := os.Getenv("MYSQL_DB")
	if mysqlDb == "" {
		panic("mysql database cannot be nil")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		panic("redis host cannot be nil")
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		panic("redis port cannot be nil")
	}

	redisPasword := os.Getenv("REDIS_PASSWORD")

	redisDB := os.Getenv("REDIS_DB")
	if redisDB == "" {
		panic("redis db cannot be nil")
	}

	cnvRedisDb, err := strconv.Atoi(redisDB)
	if err != nil {
		panic(err)
	}

	db := NewDatabase(Mysql{
		USER:     mysqlUser,
		PASSWORD: mysqlPassword,
		PORT:     mysqlPort,
		DB:       mysqlDb,
		HOST:     mysqlHost,
	})

	redis := NewRedis(Redis{
		HOST:     redisHost,
		PASSWORD: redisPasword,
		PORT:     redisPort,
		DB:       cnvRedisDb,
	})

	return &ServiceConfig{
		Logger: logger,
		DB:     db,
		Redis:  redis,
	}
}

func (sc *ServiceConfig) GetLogger() *zap.Logger {
	return sc.Logger
}

func (sc *ServiceConfig) GetDatabase() *gorm.DB {
	return sc.DB
}

func (sc *ServiceConfig) GetRedis() *redis.Client {
	return sc.Redis
}

func (sc *ServiceConfig) GetServicePort() string {
	return sc.ServicePort
}