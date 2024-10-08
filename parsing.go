/*
Provides utilities for parsing and formatting from structured data

Key Features:
  - **Parser**: A struct for converting data structures into formatted strings using Go's templating system
  - **Template Parsing**:
    The `ParseTemplateTemplateV1`: function allows flexible parsing of any data structure with optional custom
    functions for advanced formatting.
  - **Custom Functions**:
  - `FormatDecimal`: Formats float to a specified precision.
  - `FormatDate`: Formats `time.Time` objects.
  - **Use Cases**:
  - `ParseMessageForPrint`: Generates print-ready strings using a `RootObject`.
  - `ParseMessageOnly`: Formats messages using a `Ticket` structure.
*/
package main

import (
	"fmt"
	"strings"
	"text/template"
	"time"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseMessageForPrint(templateString string, ro RootObject) (string, error) {
	formattedStr, err := ParseTemplateTemplateV1[RootObject](templateString, ro, nil)
	if err != nil {
		return "", err
	}
	return formattedStr, nil
}

func (p *Parser) ParseMessageOnly(templateString string, ticket Ticket) (string, error) {
	fnMap := FuncMapTicket(ticket)
	formattedStr, err := ParseTemplateTemplateV1[Ticket](templateString, ticket, &fnMap)
	if err != nil {
		return "", err
	}
	return formattedStr, nil
}

func FormatDecimal(value float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, value)
}

func FormatDate(date time.Time, format string) string {
	return date.Format(format)
}

func FuncMapTicket(ticket Ticket) template.FuncMap {
	return template.FuncMap{
		"FormatDecimal": FormatDecimal,
		"FormatDate":    FormatDate,
	}
}

func ParseTemplateTemplateV1[T any](tmplStr string, obj T, funcMap *template.FuncMap) (string, error) {
	var tmpl *template.Template
	var err error

	switch {
	case funcMap != nil:
		tmpl, err = template.New("tmplStr").Funcs(*funcMap).Parse(tmplStr)
		if err != nil {
			return "", err
		}
	default:
		tmpl, err = template.New("tmplStr").Parse(tmplStr)
		if err != nil {
			return "", err
		}
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, obj); err != nil {
		return "", err
	}

	return result.String(), nil
}
