package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
)

// структура для данных банка Тайланда
type THData struct {
	Descriptions []Description `xml:"item"`
	ValuteList   []Valute      // список валют
	Source       string        // url источника
	XML          []byte        // xml источника
}

// структура для хранения тайтла(именно это поле будет парситься, тк в нем содержаться все нужные данные)
type Description struct {
	Desc string `xml:"title"`
}

// функция парсинга данных
func (data *THData) Parse() {
	// очистка списка валют
	data.ValuteList = nil

	var valuteList []Valute = nil
	var str string
	var tokens = make([][]string, len(data.Descriptions))

	for i := 0; i < len(data.Descriptions); i++ { // разбитие заголовка в xml на токены
		str = data.Descriptions[i].Desc
		tokens[i] = make([]string, 12)
		tokens[i] = strings.Split(str, " ")
	}

	for i := 1; i < len(data.Descriptions); i++ { // цикл прохода по всем данным
		var valute Valute

		valute.CharCode = tokens[i][5]
		valute.Nominal = tokens[i][4]

		// если в строке на 10 позиции слово average
		if tokens[i][10] == "Average" {
			// если в строке есть выражение Buying Sight, то загружаем ставку покупки
			if tokens[i][11] == "Buying" && tokens[i][12] == "Sight" {
				valute.BuyRate = tokens[i][1]
			}
			// если в строке есть слово Selling, то загружаем ставку продажи
			if tokens[i+2][5] == valute.CharCode && (tokens[i+2][11] == "Selling") {
				valute.SellRate = tokens[i+2][1]
				i += 2
			}
		} else { // иначе прочитали не average ставку, поэтому токены находятся на других позициях
			if tokens[i][10] == "Buying" {
				valute.BuyRate = tokens[i][1]
			}
			if tokens[i+1][10] == "Selling" {
				valute.SellRate = tokens[i][1]
				i += 1
			} else {
				i += 1
				continue
			}
		}
		valute.Source = "THB" // объявляем источник для валюты
		valuteList = append(valuteList, valute)
	}
	data.ValuteList = valuteList // загружаем список валют в структуру
}

// функция вывода курсов валют из структуры
func (data THData) ShowCourses() {

	var str string
	for i := 0; i < len(data.ValuteList); i++ {
		str += fmt.Sprintf("%s, %s, %s, %s. %s\n",
			data.ValuteList[i].CharCode,
			data.ValuteList[i].Nominal,
			data.ValuteList[i].SellRate,
			data.ValuteList[i].BuyRate,
			data.ValuteList[i].Source,
		)
	}
	fmt.Println(str)
}

// функция загрузки данных из источника
func (data *THData) LoadFromSource() {
	data.Descriptions = nil

	fmt.Println("HTTP GET to TH bank...")
	req, err := http.Get(data.Source)
	defer req.Body.Close()
	if err != nil {
		panic(err)
	}

	data.XML, err = ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(data.XML)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&data)
	if err != nil {
		panic(err)
	}
}

// функция получения списка валют
func (d THData) GetValuteList() []Valute {
	return d.ValuteList
}
