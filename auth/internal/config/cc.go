package env

import (
	config "dz/auth/internal"
	"errors"
	"os"
)

const pgDsnEnv = "PG_DSN


var _ config.DBConfig =
const pgDsnEnv = "PG_DSN"
var _ config.DBConfig = (*dbConfig)(nil)
var _ config.CcConfig = (*ccConfig)(nil)

type ccConfig struct {
	port string
}
type dbConfig struct {
	dsn string
}

type AppConfig struct {
	CC *ccConfig
	DB *dbConfig
}

func NewAppConfig() (*AppConfig, error) {
	port := os.Getenv(pgDsnEnv)
	if len(port) == 0 {
		return nil, errors.New("cc port not found")
	}
	db := &dbConfig{dsn: port}

	dsn := os.Getenv("PG_PORT")
	if len(dsn) == 0 {
		return nil, errors.New("dsn not found")
	}
	cc := &ccConfig{port: port}

	return &AppConfig{
		CC: cc,
		DB: db,
	}, nil
}

func (c *ccConfig) Port() string {
	return c.port
}
func (d *dbConfig) DSN() string {
	return d.dsn
}
