package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	_ "log" //log包没有被使用，_ 这只会执行包的 init()方法
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/kr/pretty"
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

//****************************************************
//---------------------Normal Test--------------------
//****************************************************

func testVariables() {
	//一般不这么赋值
	var i int
	i = 1
	var j = 0

	k := 0 //这种用的比较多，注意不要加var
	//var t 	---- Error
	//var t:= 9 ---- Error


	fmt.Println(i, j, k)


	x := 2
	y := 3
	//x := 4 ---- Error，x已经声明了
	{
		x := 4 	//这样也可以，作用域仅限于括号内，且这个变量在括号内必须被使用
		fmt.Println(x)
	}
	x, z := 5, 6 //这样就可以，因为z没有被声明
	//值互换
	x, y = y, x

	//nil 是 interface、function、pointer、map、slice 和 channel 类型变量的默认初始值。
	//但声明时不指定类型，编译器也无法推断出变量的具体类型。
	//var t = nil --- Error
	//t := nil --- Error
	var t interface{} = nil

	fmt.Println(z, t)

	//const q 	----Error，必须初始化
	const q float32 = 0.5
	//q = 2.3	----Error，不能改变

	const v = 10.5
	//const v := 10  ---- Error，const有var一样的语法

	const (
		OKCOIN_CN  = "okcoin.cn"
		OKCOIN_COM = "okcoin.com"
		OKEX       = "okex.com"
	)
	//const这种常量定义比较常用，还有就是下面的枚举类型比较常用，也就这两种
}

func testEnum() {
	const (
		Sunday = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
		numberofDays //大写字母包外可见，小写字母包外不可见
	)
}

func testKind() {
	//这个涉及到go中的类型转换，go的类型转换比较的严格
	var tt = func(a float32) {
	}
	var zz = func(a int) {}
	var a float32 = 2.0
	b := 2 // int

	//a = b	---- Error 无法自动转换
	a = float32(b)
	//tt(b) ----- Error float64，float32无法自动转换
	tt(3.4)
	tt(3)
	//zz(3.4) ----- Error 无法转换
	zz(3)

	var c int16 = 5
	//zz(c) ----- Error int int8 int16 int32 int64变量间也无法自动转换
	//c := (a == b) ----Error  不同类型无法比较

	fmt.Println(a, b, c)
}

func testStr() {
	//var s1 string = nil  ---- Error 字符串类型的零值是空串 ""，不能赋值为nil

	strNum := "34"
	n, _ := strconv.Atoi(strNum) //这个有err返回, n是int类型
	n++
	strNum = strconv.Itoa(n) //这个没有err返回

	strX := "Hello Go !"
	ch := strX[0] //无法改值，只能读值
	fmt.Println("ch is ", ch)

	//byte等价于uint8，常用来处理ascii字符，go语言没有char类型
	var c byte = 'a'
	//var cn byte = '好' ----- Error 越界
	var cn rune = '好'

	//rune等同于int32
	var uc = "hello, 你好 "

	//strX += 123 ----- Error

	// + 这种合并方式效率非常低，每合并一次，都是创建一个新的字符串,就必须遍历复制一次字符串
	tmp := "123"
	strX += tmp

	//这个速度比较快，缺点是Buffer线程不安全
	var buffer bytes.Buffer
	buffer.WriteString(tmp)
	tmp = buffer.String()

	//第一个是计算字节数，第二，三是计算字符数，二三等价
	fmt.Println(len(uc), utf8.RuneCountInString(uc), len([]rune(uc)))
	fmt.Printf("%d,%s,%c,%c", len(strX), strX, c, cn)

	//原生字符串，支持换行 ，\n会被原样输出
	var raw = `
		ttttttt
		zzzzzzzzz\n
		ssssssssssss
	`

	//go中，字符串的本质就是一个字节数组，字符串是不可变的
	//实用修改字符串方式
	sourceStr := "abcdefg"
	//tmpAry := []byte(sourceStr)  使用rune更好
	tmpAry := []rune(sourceStr)
	tmpAry[2] = '我'
	targetStr := string(tmpAry)

	strings.IndexByte(targetStr, 'z') //strings包很有用

	fmt.Println(raw, targetStr)

}

