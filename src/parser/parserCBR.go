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

type CBRData struct {
	ValuteList []Valute `xml:"Valute"`
	Source     string
	XML        []byte
}

func (d CBRData) ShowCourses() {

	var str string
	for i := 0; i < len(d.ValuteList); i++ {
		str += fmt.Sprintf("%3s	%6s	%7s %s %s\n",
			d.ValuteList[i].CharCode,
			d.ValuteList[i].Nominal,
			d.ValuteList[i].SellRate,
			d.ValuteList[i].BuyRate,
			d.ValuteList[i].Source,
		)
	}
	fmt.Println(str)
}

func (d *CBRData) Parse() {
	for i := 0; i < len(d.ValuteList); i++ {
		d.ValuteList[i].SellRate = strings.Replace(d.ValuteList[i].SellRate, ",", ".", 1)
		d.ValuteList[i].Source = "CBR"
		d.ValuteList[i].BuyRate = d.ValuteList[i].SellRate
	}
}

func (data *CBRData) LoadFromSource() {

	data.ValuteList = nil

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

func (data CBRData) GetValuteList() []Valute {
	return data.ValuteList
}
