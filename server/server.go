package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/garoque/Client-Server-API/database"
	"github.com/garoque/Client-Server-API/model"
)

const URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

var dbConn *sql.DB

func StartServer(db *sql.DB) {
	log.Println("Server Started")

	dbConn = db

	http.HandleFunc("/cotacao", HandlerFunc)
	http.ListenAndServe(":8080", nil)
}

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*200)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var data model.ResponseUSDBRL
	json.Unmarshal(b, &data)

	database.AddExchangeValue(ctx, dbConn, data)
	rows, err := database.ReadAllExchangeValues(dbConn)
	if err != nil {
		log.Println(err.Error())
	}

	for _, row := range rows {
		fmt.Println(row)
	}

	w.Write([]byte(data.USDBRL.Bid))
}
