package main

var TemplateParseString1 = `
Terminal    : {{.Terminal}}
Cashier     : {{.LoginUser}}
Date        : {{FormatDate .Date "02/01/2006"}} {{.Time}}
Bill        : {{FormatDate .PaymentDate "02/01/2006"}} {{.PaymentTime}}
[Cover: {{.TagPax}}]

#Ticket Orders#
{{range .Orders}}Name: {{.Name}} {{FormatDecimal .Quantity 2}} {{FormatDecimal .Price 2}}
{{end}}

#Ticket.Discounts#
{{range .Discounts}}Discount: {{.Name}} | Amount: {{FormatDecimal .Amount 2}}
{{end}}

#Ticket.Services#
{{range .Services}}Service: {{.Name}} | Amount: {{FormatDecimal .Amount 2}}
{{end}}

#Ticket.Taxes#
{{range .Taxes}}Tax: {{.Name}} | Amount: {{FormatDecimal .Amount 2}}
{{end}}

#Ticket.Payments#
{{range .Payments}}Tendered: {{.Name}} | Amount: {{FormatDecimal .Tendered 2}}
Change: {{FormatDecimal .Tendered 2}}
RefNo: {{.PaymentInformation.RefNo}}
{{end}}

##Tickets.Orders##
{{range .Orders}}Name {{.Name}} {{FormatDecimal .Quantity 2}} {{FormatDecimal .Price 2}}
{{end}}
`

var TemplateParseString2 = `
<C10>CURRY VILLAGE
<C10>BANANA LEAF P/L
<C>8 LIM TECK KIM ROAD
<C>TEL : 6226 2562
<F>-
<C10>Receipt No: {{.TicketNo}}
<J00>Date: |{{.Date}} {{.Time}}
{{.Entities}}
<F>-
<J00> Qty Items|Price  Amount
{{.Orders}}
<F>=
<EB>
{{.Discounts}}
<J10>Total:|{{.TicketTotal}}
{{.PaymentDetails}}
<DB>
<F>=
<C10>THANK YOU
`
