package config

import (
	"log"
	"os"
	"strings"
)

var AppConfig Config

type Env struct {
}

type HTTP struct {
	IP   string
	PORT string
}
type Database struct {
	URL        string
	HOST       string
	PORT       string
	USER       string
	PASS       string
	DB         string
	LogQueries bool
}

type Log struct {
	Level string
}

type AWSConfig struct {
	AccessKeyID            string
	AccessKeySecret        string
	Region                 string
	PlatformApplicationArn string
	TopicArn               string
}

type Config struct {
	Env string
	Database
	Log
	HTTP
	ElasticConfig
	RedisConfig
	TwilioConfig
	GoogleConfig
	SentryConfig
	AWSConfig
	BaseURL
}

type ElasticConfig struct {
	UserName string
	Password string
	URL      []string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type TwilioConfig struct {
	AccountSID  string
	AuthToken   string
	PhoneNumber string
}

type GoogleConfig struct {
	ClientID     string
	ClientSecret string
}

type SentryConfig struct {
	DSN string
}

type BaseURL struct {
	URL string
}

const SecretKey string = ""

var ConfigValue Config

func Load(appEnv string) (Config, error) {
	ConfigValue = Config{
		Env: appEnv,
		Log: Log{
			Level: loglevel(),
		},
		Database: Database{
			URL:        databaseURL(),
			HOST:       databaseHost(),
			PORT:       databasePort(),
			USER:       databaseUser(),
			PASS:       databasePass(),
			DB:         databaseName(),
			LogQueries: logQueries(),
		},
		HTTP: HTTP{
			IP:   hostIP(),
			PORT: httpPort(),
		},
		ElasticConfig: ElasticConfig{
			UserName: getElasticUserName(),
			Password: getElasticPassword(),
			URL:      getElasticURL(),
		},
		RedisConfig: RedisConfig{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASSWORD"),
		},
		TwilioConfig: TwilioConfig{
			AccountSID:  twilioAccountSID(),
			AuthToken:   twilioAuthToken(),
			PhoneNumber: twilioPhoneNumber(),
		},
		GoogleConfig: GoogleConfig{
			ClientID:     googleClientID(),
			ClientSecret: googleClientSecret(),
		},
		SentryConfig: SentryConfig{
			DSN: sentryDSN(),
		},
		AWSConfig: AWSConfig{
			AccessKeyID:            awsAccessKeyID(),
			AccessKeySecret:        awsAccessKeySecret(),
			Region:                 awsRegion(),
			PlatformApplicationArn: awsPlatformApplicationArn(),
			TopicArn:               awsTopicArn(),
		},
		BaseURL: BaseURL{
			URL: getBaseURL(appEnv),
		},
	}

	AppConfig = ConfigValue
	return ConfigValue, nil
}

func awsPlatformApplicationArn() string {
	return os.Getenv("AWS_PLATFORM_APPLICATION_ARN")
}

func awsTopicArn() string {
	return os.Getenv("AWS_TOPIC_ARN")
}

func awsRegion() string {
	return os.Getenv("AWS_REGION")
}

func awsAccessKeySecret() string {
	return os.Getenv("AWS_ACCESS_KEY_SECRET")
}

func awsAccessKeyID() string {
	return os.Getenv("AWS_ACCESS_KEY_ID")
}

func hostIP() string {
	host := os.Getenv("HOST_IP")

	if host == "" {
		log.Fatalln("host IP not found.")
	}
	return host
}

func httpPort() string {
	return os.Getenv("HTTP_PORT")
}

func loglevel() string {
	level := os.Getenv("LOG_LEVEL")

	if level == "" {
		return "info"
	}
	return level
}

func databaseURL() string {
	databaseURL := os.Getenv("DATABASE_URL")

	return databaseURL
}

func logQueries() bool {
	val := os.Getenv("LOG_DB_QUERIES")

	return strings.ToLower(val) == "true"
}

func getElasticURL() []string {
	val := os.Getenv("ELASTIC_URL")
	return strings.Split(val, ",")
}

func getElasticPassword() string {
	val := os.Getenv("ELASTIC_PASSWORD")
	return val
}

func getElasticUserName() string {
	val := os.Getenv("ELASTIC_USERNAME")
	return val
}

func databaseHost() string {
	host := os.Getenv("DATABASE_HOST")

	if host == "" {
		log.Fatalln("DB host not found.")
	}
	return host
}

func databasePort() string {
	port := os.Getenv("DATABASE_PORT")

	if port == "" {
		log.Fatalln("DB PORT not found.")
	}
	return port
}

func databaseUser() string {
	user := os.Getenv("DATABASE_USER")

	if user == "" {
		log.Fatalln("DB USER not found.")
	}
	return user
}

func databasePass() string {
	pass := os.Getenv("DATABASE_PASS")

	if pass == "" {
		log.Fatalln("DB PASS not found.")
	}
	return pass
}

func databaseName() string {
	db := os.Getenv("DATABASE_NAME")

	if db == "" {
		log.Fatalln("DB NAME not found.")
	}
	return db
}

func twilioAccountSID() string {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")

	if accountSID == "" {
		log.Fatalln("TWILIO_ACCOUNT_SID not found.")
	}

	return accountSID
}

func twilioAuthToken() string {
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")

	if authToken == "" {
		log.Fatalln("TWILIO_AUTH_TOKEN not found.")
	}

	return authToken
}

func twilioPhoneNumber() string {
	phoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	if phoneNumber == "" {
		log.Fatalln("TWILIO_PHONE_NUMBER not found.")
	}

	return phoneNumber
}

func googleClientID() string {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")

	if clientID == "" {
		log.Fatalln("GOOGLE_CLIENT_ID not found.")
	}

	return clientID
}

func googleClientSecret() string {
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientSecret == "" {
		log.Fatalln("GOOGLE_CLIENT_SECRET not found.")
	}

	return clientSecret
}

func sentryDSN() string {
	dsn := os.Getenv("SENTRY_DSN")

	if dsn == "" {
		log.Fatalln("SENTRY_DSN not found.")
	}
	return dsn
}

func getBaseURL(appENV string) string {
	if appENV == "dev" {
		return "http://localhost:8080"
	} else if appENV == "staging" {
		return "https://firstpluto.com"
	} else if appENV == "prod" {
		return "https://prod.com"
	} else {
		return ""
	}
}
