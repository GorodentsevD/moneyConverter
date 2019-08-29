package server

import (
	"../db"
	"../parser"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"html/template"
	"net/http"
	"strconv"
)

// структура для данных, полученных от клиента
type ClientData struct {
	Valute1 string
	Valute2 string
	Value   string
	Source  string
}

// структура для хранения данных, полученных из бд
type FormattedValute struct {
	Nominal  float64
	SellRate float64
	BuyRate  float64
}

// запуск приложения
func StartServer(d parser.Parser, dataBase sqlx.DB, tableName string) {

	var clientQuery ClientData
	val := db.GetAllCharCodes(dataBase, tableName)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			panic(err)
		}

		err = t.Execute(w, &val)
		if err != nil {
			panic(err)
		}
	})
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			clientQuery = ClientData{
				Valute1: r.FormValue("Valute1"),
				Valute2: r.FormValue("Valute2"),
				Value:   r.FormValue("Value"),
				Source:  r.FormValue("Source"),
			}

			res, err := json.Marshal(getResult(clientQuery, dataBase, tableName))
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	})

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		panic(err)
	}
}

// получить результат конвертирования
func getResult(data ClientData, dataBase sqlx.DB, tableName string) float64 {

	var valute1 parser.Valute
	var valute2 parser.Valute

	var formattedValute1 FormattedValute
	var formattedValute2 FormattedValute

	var err error

	valute1 = db.GetValuteByCharCode(dataBase, data.Valute1, tableName, data.Source)
	valute2 = db.GetValuteByCharCode(dataBase, data.Valute2, tableName, data.Source)

	formattedValute1.Nominal, err = strconv.ParseFloat(valute1.Nominal, 64)
	if err != nil {
		panic(err)
	}

	formattedValute2.Nominal, err = strconv.ParseFloat(valute2.Nominal, 64)
	if err != nil {
		panic(err)
	}

	formattedValute1.SellRate, err = strconv.ParseFloat(valute1.SellRate, 64)
	if err != nil {
		panic(err)
	}

	formattedValute2.BuyRate, err = strconv.ParseFloat(valute2.BuyRate, 64)
	if err != nil {
		panic(err)
	}

	value, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("\n***VALUES***\n\nvalue = %f\n\nVALUE 1: \nSell = %f\nNominal = %f\n\nVALUTE 2:\nSell = %f\nNominal = %f",
	//	value, formattedValute1.SellRate, formattedValute1.Nominal, formattedValute2.SellRate, formattedValute2.Nominal)

	result := value * (formattedValute1.SellRate / formattedValute1.Nominal) /
		(formattedValute2.BuyRate / formattedValute2.Nominal)

	return result
}
