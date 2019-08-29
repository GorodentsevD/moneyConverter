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

type THData struct {
	Descriptions []Description `xml:"item"`
	ValuteList   []Valute
	Source       string
	XML          []byte
}

type Description struct {
	Desc string `xml:"title"`
}

func (data *THData) Parse() {

	data.ValuteList = nil

	var valuteList []Valute = nil
	var str string
	var tokens = make([][]string, len(data.Descriptions))

	for i := 0; i < len(data.Descriptions); i++ {
		str = data.Descriptions[i].Desc
		tokens[i] = make([]string, 12)
		tokens[i] = strings.Split(str, " ")
	}

	for i := 1; i < len(data.Descriptions); i++ {
		var valute Valute

		valute.CharCode = tokens[i][5]
		valute.Nominal = tokens[i][4]

		if tokens[i][10] == "Average" {
			if tokens[i][11] == "Buying" && tokens[i][12] == "Sight" {
				valute.BuyRate = tokens[i][1]
			}
			if tokens[i+2][5] == valute.CharCode && (tokens[i+2][11] == "Selling") {
				valute.SellRate = tokens[i+2][1]
				i += 2
			}
		} else {
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
		valute.Source = "THB"
		valuteList = append(valuteList, valute)
	}
	data.ValuteList = valuteList
}

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

func (data *THData) LoadFromSource() {
	data.Descriptions = nil

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

func (d THData) GetValuteList() []Valute {
	return d.ValuteList
}
