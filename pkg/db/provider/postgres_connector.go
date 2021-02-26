package provider

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"time"
)


type postgresConnector struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSL      string
	db       *sql.DB
}

const connectionPattern = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

func NewPostgresConnector(host, port, username, password, dbName, dbSsl string) (PostgresConnector, error) {
	psqlInfo := fmt.Sprintf(connectionPattern, host, port, username,
		password, dbName, dbSsl)
	db, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		return nil, ErrPGCreateConnector
	}
	return postgresConnector{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		DBName:   dbName,
		SSL:      dbSsl,
		db:       db,
	}, nil
}


func (pc postgresConnector) DBConnection() *sql.DB {
	return pc.db
}

func (pc postgresConnector) PingLoop() {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for range ticker.C {
		err := pc.db.Ping()
		if err != nil {
			panic(err)
		}
	}
}
