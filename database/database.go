package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/garoque/Client-Server-API/model"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./ClientServerAPI.db")
	if err != nil {
		log.Fatalln("CreateConnection err: ", err.Error())
	}

	return db
}

func CreateTableExchange(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `exchange` (`id` VARCHAR(64) PRIMARY KEY, `code` VARCHAR(64), `codein` VARCHAR(64), `name` VARCHAR(64), `high` VARCHAR(64), `low` VARCHAR(64), `var_bid` VARCHAR(64), `pct_change` VARCHAR(64), `bid` VARCHAR(64), `ask` VARCHAR(64), `timestamp` DATETIME, `create_date` DATETIME)")
	if err != nil {
		log.Fatalln("CreateTableExchange err: ", err.Error())
	}
}

func AddExchangeValue(ctx context.Context, db *sql.DB, value model.ResponseUSDBRL) {
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	id := uuid.New().String()

	stmt, _ := db.Prepare("INSERT INTO exchange (id, code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.ExecContext(ctx, id, value.USDBRL.Code, value.USDBRL.Codein, value.USDBRL.Name, value.USDBRL.High, value.USDBRL.Low, value.USDBRL.VarBid, value.USDBRL.PctChange, value.USDBRL.Bid, value.USDBRL.Ask, value.USDBRL.Timestamp, value.USDBRL.CreateDate)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
}

func ReadAllExchangeValues(db *sql.DB) ([]model.ExchangeValue, error) {
	rows, err := db.Query("SELECT * FROM exchange")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []model.ExchangeValue{}
	for rows.Next() {
		i := model.ExchangeValue{}
		err = rows.Scan(&i.ID, &i.Ask, &i.Bid, &i.Code, &i.Codein, &i.CreateDate, &i.High, &i.Low, &i.Name, &i.PctChange, &i.Timestamp, &i.VarBid)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	return data, nil
}
