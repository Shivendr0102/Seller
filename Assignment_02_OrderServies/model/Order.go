package model

type Order struct {
	ID           string  `json:"id"`
	Status       string  `json:"status"`
	Items        []Items `json:"items"`
	Total        string  `json:"total"`
	CurrencyUnit string  `json:"currencyUnit"`
}

type Items struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Quantity    string `json:"quantity"`
}
