package main

import (
	"github.com/raufhm/levelfive-lib/parsing"
	"github.com/raufhm/levelfive-lib/printing"
	"github.com/raufhm/levelfive-lib/shared"
	"log/slog"
)

func main() {
	ps := parsing.NewParser()
	ro := shared.RootObject{
		TicketNo:       "1000",
		Date:           "12/08/2024",
		Time:           "10:15",
		Entities:       "test entities",
		Orders:         "test orders",
		Discounts:      "test discounts",
		TicketTotal:    "test total ticket",
		PaymentDetails: "test payment details",
	}
	fmtStr, err := ps.ParseMessageForPrint(shared.TemplateParseString2, ro)
	if err != nil {
		slog.Error("ParseMessageForPrint: err: %v", err)
		return
	}

	pr, err := printing.NewPrinter(printing.USB)
	if err != nil {
		slog.Error("NewPrinter: err: %v", err)
		return
	}

	if err := pr.Print(fmtStr, shared.PrinterArgs{}); err != nil {
		return
	}

}
