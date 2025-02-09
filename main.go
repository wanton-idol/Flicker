package main

import (
	"github.com/SuperMatch/zapLogger"
	"log"
	"time"

	"github.com/SuperMatch/pkg/redis"
	"github.com/getsentry/sentry-go"

	_ "github.com/SuperMatch/docs"

	"github.com/SuperMatch/config"
	"github.com/SuperMatch/env"
	pkgdb "github.com/SuperMatch/pkg/db"
	"github.com/SuperMatch/pkg/elasticSeach"
	"github.com/SuperMatch/server"
	"go.uber.org/zap"
)

//	@title			API documentation
//	@version		1.0
//	@description	This is a sample server.

//	@host	localhost:8080

//@securityDefinitions.apikey	ApiKeyAuth
//	@in				header
//	@name			token
//	@description	token used to authenticate

func main() {

	//cmd := flag.String("env", "dev", "")
	//flag.Parse()
	//appEnv := *cmd
	appEnv, err := env.GetAppEnv()

	if err != nil {
		log.Fatal("unable to load APP_ENV variable.")
	}

	err = env.LoadConfigFile(appEnv)

	if err != nil {
		log.Fatal("unable to load config file.")
	}

	config, err := config.Load(appEnv)

	if err != nil {
		log.Fatal("unable to read config file.")
	}

	log.Printf("i have read %v", config.Database)

	logger := zapLogger.InitLogger(appEnv, config)

	_, err = pkgdb.Load(&config.Database, appEnv)

	if err != nil {
		logger.Fatal("error in connecting database ", zap.Error(err))
	} else {
		logger.Log(zap.InfoLevel, "database connection successful !")
	}

	err = elasticSeach.CreateElasticClient(config)

	if err != nil {
		logger.Fatal("error in connecting elasticsearch ", zap.Error(err))
	} else {
		logger.Log(zap.InfoLevel, "elastic search connection successful !")
	}

	err = redis.CreateRedisClient(config)

	if err != nil {
		logger.Fatal("error in connecting redis", zap.Error(err))
	} else {
		logger.Log(zap.InfoLevel, "redis connection successful !")
	}

	if config.Env == "staging" || config.Env == "prod" {
		//create sentry client
		err = sentry.Init(sentry.ClientOptions{
			Dsn:         config.SentryConfig.DSN,
			Environment: appEnv,
		})
		if err != nil {
			log.Fatalf("Sentry initialization failed: %v", err)
		}

		defer sentry.Flush(2 * time.Second)
	}

	deps := &server.RouterDeps{
		AppEnv:    appEnv,
		AppConfig: config,
	}

	router, err := server.ConfigRouter(deps)

	if err != nil {
		log.Fatal("Could not initialize Router", zap.Error(err))
	}

	// signals := make(chan os.Signal, 1)
	// signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	// wg := &sync.WaitGroup{}

	// server.StartServer(router, config.HTTP, wg)

	// defer httpServer.Stop()

	// defer dbConn.Close()

	router.Run(config.HTTP.IP + ":" + config.HTTP.PORT)
}
