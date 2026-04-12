// Package config contains app config struct
package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config struct for app config
type Config struct {
	Global struct {
		ServerIP           string `yaml:"server_ip" env:"GLOBAL_SERVER_IP"`
		ProjectName        string `yaml:"project_name" env:"GLOBAL_PROJECT_NAME" env-default:""`
		ProjectFrontendUrl string `yaml:"project_frontend_url" env:"GLOBAL_PROJECT_FRONTEND_URL" env-default:""`
		ProjectApiUrl      string `yaml:"project_api_url" env:"GLOBAL_PROJECT_API_URL" env-default:""`
		CpanelFrontendUrl  string `yaml:"cpanel_frontend_url" env:"GLOBAL_CPANEL_FRONTEND_URL" env-default:""`
		CpanelApiUrl       string `yaml:"cpanel_api_url" env:"GLOBAL_CPANEL_API_URL" env-default:""`
		TenantMainDomain   string `yaml:"tenant_main_domain" env:"GLOBAL_TENANT_MAIN_DOMAIN" env-default:""`
		TenantUrlScheme    string `yaml:"tenant_url_scheme" env:"GLOBAL_TENANT_URL_SCHEME" env-default:""`
	} `yaml:"global"`
	CPanelApp struct {
		Name    string `yaml:"name" env:"CPANEL_APP_NAME" env-default:"cpanel"`
		Version string `yaml:"version" env:"CPANEL_APP_VERSION" env-default:"1.0.0"`
		Base    struct {
			StartTimeoutSec int  `yaml:"start_timeout_sec" env:"CPANEL_APP_BASE_START_TIMEOUT_SEC" env-default:"10"`
			StopTimeoutSec  int  `yaml:"stop_timeout_sec" env:"CPANEL_APP_BASE_STOP_TIMEOUT_SEC" env-default:"2"`
			IsProd          bool `yaml:"is_prod" env:"CPANEL_APP_BASE_IS_PROD" env-default:"false"`
			UseFxLogger     bool `yaml:"use_fx_logger" env:"CPANEL_APP_BASE_USE_FX_LOGGER" env-default:"true"`
			UseLogger       bool `yaml:"use_logger" env:"CPANEL_APP_BASE_USE_LOGGER" env-default:"true"`
			LogSQLQueries   bool `yaml:"log_sql_queries" env:"CPANEL_APP_BASE_LOG_SQL_QUERIES" env-default:"true"`
			LogHTTP         bool `yaml:"log_http" env:"CPANEL_APP_BASE_LOG_HTTP" env-default:"true"`
		} `yaml:"base"`
		HTTP struct {
			Port             int      `yaml:"port" env:"CPANEL_APP_HTTP_PORT" env-default:"8081"`
			Prefix           string   `yaml:"prefix" env:"CPANEL_APP_HTTP_PREFIX" env-default:"/admin-api"`
			UnderProxy       bool     `yaml:"under_proxy" env:"CPANEL_APP_HTTP_UNDER_PROXY" env-default:"false"`
			StopTimeoutSec   int      `yaml:"stop_timeout_sec" env:"CPANEL_APP_HTTP_STOP_TIMEOUT_SEC" env-default:"3"`
			CorsAllowOrigins []string `yaml:"cors_allow_origins" env:"CPANEL_APP_HTTP_CORS_ALLOW_ORIGINS" env-default:""`
		} `yaml:"http"`
		Auth struct {
			// nolint
			LifeTimeWithoutActivitySec int `yaml:"life_time_without_activity_sec" env:"CPANEL_APP_AUTH_LIFE_TIME_WITHOUT_ACTIVITY_SEC" env-default:"3600"`
			// nolint
			JWTSecret string `yaml:"jwt_secret" env:"CPANEL_APP_AUTH_JWT_SECRET" env-default:""`
		} `yaml:"auth"`
	} `yaml:"cpanel_app"`
	BackendApp struct {
		Name    string `yaml:"name" env:"BACKEND_APP_NAME" env-default:"backend"`
		Version string `yaml:"version" env:"BACKEND_APP_VERSION" env-default:"1.0.0"`
		Base    struct {
			StartTimeoutSec int  `yaml:"start_timeout_sec" env:"BACKEND_APP_BASE_START_TIMEOUT_SEC" env-default:"10"`
			StopTimeoutSec  int  `yaml:"stop_timeout_sec" env:"BACKEND_APP_BASE_STOP_TIMEOUT_SEC" env-default:"2"`
			IsProd          bool `yaml:"is_prod" env:"BACKEND_APP_BASE_IS_PROD" env-default:"false"`
			UseFxLogger     bool `yaml:"use_fx_logger" env:"BACKEND_APP_BASE_USE_FX_LOGGER" env-default:"true"`
			UseLogger       bool `yaml:"use_logger" env:"BACKEND_APP_BASE_USE_LOGGER" env-default:"true"`
			LogSQLQueries   bool `yaml:"log_sql_queries" env:"BACKEND_APP_BASE_LOG_SQL_QUERIES" env-default:"true"`
			LogHTTP         bool `yaml:"log_http" env:"BACKEND_APP_BASE_LOG_HTTP" env-default:"true"`
		} `yaml:"base"`
		HTTP struct {
			Port             int      `yaml:"port" env:"BACKEND_APP_HTTP_PORT" env-default:"8080"`
			Prefix           string   `yaml:"prefix" env:"BACKEND_APP_HTTP_PREFIX" env-default:""`
			StopTimeoutSec   int      `yaml:"stop_timeout_sec" env:"BACKEND_APP_HTTP_STOP_TIMEOUT_SEC" env-default:"3"`
			CorsAllowOrigins []string `yaml:"cors_allow_origins" env:"BACKEND_APP_HTTP_CORS_ALLOW_ORIGINS" env-default:""`
		} `yaml:"http"`
		GRPC struct {
			Port               int  `yaml:"port" env:"BACKEND_APP_GRPC_PORT" env-default:"50051"`
			LogResponseSent    bool `yaml:"log_response_sent" env:"BACKEND_APP_GRPC_LOG_RESPONSE_SENT" env-default:"false"`
			LogPayloadReceived bool `yaml:"log_payload_received" env:"BACKEND_APP_GRPC_LOG_PAYLOAD_RECEIVED" env-default:"false"`
		} `yaml:"grpc"`
	} `yaml:"backend_app"`
	CronjobApp struct {
		Name    string `yaml:"name" env:"CRONJOB_APP_NAME" env-default:"cronjob"`
		Version string `yaml:"version" env:"CRONJOB_APP_VERSION" env-default:"1.0.0"`
		Base    struct {
			StartTimeoutSec int  `yaml:"start_timeout_sec" env:"CRONJOB_APP_BASE_START_TIMEOUT_SEC" env-default:"10"`
			StopTimeoutSec  int  `yaml:"stop_timeout_sec" env:"CRONJOB_APP_BASE_STOP_TIMEOUT_SEC" env-default:"2"`
			IsProd          bool `yaml:"is_prod" env:"CRONJOB_APP_BASE_IS_PROD" env-default:"false"`
			UseFxLogger     bool `yaml:"use_fx_logger" env:"CRONJOB_APP_BASE_USE_FX_LOGGER" env-default:"true"`
			UseLogger       bool `yaml:"use_logger" env:"CRONJOB_APP_BASE_USE_LOGGER" env-default:"true"`
			LogSQLQueries   bool `yaml:"log_sql_queries" env:"CRONJOB_APP_BASE_LOG_SQL_QUERIES" env-default:"true"`
			LogJob          bool `yaml:"log_job" env:"CRONJOB_APP_BASE_LOG_JOB" env-default:"true"`
		} `yaml:"base"`
		Jobs struct {
			Autostart            bool `yaml:"autostart" env:"CRONJOB_APP_JOBS_AUTOSTART" env-default:"false"`
			ProcessFilesToDelete struct {
				// nolint
				TimeoutMillisec int `yaml:"timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_FILES_TO_DELETE_TIMEOUT_MILLISEC" env-default:"5000"`
				// nolint
				FailedTimeoutMillisec int `yaml:"failed_timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_FILES_TO_DELETE_FAILED_TIMEOUT_MILLISEC" env-default:"5000"`
			} `yaml:"process_files_to_delete"`
			ProcessUnusedFiles struct {
				// nolint
				TimeoutMillisec int `yaml:"timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_UNUSED_FILES_TIMEOUT_MILLISEC" env-default:"5000"`
				// nolint
				FailedTimeoutMillisec int `yaml:"failed_timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_UNUSED_FILES_FAILED_TIMEOUT_MILLISEC" env-default:"5000"`
				// nolint
				UnusedTtlMin int `yaml:"unused_ttl_min" env:"CRONJOB_APP_JOBS_PROCESS_UNUSED_FILES_UNUSED_TTL_MIN" env-default:"1440"`
			} `yaml:"process_unused_files"`
			ProcessEmailsToSend struct {
				// nolint
				TimeoutMillisec int `yaml:"timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_EMAILS_TO_SEND_TIMEOUT_MILLISEC" env-default:"2000"`
				// nolint
				FailedTimeoutMillisec int `yaml:"failed_timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_EMAILS_TO_SEND_FAILED_TIMEOUT_MILLISEC" env-default:"5000"`
			} `yaml:"process_emails_to_send"`
			ProcessEmailsToDelete struct {
				// nolint
				TimeoutMillisec int `yaml:"timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_EMAILS_TO_DELETE_TIMEOUT_MILLISEC" env-default:"3600000"`
				// nolint
				FailedTimeoutMillisec int `yaml:"failed_timeout_millisec" env:"CRONJOB_APP_JOBS_PROCESS_EMAILS_TO_DELETE_FAILED_TIMEOUT_MILLISEC" env-default:"5000"`
				// nolint
				TtlMin int `yaml:"ttl_min" env:"CRONJOB_APP_JOBS_PROCESS_EMAILS_TO_DELETE_TTL_MIN" env-default:"43200"`
			} `yaml:"process_emails_to_delete"`
		} `yaml:"jobs"`
	} `yaml:"cronjob_app"`
	Postgres struct {
		MaxAttempts         int    `yaml:"max_attempts" env:"POSTGRES_MAX_ATTEMPTS" env-default:"3"`
		AttemptSleepSeconds int    `yaml:"attempt_sleep_seconds" env:"POSTGRES_ATTEMPT_SLEEP_SECONDS" env-default:"1"`
		MigrationsPath      string `yaml:"migrations_path" env:"POSTGRES_MIGRATIONS_PATH" env-default:"migrations"`
		Master              struct {
			DSN string `yaml:"dsn" env:"POSTGRES_MASTER_DSN"`
		} `yaml:"master"`
	} `yaml:"postgres"`
	Storage struct {
		UpMigrations     bool   `yaml:"up_migrations" env:"STORAGE_UP_MIGRATIONS" env-default:"false"`
		S3Endpoint       string `yaml:"s3_endpoint" env:"STORAGE_S3_ENDPOINT" env-default:""`
		S3AccessKey      string `yaml:"s3_access_key" env:"STORAGE_S3_ACCESS_KEY" env-default:""`
		S3SecretKey      string `yaml:"s3_secret_key" env:"STORAGE_S3_SECRET_KEY" env-default:""`
		S3Region         string `yaml:"s3_region" env:"STORAGE_S3_REGION" env-default:""`
		S3URL            string `yaml:"s3_url" env:"STORAGE_S3_URL" env-default:""`
		S3URLIsHost      bool   `yaml:"s3_url_is_host" env:"STORAGE_S3_URL_IS_HOST" env-default:"false"`
		S3URLHostPrefix  string `yaml:"s3_url_host_prefix" env:"STORAGE_S3_URL_HOST_PREFIX" env-default:""`
		S3URLHostPostfix string `yaml:"s3_url_host_postfix" env:"STORAGE_S3_URL_HOST_POSTFIX" env-default:""`
	} `yaml:"storage"`
	Emailing struct {
		SmtpHost         string `yaml:"smtp_host" env:"EMAILING_SMTP_HOST" env-default:""`
		SmtpPort         int    `yaml:"smtp_port" env:"EMAILING_SMTP_PORT" env-default:""`
		SmtpUser         string `yaml:"smtp_user" env:"EMAILING_SMTP_USER" env-default:""`
		SmtpPassword     string `yaml:"smtp_password" env:"EMAILING_SMTP_PASSWORD" env-default:""`
		DefaultFromEmail string `yaml:"default_from_email" env:"EMAILING_DEFAULT_FROM_EMAIL" env-default:""`
		DefaultFromTitle string `yaml:"default_from_title" env:"EMAILING_DEFAULT_FROM_TITLE" env-default:"api"`
	} `yaml:"emailing"`
	Alerts struct {
		BotToken      string `yaml:"bot_token" env:"ALERTS_BOT_TOKEN" env-default:""`
		TargetChannel int64  `yaml:"target_channel" env:"ALERTS_TARGET_CHANNEL"`
	} `yaml:"alerts"`
	Auth struct {
		AccessTokenLifetime  time.Duration `yaml:"access_token_lifetime" env:"AUTH_ACCESS_TOKEN_LIFETIME"`
		RefreshTokenLifetime time.Duration `yaml:"refresh_token_lifetime" env:"AUTH_REFRESH_TOKEN_LIFETIME"`
		JwtAccessSecret      string        `yaml:"jwt_access_secret" env:"AUTH_JWT_ACCESS_SECRET"`
		JwtRefreshSecret     string        `yaml:"jwt_refresh_secret" env:"AUTH_JWT_REFRESH_SECRET"`
	} `yaml:"auth"`
	Openai struct {
		Token        string `yaml:"token" env:"OPENAI_TOKEN"`
		BaseURL      string `yaml:"base_url" env:"OPENAI_BASE_URL"`
		DefaultModel string `yaml:"default_model" env:"OPENAI_DEFAULT_MODEL"`
	} `yaml:"openai"`
}

// LoadConfig loads app config from file
func LoadConfig(files ...string) Config {
	var Config Config

	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			err := cleanenv.ReadConfig(file, &Config)
			if err != nil {
				log.Println("config file error", err)
			}
		} else {
			log.Println("config file not found", file)
		}
	}

	return Config
}
