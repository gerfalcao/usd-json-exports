package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300 * time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET","http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	
	var quote Quote

	err = json.NewDecoder(res.Body).Decode(&quote)
	if err != nil {
		panic(err)
	}
	createFile(quote)
}

func createFile(quote Quote) {
	f, err := os.Create("cotacao.txt")
	_, err = f.Write([]byte(`Dólar: ` + quote.USDBRL.Bid))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file, err := os.ReadFile("cotacao.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file))
}