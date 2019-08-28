package main

import (
	. "./db"
	. "./parser"
	"bytes"
	"encoding/xml"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/html/charset"

	. "./server"
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

	var cbrData CBRData
	var CBRXML = LoadFromSource(CBRSource)

	var thData THData
	var THXML = LoadFromSource(THSource)

	var dataBase sqlx.DB

	// *****
	// данный кусок кода пошел на замену xml.UnMarshall
	// т.к. в источнике используется кодировка, отличная от UTF-8

	// считывание xml от ЦБР
	reader := bytes.NewReader(CBRXML)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err := decoder.Decode(&cbrData)
	if err != nil {
		panic(err)
	}

	// считывание xml от Банка Тайланда
	reader = bytes.NewReader(THXML)
	decoder = xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&thData)
	if err != nil {
		panic(err)
	}
	// *****

	cbrData.Parse()
	thData.Parse()

	cbrData.ShowCourses()
	thData.ShowCourses()

	dataBase = ConnectToDB(Driver, DbName, User, Password)
	CreateTable(&dataBase, Table)

	AddToTable(&dataBase, Table, thData.ValuteList)
	AddToTable(&dataBase, Table, cbrData.ValuteList)

	// запуск веб-сервера
	StartServer(cbrData, dataBase, Table)

}
