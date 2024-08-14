package main

import (
	"fmt"
	"log/slog"
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
		slog.Error("ParseMessageForPrint: err: %v", err)
		return
	}
	printerName, err := FindPrinterName()
	if err != nil {
		slog.Error("FindPrinterName: err: %v", err)
		return
	}
	if printerName == "" {
		slog.Error("FindPrinterName: msg: printer not found")
		return
	}

	printer := Printer{Name: printerName, Remote: false, RawFile: fmtStr}
	if err := printer.Do(); err != nil {
		slog.Error("Do: err: %v", err)
		return
	}

	fmt.Println("success")
}
