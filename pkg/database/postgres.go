package database

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"websockets/pkg/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustGetDB() *sqlx.DB {
	connString := fmt.Sprintf(
		"user=%v password=%v host=%v port=%v dbname=%v sslmode=disable",
		viper.GetString(config.DBUser),
		viper.GetString(config.DBPassword),
		viper.GetString(config.DBHost),
		viper.GetInt(config.DBPort),
		viper.GetString(config.DBName),
	)

	//db, err := sqlx.Open("postgres", connString)
	//if err != nil {
	//	panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	//}

	db, err := otelsqlx.Open("postgres", connString,
		otelsql.WithAttributes(semconv.DBSystemSqlite))
	if err != nil {
		panic(fmt.Sprintf("Error while connecting to DB. Error: %v", err.Error()))
	}

	db.SetMaxOpenConns(viper.GetInt(config.DBMaxOpenConns))

	return db
}
