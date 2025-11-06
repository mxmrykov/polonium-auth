package config

import (
	"log"
	"os"
	"strconv"
)

type (
	PAuth struct {
		Server Server
		Psql   Psql
		Smtp   Smtp
		Redis  Redis
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
)

func Init() *PAuth {
	return &PAuth{
		Server: Server{Port: ":8080"},
		Psql:   loadPsql(),
		Smtp:   loadSmtp(),
		Redis:  loadRedis(),
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
		Db:       envRequired[int]("REDIS_DB"),
	}
}

func envRequired[T interface{ string | int }](name string) T {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("environment variable %s is required", name)
	}

	var result T
	switch any(result).(type) {
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
