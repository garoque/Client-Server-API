package client

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const URL = "http://localhost:8080/cotacao"

func SendRequest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal(err)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	f.Write([]byte("Dólar: " + string(b) + "\n"))
	defer f.Close()
}
