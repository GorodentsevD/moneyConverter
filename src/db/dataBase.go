package db

import (
	"fmt"

	. "../parser"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type ValuteCode struct {
	CharCode string `db:"CharCode" json:"CharCode"`
	Source   string `db:"Source" json:"Source"`
}

func ConnectToDB(driver string, dbName string, user string, password string) sqlx.DB {

	connectString := fmt.Sprintf("%s:%s@tcp(localhost:3306)/", user, password)

	conn, err := sqlx.Connect(driver, connectString)
	if err != nil {
		panic(err)
	}

	err = conn.Ping()
	if err != nil {
		panic(err)
	}

	conn.MustExec("CREATE DATABASE IF NOT EXISTS " + dbName)
	conn.MustExec("USE " + dbName)

	return *conn
}

func CreateTable(conn *sqlx.DB, tableName string) {
	err := conn.Ping()
	if err != nil {
		panic(err)
	}

	drop := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	conn.MustExec(drop)

	schema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ("+
		"`CharCode` VARCHAR(3) NOT NULL, "+
		"`Nominal` INT NOT NULL, "+
		"`SellRate` DEC(10,4) NOT NULL, "+
		"`BuyRate` DEC(10,4) NOT NULL,"+
		"`Source` VARCHAR(5) NOT NULL)", tableName)

	conn.MustExec(schema) // создание таблицы по шаблону schema
}

func AddToTable(conn *sqlx.DB, tableName string, ValuteList []Valute) {
	err := conn.Ping()
	if err != nil {
		panic(err)
	}

	insertIntoDB := fmt.Sprintf("INSERT INTO %s VALUES", tableName)

	// запись в бд
	for i := 0; i < len(ValuteList); i++ {
		insertIntoDB += fmt.Sprintf(" (\"%s\", %s, %s, %s, \"%s\")",
			ValuteList[i].CharCode,
			ValuteList[i].Nominal,
			ValuteList[i].SellRate,
			ValuteList[i].BuyRate,
			ValuteList[i].Source,
		)
		if i != len(ValuteList)-1 {
			insertIntoDB += ","
		}
	}

	fmt.Println(ValuteList)
	_, err = conn.Exec(insertIntoDB)
	if err != nil {
		panic(err)
	}
}

func GetValuteByCharCode(conn sqlx.DB, CharCode string, tableName string, source string) Valute {
	var valute Valute

	query := fmt.Sprintf("SELECT * FROM %s WHERE CharCode=\"%s\" AND Source=\"%s\"", tableName, CharCode, source)
	err := conn.Get(&valute, query)
	if err != nil {
		panic(err)
	}

	return valute
}

func GetAllCharCodes(conn sqlx.DB, tableName string) []ValuteCode {

	query := fmt.Sprintf("SELECT CharCode, Source FROM %s", tableName)
	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	var valuteList []ValuteCode
	var valute ValuteCode
	for rows.Next() {
		err = rows.Scan(&valute.CharCode, &valute.Source)
		if err != nil {
			panic(err)
		}
		valuteList = append(valuteList, valute)
	}

	return valuteList
}

func RefreshTable(conn *sqlx.DB, tableName string, data ...Parser) {
	fmt.Println("Refresh is start")
	err := conn.Ping()
	if err != nil {
		panic(err)
	}

	// используем данную функцию для сброса таблицы
	CreateTable(conn, tableName)

	var valuteList []Valute
	for i, singleData := range data {
		fmt.Printf("Refresh #%d\n", i)
		singleData.LoadFromSource()
		singleData.Parse()
		valuteList = singleData.GetValuteList()
		fmt.Println(valuteList)
		AddToTable(conn, tableName, valuteList)
	}

	fmt.Println("Refresh is done")
}
