package config

type MysqlDB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	Tz       string
}

type RedisDB struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type Jwt struct {
	Key string
	Exp int
}

type Server struct {
	Host string
	Port string
}

type Config struct {
	Server  Server
	MysqlDB MysqlDB
	RedisDB RedisDB
	Jwt     Jwt
}