func testArrayAndSlice() {
	const nX = 5
	//arryX := [5]int{1,2,3,4,5} ----- Correct
	//arryX := [5]int{}	------ Correct,可以后面再赋值
	//arryX = [...]int{1,2,3,4,5} ------ Correct
	arryX := [nX]int{1, 2, 3, 4, 5} //这个是定义数组，数组是值传递，在函数调用时参数将发生数据复制，此时nX必须是const定义的，固定长度了
	sliceX := []int{1, 2, 3, 4, 5}  //试一试这么定义，很奇怪的，因为这是定义数组切片，切片是引用传递
	//sliceX := []int ---- Error

	arryY := arryX  //完整拷贝
	arryZ := &arryX //引用

	var px *[5]int = &arryX //数组指针
	a := 1
	b := 2
	pArry := [...]*int{&a, &b} //指针数组
	fmt.Println(px, pArry)

	arryY[2] = 15
	arryZ[2] = 13
	fmt.Println(arryX, arryY, arryZ)
	fmt.Println(len(arryX))

	//遍历数组,切片
	for i, v := range arryX {
		fmt.Println(i, v)
	}

	//同上
	for i := 0; i < len(sliceX); i++ {
		fmt.Println(i, sliceX[i])
	}

	fmt.Println("*******************")

	var tArry = func(arr [5]int) {
		arr[0] = 10
		fmt.Println(arr)
	}

	var tSlice = func(arr []int) {
		arr[0] = 10
		fmt.Println(arr)
	}

	tArry(arryX) //值传递
	//tArry(sliceX)  ----- Error
	//tSlice(&arryX) ----- Error
	//tSlice(arryX)  ----- Error
	tSlice(sliceX) //引用传递

	sliceT := arryX[:3] //对数组取切片
	sliceT = arryX[:]   //将数组转换为切片
	//var sliceT []int
	//var sliceT []float32
	tSlice(sliceT)

	//三种定义切片的方式
	mySlice := make([]int, 5, 10) //5代表元素个数，10代表存储能力（超过存储能力，将自动扩容），初始化为0，
	//mySlice := arryX[:3]
	//mySlice := []int{1,2,3,4,5}



	mySlice = append(mySlice, 8, 9, 10)
	var mySlice1 []int	//mySlice1 is nil
	//mySlice1[1] = 0  ---- Error
	mySlice2 := append(mySlice1, 1,2) //虽然是nil，但可以这样
	mySlice = append(mySlice, mySlice2...) //正确写法，三个点必须
	mySlice = append(mySlice, []int{23, 34}...)

	copy(mySlice, mySlice1) //复制mySlice1给mySlice
	fmt.Println(len(mySlice), cap(mySlice))
}

func testMap() {
	type Person struct {
		name string
		age  int
	}

	p := Person{
		name: "wxy",
		age:  23}
	p.age = 34
	//p.age := 35 --- Error 无论什么情况也不能这么用
	p1 := Person{
		name: "shilf",
		age:  34, //逗号必须有
	}
	//p,p1的区别在与最后一个括号是否另起一行

	mapX := make(map[string]Person)
	mapX["wxy"] = p
	mapX["shilf"] = p1

	var mapY map[string]Person
	mapY = make(map[string]Person, 100) //100代表存储能力
	//注意观察下面整行，其实这里的Person可以不要
	mapY = map[string]Person{"12": Person{"shily", 36}, "13": {"shilp", 45}}

	//两种取值方式
	v := mapX["shilf"]
	v1, ok := mapY["ss"]
	if ok {
		fmt.Println(v, v1)
	} else {
		err := errors.New("NO Data")
		fmt.Println(err)
	}

	//遍历字典
	for k, v := range mapY {
		fmt.Println(k, v)
	}

	//删除key
	delete(mapX, "shilf")
}

