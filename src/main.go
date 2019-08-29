package main

import (
	. "./db"
	. "./parser"
	. "./server"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const CBRSource = "http://www.cbr.ru/scripts/XML_daily.asp"
const THSource = "https://www.bot.or.th/App/RSS/fxrate-all.xml"

const DbName = "valutes"

const CBRTable = "cbrValutes"
const THTABLE = "thValutes"

const Table = "valutes"

const Driver = "mysql"

const User = "root"
const Password = "00000000"

func main() {

	var cbrData = &CBRData{Source: CBRSource}
	var thData = &THData{Source: THSource}

	cbrData.LoadFromSource()
	thData.LoadFromSource()

	var dataBase sqlx.DB

	cbrData.Parse()
	thData.Parse()

	cbrData.ShowCourses()
	thData.ShowCourses()

	dataBase = ConnectToDB(Driver, DbName, User, Password)
	CreateTable(&dataBase, Table)

	AddToTable(&dataBase, Table, thData.ValuteList)
	AddToTable(&dataBase, Table, cbrData.ValuteList)

	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			fmt.Println("goroutine #")
			RefreshTable(&dataBase, Table, cbrData, thData)
		}
	}()

	// запуск веб-сервера
	StartServer(cbrData, dataBase, Table)

}
