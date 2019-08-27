package db

import (
	"fmt"

	. "../parser"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

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

	conn.MustExec("CREATE DATABASE IF NOT EXISTS "+dbName)
	conn.MustExec("USE "+dbName)

	return *conn
}

func AddToTable(conn *sqlx.DB, tableName string, ValuteList []Valute) {

	err := conn.Ping()
	if err != nil {
		panic(err)
	}

	drop :=  fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	conn.MustExec(drop)


	schema := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (" +
		"`CharCode` VARCHAR(3) NOT NULL PRIMARY KEY, " +
		"`Nominal` INT NOT NULL, " +
		"`SellRate` DEC(10,4) NOT NULL, " +
		"`BuyRate` DEC(10,4) NOT NULL)", tableName)

	conn.MustExec(schema) // создание таблицы по шаблону schema

	insertIntoDB := fmt.Sprintf("INSERT INTO %s VALUES", tableName)

	// запись в бд
	for i := 0; i < len(ValuteList); i++ {
		insertIntoDB += fmt.Sprintf(" (\"%s\", %s, %s, %s)",
			ValuteList[i].CharCode,
			ValuteList[i].Nominal,
			ValuteList[i].SellRate,
			ValuteList[i].SellRate,
		)
		if i != len(ValuteList) - 1 {
			insertIntoDB += ","
		}
	}

	_, err = conn.Exec(insertIntoDB)
	if err != nil {
		panic(err)
	}
}

func GetValuteByCharCode(conn sqlx.DB, CharCode string , tableName string) Valute {

	var valute Valute

	query := fmt.Sprintf("SELECT * FROM %s WHERE CharCode=\"%s\"", tableName, CharCode)

	err := conn.Get(&valute, query)
	if err != nil {
		panic(err)
	}

	return valute
}

func GetAllCharCodes(conn sqlx.DB, tableName string) []string {

	query := fmt.Sprintf("SELECT CharCode FROM %s", tableName)

	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}

	var CharCodeList []string
	var tmp string
	for rows.Next() {
		rows.Scan(&tmp)
		CharCodeList = append(CharCodeList, tmp)
	}

	return CharCodeList
}