func testAnonymous() {
	x := 2
	y := 3

	//普通匿名函数
	max := func(a int, b int) (int, bool) {
		if a > b {
			return a, true
		}
		return b, false
	}
	z, b := max(x, y)
	fmt.Printf("%d,%t\n", z, b)
	fmt.Printf("z type is %T, &z type is %T\n", z, &z)

	manyArg := func(a int, b int, params ...string) (n int, s string, args []string) {
		//...type类似于[]type
		n = 5
		s = "kaka"
		params[0] = "cc"
		for _, v := range params {
			fmt.Println(v)
		}
		args = params
		return
	}

	manyAnyArg := func(c int, params ...interface{}) {
		//任意类型
		fmt.Println(params)
		for _, arg := range params {
			switch arg.(type) {
			case int:
				fmt.Println("int")
			case string:
				fmt.Println("string")
			}
		}
	}

	manyArg(1, 2, "tt", "ss", "gg", "zz")
	manyArg(1, 2, []string{"3", "4", "5"}...) //重要，三个点必须要有

	manyAnyArg(2, []string{"1", "2", "3"}) //重要
	manyAnyArg(2, []int{3, 4, 5})
	manyAnyArg(2, "3", true, 3.12)

	/*func innerFun() string {
		return "innerFun"
	}
	在函数内部定义一般函数是不行的
	*/

	func() string {
		return "innerFun"
	}() //可以这样定义匿名函数，但必须使用
}

func testClosure() {
	//ToDO 闭包待学习研究
	j := 5
	a := func() func() {
		i := 10
		return func() {
			fmt.Printf("Close i , j : %d, %d\n", i, j)
		}
	}
	a()()
}

func testDeferAndPanicAndRecover() {
	//ToDO 待学习研究
}

func testFlag() {
	//CLI功能，cobra包提供

	//-i sss
	var iniStr *string = flag.String("i", "test.ini", "File contains values for sorting")
	//-o ttt
	var outStr *string = flag.String("o", "out.csv", "File to receive sorted values")
	//-n 10
	var numCores = flag.Int("n", 2, "number of Cpu cores to use")

	//runtime运行时包
	runtime.GOMAXPROCS(*numCores)

	//这句很重要，此时才会解析参数的值，在此之前是自己设的默认值
	flag.Parse()

	if iniStr != nil && outStr != nil {
		fmt.Println(*iniStr, *outStr)

		iniFile, err := os.Open(*iniStr)
		if err != nil {
			fmt.Println("Open File Error!")
			return
		}
		defer iniFile.Close()

		//读文件的方式
		br := bufio.NewReader(iniFile)
		byteAry := make([]byte, 0)
		for {
			line, isPrefix, err := br.ReadLine()
			if err == io.EOF {
				break
			}
			fmt.Println(string(line)) //字符数组要转换
			fmt.Println(isPrefix)
			fmt.Println(err)
			byteAry = append(byteAry, line...)
		}
		fmt.Println(string(byteAry))
	}
}

/*----------------TestStruct Start--------------*/

//********************AliasAndDefinition********************

//类型别名----类型别名和原类型完全一样，主要用在其他包类型的Alias
type sInt = int

//类型定义----类型定义和原类型是不同的两个类型，记住是不同类型无法自动转换
type TradeStatus int

func (ts TradeStatus) String() string {
	return tradeStatusSymbol[ts]
}

var tradeStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT"}

//这个是struct中定义的Json解析说明
type Trade struct {
	Tid    int64       `json:"tid"`
	Amount float64     `json:"amount,string"`
	Price  float64     `json:"price,string"`
	Date   int64       `json:"date_ms"`
	Pair   TradeStatus `json:"omitempty"`
}

type IntegerX int

func (a IntegerX) Add(b IntegerX) {
	//值传递
	a += b
}

func (a *IntegerX) AddX(b IntegerX) {
	//引用传递
	*a += b
}

//********************AliasAndDefinition********************

type User struct {
	age    int
	color, //因为跨行了，所以要有逗号，等价于color, name, sex string
	name, sex string
	length   float32
	Person   //匿名结构体字段
	int      //匿名普通字段
	*Company //匿名结构体指针
}

type Person struct {
	country   string
	telephone string
}

type Company struct {
	address string
}

