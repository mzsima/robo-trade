package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type MyEvent struct {
	Name string `json:"name"`
}

type MyTrade struct {
	ID       string `dynamodbav:"id" json:"id"`
	MidPrice string
	Orders   []string
	OrderAt  time.Time
}

type Candlestick struct {
	Open   int64   `json:"open"`
	High   int64   `json:"high"`
	Low    int64   `json:"low"`
	Close  int64   `json:"close"`
	Volume float64 `json:"volume"`
	Time   int64   `json:"time"`
}

type CandlestickResponse struct {
	SymbolID     int           `json:"symbolId"`
	Candlesticks []Candlestick `json:"candlesticks"`
	Timestamp    int64         `json:"timestamp"`
}

func (myTrade MyTrade) GetKey() map[string]types.AttributeValue {
	id, err := attributevalue.Marshal(myTrade.ID)
	if err != nil {
		panic(err)
	}
	return map[string]types.AttributeValue{"id": id}
}

type Orderbook map[string]interface{}

func getGetMyTrade(ctx context.Context) (MyTrade, error) {
	myTrade := MyTrade{ID: "mytrade"}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return MyTrade{}, err
	}

	svc := dynamodb.NewFromConfig(cfg)

	out, err := svc.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("my-table"),
		Key:       myTrade.GetKey(),
	})
	if err != nil {
		panic(err)
	} else {
		err = attributevalue.UnmarshalMap(out.Item, &myTrade)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}

	fmt.Printf("%+v \n", myTrade)
	return myTrade, nil
}

func putMyTrade(ctx context.Context, mytrade MyTrade) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return err
	}

	item, err := attributevalue.MarshalMap(mytrade)
	if err != nil {
		fmt.Printf("marshal map: %s\n", err.Error())
		return err
	}

	svc := dynamodb.NewFromConfig(cfg)
	out, err := svc.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("my-table"),
		Item:      item,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(out.Attributes)
	return err
}

func getOrderbook(ctx context.Context) (Orderbook, error) {

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://exchange.rakuten-wallet.co.jp/api/v1/orderbook", nil)
	if err != nil {
		log.Fatal(err)
	}

	// appending to existing query args
	q := req.URL.Query()
	q.Add("symbolId", "7")

	// assign encoded query string to http request
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server", err)
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// fmt.Println(resp.Status)
	// fmt.Println(string(responseBody))
	orderbook := Orderbook{}
	json.Unmarshal(responseBody, &orderbook)
	return orderbook, nil
}

func putOrderbook(ctx context.Context, orderbook Orderbook) {
	id := "ticker"
	orderbook["id"] = id

	// c, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	// if err != nil {
	// 	fmt.Printf("load aws config: %s\n", err.Error())
	// 	return
	// }
	// client := dynamodb.NewFromConfig(c)

	// tableName := "my-table"
	// input := &dynamodb.PutItemInput{
	// 	Item:      av,
	// 	TableName: aws.String(tableName),
	// }

	// _, err = svc.PutItem(input)
	// if err != nil {
	// 	log.Fatalf("Got error calling PutItem: %s", err)
	// }
}

func makePrices(mid string) []decimal.Decimal {
	x, _ := decimal.NewFromString(mid)
	res := []decimal.Decimal{}
	for i := 1; i <= 5; i++ {
		res = append(res, x.Add(x.Mul(decimal.NewFromFloat(0.02*float64(i)))))
		res = append(res, x.Add(x.Mul(decimal.NewFromFloat(-0.02*float64(i)))))
	}
	fmt.Println(res)
	return res
}

func orderPrices(prices []decimal.Decimal) []string {
	ids := []string{}
	for _, price := range prices {
		orderid := order(price)
		ids = append(ids, orderid)
	}
	return ids
}

func order(price decimal.Decimal) string {
	// TODO ---- order
	id := uuid.New()
	return id.String()
}

func getCandleStick(ctx context.Context) ([]Candlestick, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://exchange.rakuten-wallet.co.jp/api/v1/candlestick", nil)
	if err != nil {
		log.Fatal(err)
	}

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	// appending to existing query args
	q := req.URL.Query()
	q.Add("symbolId", "7")
	q.Add("candlestickType", "PT5M")
	q.Add("dateFrom", fmt.Sprintf("%d", fiveMinutesAgo.UnixMilli()))

	// assign encoded query string to http request
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Errored when sending request to the server", err)
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println(resp.Status)
	fmt.Println(string(responseBody))
	candlestickResponse := CandlestickResponse{}
	json.Unmarshal(responseBody, &candlestickResponse)
	return candlestickResponse.Candlesticks, nil
}

func insertCandleStick(ctx context.Context, candlesticks []Candlestick) {
	db, err := sql.Open("mysql", "admin:password@tcp(host:3306)/test")
	if err != nil {
		fmt.Println("Failed to open database:", err)
		return
	}
	defer db.Close()

	for _, candlestick := range candlesticks {
		timestamp := time.Unix(candlestick.Time/1000, 0)
		_, err = db.Exec("INSERT INTO candlestick (time, open, high, low, close, volume) VALUES (?, ?, ?, ?, ?, ?)",
			timestamp, candlestick.Open, candlestick.High, candlestick.Low, candlestick.Close, candlestick.Volume)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Data inserted successfully.")
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	candlesticks, err := getCandleStick(ctx)
	if err != nil {
		log.Fatal(err)
	}
	insertCandleStick(ctx, candlesticks)

	return fmt.Sprintf("Hello %s!", name.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
