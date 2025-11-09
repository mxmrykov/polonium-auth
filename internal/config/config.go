package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type (
	PAuth struct {
		Server Server
		Psql   Psql
		Smtp   Smtp
		Redis  Redis
		Vault  Vault
		Auth   Auth
	}

	Auth struct {
		Issuer          string
		TokenSecret     string
		CertSecret      string
		Access, Refresh time.Duration
		CertExp         time.Duration
	}

	Server struct {
		Port string
	}

	Psql struct {
		Host, DB, User, Pass string
		MaxConn              int
	}

	Smtp struct {
		Host, Port, Sender, Password string
	}

	Redis struct {
		Host, Password string
		Db             int
	}

	Vault struct {
		Address, Token, MountPath string
	}
)

func Init() *PAuth {
	return &PAuth{
		Server: Server{Port: ":8080"},
		Psql:   loadPsql(),
		Smtp:   loadSmtp(),
		Redis:  loadRedis(),
		Vault:  loadVault(),
		Auth:   loadAuth(),
	}
}

func loadPsql() Psql {
	return Psql{
		Host:    envRequired[string]("PGSQL_HOST"),
		DB:      envRequired[string]("PGSQL_DB_NAME"),
		User:    envRequired[string]("PGSQL_USER"),
		Pass:    envRequired[string]("PGSQL_PASSWORD"),
		MaxConn: 10,
	}
}
func loadSmtp() Smtp {
	return Smtp{
		Host:     envRequired[string]("SMTP_HOST"),
		Port:     envRequired[string]("SMTP_PORT"),
		Sender:   envRequired[string]("SMTP_SENDER"),
		Password: envRequired[string]("SMTP_PASSWORD"),
	}
}
func loadRedis() Redis {
	return Redis{
		Host:     envRequired[string]("REDIS_HOST"),
		Password: envRequired[string]("REDIS_PASSWORD"),
		Db:       envDefault[int]("REDIS_DB", 0),
	}
}

func loadVault() Vault {
	return Vault{
		Address:   envRequired[string]("VAULT_ADDRESS"),
		Token:     envRequired[string]("VAULT_TOKEN"),
		MountPath: envDefault[string]("VAULT_MOUNT_PATH", "secret"),
	}
}

func loadAuth() Auth {
	return Auth{
		Issuer:      envDefault[string]("APP_AUTH_ISSUER", "polonium-authorization"),
		TokenSecret: envRequired[string]("APP_AUTH_TOKEN_SECRET"),
		CertSecret:  envRequired[string]("APP_AUTH_CERT_SECRET"),
		Access:      envDefault[time.Duration]("APP_ACCESS_TTL", time.Minute),
		Refresh:     envDefault[time.Duration]("APP_REFRESH_TTL", time.Hour),
		CertExp:     envDefault[time.Duration]("APP_CERT_TTL", time.Hour),
	}
}

func envRequired[T interface{ time.Duration | string | int }](name string) T {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("environment variable %s is required", name)
	}

	var result T
	switch any(result).(type) {
	case time.Duration:
		dur, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("environment variable %s must be a valid duration: %v", name, err)
		}
		result = any(dur).(T)
	case string:
		result = any(v).(T)
	case int:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("environment variable %s must be a valid integer: %v", name, err)
		}
		result = any(intVal).(T)
	}

	return result
}

func envDefault[T interface{ time.Duration | string | int }](name string, def T) T {
	v := os.Getenv(name)
	if v == "" {
		return def
	}

	var result T
	switch any(result).(type) {
	case time.Duration:
		dur, err := time.ParseDuration(v)
		if err != nil {
			log.Fatalf("environment variable %s must be a valid duration: %v", name, err)
		}
		result = any(dur).(T)
	case string:
		result = any(v).(T)
	case int:
		intVal, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("environment variable %s must be a valid integer: %v", name, err)
		}
		result = any(intVal).(T)
	}

	return result
}
