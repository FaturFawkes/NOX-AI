package config

import (
	"github.com/FaturFawkes/NOX-AI/pkg/utils"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ServiceConfig struct {
	ServicePort string
	Logger      *zap.Logger
	DB          *gorm.DB
	Redis       *redis.Client
	GPT         *openai.Client
	Whatsapp    Whatsapp
}

type IServiceConfig interface {
	GetServicePort() string
	GetLogger() *zap.Logger
	GetDatabase() *gorm.DB
	GetRedis() *redis.Client
	GetGPT() *openai.Client
	GetWhatsapp() Whatsapp
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

	apikeyGPT := os.Getenv("APIKEY_OPENAI")
	if apikeyGPT == "" {
		panic("apikey gpt cannot be nil")
	}

	tokenWa := os.Getenv("TOKEN_WA")
	if tokenWa == "" {
		panic("Whatsapp Token cannot be nil")
	}

	hostWa := os.Getenv("HOST_WA")
	if hostWa == "" {
		panic("Whatsapp Host cannot be nil")
	}

	versionWa := os.Getenv("VERSION_WA")
	if versionWa == "" {
		panic("Whatsapp Version cannot be nil")
	}

	numberWa := os.Getenv("NUMBER_ID_WA")
	if numberWa == "" {
		panic("Whatsapp phone number id cannot be nil")
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

	gpt := NewGPT(apikeyGPT)

	whatsapp := Whatsapp{
		Token:   tokenWa,
		Host:    hostWa,
		Version: versionWa,
		Number:  numberWa,
	}

	return &ServiceConfig{
		Logger:      logger,
		GPT:         gpt,
		DB:          db,
		Redis:       redis,
		ServicePort: servicePort,
		Whatsapp:    whatsapp,
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

func (sc *ServiceConfig) GetGPT() *openai.Client {
	return sc.GPT
}

func (sc *ServiceConfig) GetWhatsapp() Whatsapp {
	return sc.Whatsapp
}
