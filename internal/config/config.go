package config

type (
	PAuth struct {
		Server Server
		Psql   Psql
		Smtp   Smtp
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
)
