package main

import (
	"context"
	"database/sql"
	"flag"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
	"github.com/rrebeiz/quicktasks/internal/data"
	"log"
	"os"
	"time"
)

type config struct {
	port        int
	environment string
	db          struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}
type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	models   data.Models
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 4000, "application port")
	flag.StringVar(&cfg.environment, "environment", "dev", "application environment dev|prod")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DSN"), "database data source name")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "database max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 30, "database max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "database max idle time")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}
	app := application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		models:   data.NewModels(db),
	}

	err = app.serve()

	if err != nil {
		errorLog.Fatal(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)
	err = db.PingContext(context.Background())
	if err != nil {
		return nil, err
	}
	return db, nil

}
