package main

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestGetOrderbook(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	getOrderbook(ctx)
	defer cancel()
}

func TestPutOrderbook(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()
	ticker, _ := getOrderbook(ctx)
	putOrderbook(ctx, ticker)
}

func TestMakePrices(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	// defer cancel()
	// orderbook, _ := getOrderbook(ctx)
	// makePrices(fmt.Sprintf("%v", orderbook["midPrice"]))
	makePrices("3043050")
}

func TestGetMyTrade(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()
	getGetMyTrade(ctx)
}

func TestPutMyTrade(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()
	myTrade := MyTrade{
		ID:       "mytrade",
		MidPrice: "100.00",
		Orders:   []string{"12", "12321", "bbb13"},
		OrderAt:  time.Now(),
	}
	putMyTrade(ctx, myTrade)
}

func TestGetCandlestick(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	getCandleStick(ctx)
	defer cancel()
}

func TestInsertCandlestick(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer cancel()
	candlesticks, err := getCandleStick(ctx)
	if err != nil {
		log.Fatal(err)
	}
	insertCandleStick(ctx, candlesticks)
}
