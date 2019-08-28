package parser

type Valute struct {
	CharCode string `xml:"CharCode" json:"CharCode" db:"CharCode"` // буквенное значение валюты
	Nominal  string `xml:"Nominal" json:"Nominal" db:"Nominal"`    // номинал валюты
	SellRate string `xml:"Value" json:"Value" db:"SellRate"`       // стоимость валюты относительно рубля
	BuyRate  string `db:"BuyRate"`
	Source   string `db:"Source"`
}

type Parser interface {
	Parse()
	ShowCourses()
}
