package parser

type Valute struct {
	CharCode string `xml:"CharCode" json:"CharCode" db:"CharCode"` // буквенное значение валюты
	Nominal  string `xml:"Nominal" json:"Nominal" db:"Nominal"`    // номинал валюты
	SellRate string `xml:"Value" json:"Value" db:"SellRate"`       // ставка продажи
	BuyRate  string `db:"BuyRate"`                                 // ставка покупки
	Source   string `db:"Source"`                                  // источник валюты
}

type Parser interface {
	Parse()
	ShowCourses()
	LoadFromSource()
	GetValuteList() []Valute
}
