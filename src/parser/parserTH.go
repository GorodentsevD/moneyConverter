package parser

import (
	"fmt"
	"strings"
)

type THData struct {
	Descriptions []Description `xml:"item"`
	ValuteList []Valute
}


type Description struct {
	Desc string `xml:"title"`
}


func (d *THData) Parse() {

	var str string
	var tokens = make([][]string, len(d.Descriptions))

	for i := 0; i < len(d.Descriptions); i++ {
		str = d.Descriptions[i].Desc
		tokens[i] = make([]string, 12)
		tokens[i] = strings.Split(str, " ")
	}

	for i := 1; i < len(d.Descriptions); i++ {
		var valute Valute

		valute.CharCode = tokens[i][5]
		valute.Nominal = tokens[i][4]

		if tokens[i][10] == "Average" {
			if tokens[i][11] == "Buying" && tokens[i][12] == "Sight" {
				valute.BuyRate = tokens[i][1]
			}
			if tokens[i + 2][5] == valute.CharCode && (tokens[i + 2][11] == "Selling") {
				valute.SellRate = tokens[i + 2][1]
				i += 2
			}
		} else {
			if tokens[i][10] == "Buying" {
				valute.BuyRate = tokens[i][1]
			}
			if tokens[i + 1][10] == "Selling" {
				valute.SellRate = tokens[i][1]
				i += 1
			} else {
				i += 1
				continue
			}
		}
		d.ValuteList = append(d.ValuteList, valute)
	}
}

func ShowValutes(valuteList []Valute) {

	for i := 0; i < len(valuteList); i++ {
		fmt.Printf("CHarCode = %s, Nominal = %s, SellRate = %s, BuyRate = %s\n",
			valuteList[i].CharCode,
			valuteList[i].Nominal,
			valuteList[i].SellRate,
			valuteList[i].BuyRate,
			)
	}
}


func (d THData) ShowCourses() string{

	var str string
	for i := 0; i < len(d.Descriptions); i++ {
		str += fmt.Sprintf("%s\n",
			d.Descriptions[i].Desc,
		)
	}
	return str
}