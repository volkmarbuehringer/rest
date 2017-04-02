package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable application_name=gogogo connect_timeout=3",
		os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"), os.Getenv("PGHOST"), os.Getenv("PGPORT"))
	DB, err = sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	//	fmt.Println("stats", DB.Stats())

}

func ParseLines(filePath string, parse func(string) error) (err error) {

	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	for i := 0; scanner.Scan(); i++ {
		if i > 0 {
			if err = parse(scanner.Text()); err != nil {
				//		results = append(results, output)
				return
			}
		}

	}
	if err = scanner.Err(); err != nil {
		return
	}

	return
}

type tf struct {
	from int
	to   int
}

type t struct {
	Year   int `db:"year" file:"0:4"`
	Anzsic string
	Area   string `db:"area" file:"5:5"`
	Geo    int    `db:"geo" file:"0:4"`
	Ec     int    `db:"ec" file:"0:4"`
}

func columns(r interface{}) (cols []string, cols1 []string, lpfun func(s string) (res []interface{}, err error)) {
	//var r t
	cols = make([]string, 0)
	var flag bool
	tt := reflect.TypeOf(r)
	pc := make([]int, tt.NumField())
	pcf := make([]tf, tt.NumField())
	for i := 0; i < tt.NumField(); i++ {
		pc[i] = -1

		//tag := string(tt.Field(i).Tag)

		if alias, ok := tt.Field(i).Tag.Lookup("db"); ok {

			pc[i] = i
			//	pre := strings.Replace(v[4:], "\"", "", 1)
			cols = append(cols, alias)

			switch reflect.ValueOf(r).Field(i).Interface().(type) {
			case int, int64:
				cols1 = append(cols1, alias+" bigint")
			case string:
				cols1 = append(cols1, alias+" text")
			}
		}

		if alias, ok := tt.Field(i).Tag.Lookup("file"); ok {
			//pre := strings.Replace(v[6:], "\"", "", 1)
			x := strings.Split(alias, ":")
			x1, _ := strconv.Atoi(x[0])
			x2, _ := strconv.Atoi(x[1])
			pcf[i] = tf{x1, x2}
			flag = true
		}

		fmt.Printf("%+v  %+v %+v \n", tt.Field(i), reflect.TypeOf(tt.Field(i)), cols)
	}

	innerfun := func(x string, i int, res *[]interface{}) (err error) {

		gaga := reflect.ValueOf(r).Field(i)
		var ga interface{}
		switch gaga.Interface().(type) {
		case int:
			ga, err = strconv.Atoi(x)
		//	gaga.Set(reflect.ValueOf(77))
		//item.Elem().FieldByName(name).Set(reflect.ValueOf(lala))
		case string:
			ga = x
		//	gaga.Set(reflect.ValueOf("lala"))
		case float64:
			//gaga.Set(reflect.ValueOf(33.5))
			ga, err = strconv.ParseFloat(x, 64)
		case time.Time:
			gaga.Set(reflect.ValueOf(time.Now()))
		default:
			fmt.Print("hier anders")
		}
		*res = append(*res, ga)
		return
	}

	if flag {
		lpfun = func(s string) (res []interface{}, err error) {

			if len(s) > 100 || len(s) < 10 {
				err = fmt.Errorf("fehler %s", s)
				return
			}
			res = make([]interface{}, 0)

			for i := 0; i < tt.NumField() && err == nil; i++ {
				//			fmt.Println(i, pc, tt)
				if pcf[i].to > 0 {

					x := s[pcf[i].from:pcf[i].to]
					//	fmt.Println("hier da", s, pcf[i], i, x)

					innerfun(x, i, &res)

				}

			}
			//fmt.Println("hier fertig", r)
			return
		}

	} else {
		lpfun = func(s string) (res []interface{}, err error) {

			//	itemTyp := reflect.TypeOf(&r).Elem()
			//	item := reflect.New(itemTyp)
			x := strings.Split(s, ",")

			if len(x) < tt.NumField() || len(s) > 100 || len(s) < 10 {
				err = fmt.Errorf("fehler %s", s)
				return
			}
			res = make([]interface{}, 0)

			for i := 0; i < tt.NumField() && err == nil; i++ {
				//			fmt.Println(i, pc, tt)
				if pc[i] >= 0 {
					innerfun(x[pc[i]], i, &res)

				}

			}
			//fmt.Println("hier fertig", r)
			return
		}
	}

	return

}

func packer(s string) (res []interface{}, err error) {
	x := strings.Split(s, ",")
	res = make([]interface{}, len(x))
	for i, v := range x {
		//	fmt.Println(i, v)
		switch i {

		case 0, 3, 4:
			res[i], err = strconv.Atoi(v)
		default:
			res[i] = v
		}

	}
	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: line_parser <path>")
		return
	}

	txn, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	cols, _, packer1 := columns(t{})

	//if _, err1 := txn.Exec(`create temporary table csvtest ( ` + strings.Join(cols1, ",") + ` ) on commit preserve rows`); err != nil {
	//	log.Fatal(err1)
	//}

	stmt, err := txn.Prepare(pq.CopyIn("csvtest", cols...))
	if err != nil {
		log.Fatal(err)
	}

	fun := func(s string) (err error) {
		var xx []interface{}
		if xx, err = packer1(s); err != nil {
			log.Fatal(err)
		}

		if _, err = stmt.Exec(xx...); err != nil {
			return
		}
		return
	}

	if err = ParseLines(os.Args[1], fun); err != nil {
		fmt.Println("Error while parsing file", err)
		return
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal(err)
	}

}