func testStruct() {
	//生产中比较复杂的数据结构，也不过如此
	type ComplexStruct struct {
		WsUrl                 string
		ProxyUrl              string
		ReqHeaders            map[string][]string
		HeartbeatIntervalTime time.Duration
		HeartbeatData         []byte
		ReconnectIntervalTime time.Duration
		ProtoHandleFunc       func([]byte) error
		UnCompressFunc        func([]byte) ([]byte, error)
		ErrorHandleFunc       func(err error)
		IsDump                bool
	}

	//Json转换
	trade := &Trade{123456, 34, 45.6, 155555334344, TradeStatus(2)}
	jsonByteAry, _ := json.Marshal(trade)
	fmt.Println(string(jsonByteAry))

	var x IntegerX = 1
	var y IntegerX = 2
	x.Add(y)
	x.AddX(2)

	z := 1
	//x.Add(z) ----- Error，z的类型无法匹配
	fmt.Println(x, y, z)

	//结构体对象指针
	var u *User = NewUserX(12, "yellow", "shilf", "male", 1.70)
	u.age = 14
	//u := new(User)
	//u := &User{age:12,name:"shilp"}
	fmt.Printf("u Type : %T \n", u)
	//打印结构体的三种方式
	fmt.Printf("u : %v\n", u)
	fmt.Printf("u+ : %+v\n", u)
	fmt.Printf("u# : %#v\n", u)

	//结构体对象
	u1 := NewUser(10, "shilp", "male", 1.74)
	u1.age = 13
	fmt.Println("u1 Type : ", reflect.TypeOf(u1))

	//以下都可以
	u.print()
	u.printX()
	u1.print()
	u1.printX()

	//以下必须对应
	printUserX(u)
	printUser(u1)
	//printUser(u)	----- Error, 此时必须对应

}

func (u User) print() {
	//值传递
	u.age = 0
	fmt.Println(u)
}

func (u *User) printX() {
	//引用传递
	u.age = 0
	fmt.Println(u)
}

func printUser(u User) {
	fmt.Println(u.age, u.name, u.sex, u.length, u.country)
}

func printUserX(u *User) {
	fmt.Println(u.age, u.name, u.sex, u.length, u.country)
	//u.country等价于u.Person.country,也可以那么写
}

func NewUser(age int, name string, sex string, length float32) User {
	return User{age: age, name: name, sex: sex, length: length}
}

func NewUserX(age int, color string, name string, sex string, length float32) *User {
	//这种方式必须写全
	return &User{age, color, name, sex, length, Person{"china", "18601637777"}, 0, &Company{"pudong"}}
}

/*----------------TestStruct End--------------*/

/*----------------TestInterFace Start-------------*/

type ICar interface {
	/*
		1.接口是一组方法的集合，但不包含方法的实现
		2.接口中也不能包含变量
	*/
	Driver() string
	Run() string
}

type IBus interface {
	Close() string
}

type IAllChe interface {
	//组合接口
	ICar
	IBus
}

type Benz struct {
	producer string
	speed    float32
}

type Leno struct {
	name string
	age  int
}

type Mazda struct {
	country string
	color   string
}

func (car Benz) Driver() string {
	return strconv.FormatFloat(float64(car.speed), 'f', 6, 64)
}

func (car Benz) Run() string {
	return car.producer + "--run"
}

func (car Leno) Driver() string {
	return car.name + "--drive"
}

func (car Leno) Run() string {
	return "run--" + strconv.Itoa(car.age)
}

func (car *Mazda) Driver() string {
	return car.country + "--driver"
}

func (car *Mazda) Run() string {
	return car.color + "--run"
}

func testInterFace() {
	testCar := func(car ICar) {
	}

	var bcar ICar = Benz{"MEISAIDES", 250.4}
	var lcar ICar = Leno{"hero", 2}
	var mcar ICar = &Mazda{"japan", "red"}

	testCar(bcar)
	testCar(lcar)
	testCar(Benz{"Germany", 250})
	testCar(&Benz{"Germany", 250}) //这个却可以，方法集的问题[结构体指针的方法集包含结构体对象的方法集]
	testCar(mcar)
	//testCar(Mazda{}) --- Error

	var tmpCar *Benz
	if tmpCar == nil {
		//走这里
		fmt.Println("tmpCar is nil")
	} else {
		fmt.Println("tmpCar is not nil")
	}

	var iCar ICar = tmpCar
	if iCar == nil {
		fmt.Println("iCar is nil")
	} else {
		//走这里哦
		//原因：当且仅当接口的动态值和动态类型都为 nil 时，接口类型值才为 nil
		//此时iCar的动态值为nil，但是动态类型却是Benz
		fmt.Println("iCar is not nil")
	}
	fmt.Printf("iCar is nil pointer : %t", iCar == (*Benz)(nil))

	/* 类似
	任何类型都实现了空接口
	var x interface{}
	var v *User
	x = v
	fmt.Printf("x is nil pointer : %t", x == (*User)(nil))
	*/
}

