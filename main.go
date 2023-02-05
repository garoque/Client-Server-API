package main

import (
	"time"

	"github.com/garoque/Client-Server-API/client"
	"github.com/garoque/Client-Server-API/database"
	"github.com/garoque/Client-Server-API/server"
)

func main() {
	dbConn := database.CreateConnection()
	defer dbConn.Close()

	database.CreateTableExchange(dbConn)

	go func() {
		time.Sleep(time.Second * 3)

		client.SendRequest()
	}()

	server.StartServer(dbConn)
}
