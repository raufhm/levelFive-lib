package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

func main() {
	ps := NewParser()
	ro := RootObject{
		TicketNo:       "1000",
		Date:           "12/08/2024",
		Time:           "10:15",
		Entities:       "test entities",
		Orders:         "test orders",
		Discounts:      "test discounts",
		TicketTotal:    "test total ticket",
		PaymentDetails: "test payment details",
	}
	fmtStr, err := ps.ParseMessageForPrint(TemplateParseString2, ro)
	if err != nil {
		log.Error("ParseMessageForPrint: err: " + err.Error())
		return
	}
	printerName, err := FindPrinterName()
	if err != nil {
		log.Error("FindPrinterName: err: " + err.Error())
		return
	}
	if printerName == "" {
		log.Error("FindPrinterName: msg: printer not found")
		return
	}

	printer := Printer{Name: printerName, Remote: false, RawFile: fmtStr}
	if err := printer.Do(); err != nil {
		log.Error("Do: err: " + err.Error())
		return
	}

	fmt.Println("success")
}