/*----------------TestInterFace End-------------*/

func testAssert() {
	var bcar ICar = Benz{}
	a := bcar.(Benz)
	fmt.Printf("%T\n", a)

	Assert := func(i interface{}) (value int, ok bool) {
		value, ok = i.(int)
		return
	}
	//bcar.(Benz)  i.(int) 必须是接口才能这么用

	/*
		比如int，float32这些都是静态类型，接口既有静态类型又动态类型
		接口的静态类型就是接口本身，接口没有静态值
		接口的动态类型就是实现接口的类型，该类型的对象就是接口的动态值
	*/
	FindType := func(i interface{}) {
		switch i.(type) {
		case string:
			fmt.Printf("string : %s\n", i.(string))
		case int:
			fmt.Printf("int : %d\n", i.(int))
		case bool:
			fmt.Printf("bool : %t\n", i.(bool))
		case Benz:
			fmt.Printf("BenZ : %+v\n", i.(Benz))
		default:
			fmt.Printf("UnKnown type\n")
		}
		//i必须是接口
		// i.(type)就是获取接口的动态类型
		// 必须在switch中使用
	}

	FindCar := func(i ICar) {
		switch i.(type) {
		case Benz:
			fmt.Printf("BenZ : %+v\n", i.(Benz))
		case Leno:
			fmt.Printf("Leno : %+v\n", i.(Leno))
		}
	}

	FindType(bcar)
	FindCar(bcar)
	FindCar(Leno{})

	if value, ok := Assert(2); ok {
		fmt.Println("value is ", value)
	}
}

func testOs() {
	envPath := os.Getenv("Path")
	fmt.Println(envPath)
	args := os.Args
	fmt.Println(args)
	if len(args) > 2 {
		arg, _ := strconv.Atoi(args[1])
		fmt.Println(arg)
		fmt.Println("the length of args is bigger than 2")
	}
}

func testGoroutineAndChan() {
	var ch chan int
	var m map[string]chan bool
	ch = make(chan int)
	chx := make(chan int, 10)
	m = make(map[string]chan bool)
	m["lala"] = make(chan bool)
	fmt.Println(ch, chx, m)

	Count := func(ch chan int, i int) {
		fmt.Println("Counting")
		time.Sleep(time.Duration(5) * time.Second)
		ch <- i
	}

	chs := make([]chan int, 10)
	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		go Count(chs[i], i)
	}

	for _, ch := range chs {
		fmt.Println("zzzzzzzzzz")
		fmt.Println(<-ch)
	}
}

func test() {
	testVariables()
	testEnum()
	testKind()
	testStr()
	testArrayAndSlice()
	testMap()
	testAnonymous()
	testFlag()
	testClosure()
	testDeferAndPanicAndRecover()
	testStruct()
	testInterFace()
	testAssert()
	testOs()
	testGoroutineAndChan()

	testRedis()
	testMux()
	testGorillaWebSocket()
	testDecimal()
	testCast()
	testCobra()
	testPretty()
	testViper()
	testJSON()
}

//函数外定义变量
//helloWorld := "Hello World !" ---- Error 简短声明的变量只能在函数内部使用
var helloWorld = "Hello World !"

//var helloWorld string = "Hello World !" ---- Correct

func main() {
	fmt.Println(helloWorld)
	/*
		------------------Go 语言特点------------------
		1.很少使用（）,左大括号不能另起一行
		2.结束不使用";"冒号
		3.go的基础结构比较少，数组（不常用），切片（常用），字典，结构，接口
		4.a int类型在变量名后面
		5.要生成可执行程序，必须要有一个名字为main的包，并且在该包中包含main()函数
	*/


	//ch := make(chan int)
	/*for i := 0; i < 10; i++ {
		select {
		case ch <- 0:
			fmt.Println("-------aa")
		case ch <- 1:
			fmt.Println("-------bb")
		case ch <- 2:
			fmt.Println("-------cc")
		case ch <- 3:
			fmt.Println("-------dd")
		}
		vaule := <-ch
		fmt.Println(vaule)
	}*/
}
