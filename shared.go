package main

import "time"

//const DefaultPrinterName = "Test_Printer_1"

type Ticket struct {
	Terminal    string
	LoginUser   string
	Date        time.Time
	Time        string
	PaymentDate time.Time
	PaymentTime string
	TagPax      string
	Orders      []Order
	Discounts   []Discount
	Services    []Service
	Taxes       []Tax
	Payments    []Payment
}

type Order struct {
	Name     string
	Quantity float64
	Price    float64
}

type Discount struct {
	Name   string
	Amount float64
}

type Service struct {
	Name   string
	Amount float64
}

type Tax struct {
	Name   string
	Amount float64
}

type Payment struct {
	Name               string
	Tendered           float64
	PaymentInformation PaymentInfo
}

type PaymentInfo struct {
	RefNo string
}

type RootObject struct {
	TicketNo       string
	Date           string
	Time           string
	Entities       string
	Orders         string
	Discounts      string
	TicketTotal    string
	PaymentDetails string
}
