package main

import (
	. "./db"
	. "./parser"
	. "./server"
	"github.com/jmoiron/sqlx"
	"time"
)

const CBRSource = "http://www.cbr.ru/scripts/XML_daily.asp"
const THSource = "https://www.bot.or.th/App/RSS/fxrate-all.xml"

const Driver = "mysql"

const User = "root"
const Password = "00000000"

const DbName = "valutes"
const Table = "valutes"

func main() {

	var cbrData = &CBRData{Source: CBRSource}
	var thData = &THData{Source: THSource}
	var dataBase sqlx.DB

	cbrData.LoadFromSource()
	thData.LoadFromSource()

	cbrData.Parse()
	thData.Parse()

	dataBase = ConnectToDB(Driver, DbName, User, Password)
	CreateTable(&dataBase, Table)

	AddToTable(&dataBase, Table, thData.ValuteList)
	AddToTable(&dataBase, Table, cbrData.ValuteList)

	// горутина для обновления базы данных
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for range ticker.C {
			RefreshTable(&dataBase, Table, cbrData, thData)
		}
	}()

	// запуск веб-сервера
	StartServer(cbrData, dataBase, Table)

}
