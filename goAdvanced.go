package main
import (
	"bytes"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/kr/pretty"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"time"
)


/*
import "github.com/streadway/amqp"
import "golang.org/x/net/websocket"
import _ "github.com/go-sql-driver/mysql"
import "github.com/emirpasic/gods/sets"
import "github.com/emirpasic/gods/lists"
import "github.com/emirpasic/gods/stacks"
import "github.com/emirpasic/gods/maps"
import "github.com/emirpasic/gods/utils"
import "github.com/pkg/errors"
import "gopkg.in/natefinch/lumberjack.v2"
*/


//--------------redis-----------
func testRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

//--------------mux-----------
func testMux() {
	r := mux.NewRouter()
	r.HandleFunc("/test/{category}", testHandler)
	http.Handle("/test/1234", r)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}

//--------------gorilla(websocket)-----------
func testGorillaWebSocket() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.EnableCompression = true
}

//--------------decimal-----------
func testDecimal() {
	decimal.DivisionPrecision = 18
	price1, err := decimal.NewFromString("1.23333333333333")
	//price2, _ := decimal.NewFromString("1.245780147245780148")

	var tmpPrice = decimal.Zero
	startTime := time.Now()
	for i := 0; i < 1000000; i++ {
		price1.Add(price1)
		//tmpPrice = price1
		//tmpPrice = tmpPrice.Mul(tmpPrice)
		//tmpPrice = tmpPrice.Div(tmpPrice)
		/*if price1.LessThan(price2) {
			//fmt.Println("price1 is LessThan price2 !")
		}*/
	}
	if err != nil {
		panic(err)
	}
	endTime := time.Now()
	fmt.Println("DiffTime is ", endTime.Sub(startTime))
	fmt.Println(tmpPrice)
}

//--------------cast-----------
func testCast() {
	date, _ := cast.StringToDate("1985-10-16 12:23:45")
	fmt.Println(date)
}

//-------------cobra-----------
func testCobra() {
	var cmdPrint = &cobra.Command{
		Use:   "Use",
		Short: "short",
	}
	fmt.Println(cmdPrint)
}

//--------------pretty-----------
func testPretty() {
	type myType struct {
		a, b int
	}
	var x = []myType{{1, 2}, {3, 4}, {5, 6}}
	fmt.Println(x)
	fmt.Printf("%# v\n", pretty.Formatter(x))
}

//--------------viper-----------
func testViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	var yamlExample = []byte(`
		Hacker: true
		name: Robin Shilf
		hobbies:
		- skateboarding
		- snowboarding
		- go
		clothing:
		jacket: leather
		trousers: denim
		age: 35
		eyes : brown
		beard: true
		`)
	viper.ReadConfig(bytes.NewBuffer(yamlExample))
	fmt.Println(viper.Get("name"))
}

//--------------json-----------
func testJSON() {
	type ColorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	group := ColorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	b, _ := jsoniter.Marshal(group)
	fmt.Println(b)
}

