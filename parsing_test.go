package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFormatDecimal(t *testing.T) {
	type args struct {
		value     float64
		precision int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "precise 2",
			args: args{
				value:     100,
				precision: 2,
			},
			want: "100.00",
		},
		{
			name: "precise 0",
			args: args{
				value:     1,
				precision: 0,
			},
			want: "1",
		},
		{
			name: "precise 1",
			args: args{
				value:     50,
				precision: 1,
			},
			want: "50.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDecimal(tt.args.value, tt.args.precision); got != tt.want {
				t.Errorf("FormatDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	type args struct {
		date   time.Time
		format string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "date format 02/01/2006",
			args: args{
				date:   time.Now(),
				format: "02/01/2006",
			},
			want: time.Now().Format("02/01/2006"),
		},
		{
			name: "date format 2006-01-02",
			args: args{
				date:   time.Now(),
				format: time.DateOnly,
			},
			want: time.Now().Format(time.DateOnly),
		},
		{
			name: "date format 2006-01-02T15:04:05Z07:00",
			args: args{
				date:   time.Now(),
				format: time.RFC3339,
			},
			want: time.Now().Format(time.RFC3339),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDate(tt.args.date, tt.args.format); got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

var wantTemplate1 = `
Terminal    : Terminal 1
Cashier     : login user test
Date        : %s 14:30
Bill        : %s 14:31
[Cover: 4]

#Ticket Orders#
Name: Burger 2.00 5.00
Name: Fries 1.00 2.50


#Ticket.Discounts#
Discount: Promo 10%% | Amount: -1.00


#Ticket.Services#
Service: Service Charge | Amount: 1.50


#Ticket.Taxes#
Tax: GST | Amount: 0.35


#Ticket.Payments#
Tendered: Cash | Amount: 9.00
Change: 9.00
RefNo: 123456


##Tickets.Orders##
Name Burger 2.00 5.00
Name Fries 1.00 2.50

`

func getDateToday() string {
	now := time.Now()
	return now.Format("02/01/2006")
}

func TestParseTemplateTicket(t *testing.T) {
	type args struct {
		tmplStr string
		ticket  Ticket
	}

	tdyDt := getDateToday()

	nParser := NewParser()
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty template",
			args: args{
				tmplStr: "",
				ticket:  Ticket{},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "parse 1",
			args: args{
				tmplStr: TemplateParseString1,
				ticket: Ticket{
					Terminal:    "Terminal 1",
					LoginUser:   "login user test",
					Date:        time.Now(),
					Time:        "14:30",
					PaymentDate: time.Now(),
					PaymentTime: "14:31",
					TagPax:      "4",
					Orders: []Order{
						{Name: "Burger", Quantity: 2, Price: 5.0},
						{Name: "Fries", Quantity: 1, Price: 2.5},
					},
					Discounts: []Discount{
						{Name: "Promo 10%", Amount: -1.0},
					},
					Services: []Service{
						{Name: "Service Charge", Amount: 1.5},
					},
					Taxes: []Tax{
						{Name: "GST", Amount: 0.35},
					},
					Payments: []Payment{
						{
							Name:               "Cash",
							Tendered:           9.0,
							PaymentInformation: PaymentInfo{RefNo: "123456"}},
					},
				},
			},
			want:    fmt.Sprintf(wantTemplate1, tdyDt, tdyDt),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nParser.ParseMessageOnly(tt.args.tmplStr, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

var wantTemplate2 = `
<C10>CURRY VILLAGE
<C10>BANANA LEAF P/L
<C>8 LIM TECK KIM ROAD
<C>TEL : 6226 2562
<F>-
<C10>Receipt No: 1000
<J00>Date: |12/08/2024 10:15
test entities
<F>-
<J00> Qty Items|Price  Amount
test orders
<F>=
<EB>
test discounts
<J10>Total:|test total ticket
test payment details
<DB>
<F>=
<C10>THANK YOU
`

func TestParseTemplatePrinter(t *testing.T) {
	type args struct {
		tmplStr    string
		rootObject RootObject
	}

	nParser := NewParser()
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty template",
			args: args{
				tmplStr:    "",
				rootObject: RootObject{},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "parse 1",
			args: args{
				tmplStr: TemplateParseString2,
				rootObject: RootObject{
					TicketNo:       "1000",
					Date:           "12/08/2024",
					Time:           "10:15",
					Entities:       "test entities",
					Orders:         "test orders",
					Discounts:      "test discounts",
					TicketTotal:    "test total ticket",
					PaymentDetails: "test payment details",
				},
			},
			want:    wantTemplate2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := nParser.ParseMessageForPrint(tt.args.tmplStr, tt.args.rootObject)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
