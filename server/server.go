package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Quote struct {
	ID    		int 	`gorm:"primaryKey"`
	UsdValue	float64 `gorm:"column:usd_value"`
	gorm.Model
}

type QuoteApi struct {
	USDBRL struct {
        Bid         string
    }	
}

func main() {
	http.HandleFunc("/cotacao", handlerExchange)
	http.ListenAndServe(":8080", nil)
}

func handlerExchange(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 200 * time.Millisecond)
	defer cancel()
	exchangeUrl := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequestWithContext(ctx, "GET", exchangeUrl, nil)
	if err 	!= nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	} 
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(os.Stdout, bytes.NewReader(body))

	db := runDatabase()
	if db == nil {
		http.Error(w, "Failed to initialize database", http.StatusInternalServerError)
		return
	}
	recordQuote(db, body)
	
	w.Write(body)
}

func recordQuote(db *gorm.DB,res []byte) *gorm.DB {
	fmt.Println("Recording exchange")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Millisecond)
	defer cancel()

	var quoteApi QuoteApi
	err := json.Unmarshal(res, &quoteApi)
	if err != nil {
		panic(err)
	}

	bid, err := strconv.ParseFloat(quoteApi.USDBRL.Bid, 64)
    if err != nil {
        panic(err)
    }

	fmt.Println(bid)

	quote := Quote{
		UsdValue: bid,
	}
	db.WithContext(ctx).Create(&quote)
	return db
}

func runDatabase() *gorm.DB {
	fmt.Println("Running database")
	dsn := "root:root@tcp(localhost:3306)/dolar_db?utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if db.Migrator().HasTable(&Quote{}) {
		return db
	}
	db.AutoMigrate(&Quote{})
	return db
}
