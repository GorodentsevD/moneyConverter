package parser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type CBRData struct {
	ValuteList []Valute `xml:"Valute"`
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

func (d CBRData) Parse() {
	for i := 0; i < len(d.ValuteList); i++ {
		d.ValuteList[i].SellRate = strings.Replace(d.ValuteList[i].SellRate, ",", ".", 1)
		d.ValuteList[i].Source = "CBR"
		d.ValuteList[i].BuyRate = d.ValuteList[i].SellRate
	}
}

func LoadFromSource(url string) []byte {
	req, err := http.Get(url)
	defer req.Body.Close()

	if err != nil {
		panic(err)
	}

	b, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	return b
}
