package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

var (
	templateFile = "invoice.html"
)

func main() {
	w, err := parseTemplate(templateFile)
	if err != nil {
		log.Fatalf("cannot parse template: %s", err)
	}

	log.Println(w.String())
}

func parseTemplate(name string) (*bytes.Buffer, error) {
	fileBytes, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	t, err := template.New("invoice").Parse(string(fileBytes))
	if err != nil {
		return nil, err
	}

	var (
		w    = &bytes.Buffer{}
		data = struct {
			Name, AddressLine1, AddressLine2, Email                             string
			IBANNumber, SwiftCode                                               string
			CompanyName, CompanyAddressLine1, CompanyAddressLine2, CompanyEmail string
			ServiceDescription                                                  string
			NetDays                                                             int
			TotalDue                                                            int
		}{
			Name:                os.Getenv("NAME"),
			AddressLine1:        os.Getenv("ADDRESS_LINE_1"),
			AddressLine2:        os.Getenv("ADDRESS_LINE_2"),
			Email:               os.Getenv("EMAIL"),
			IBANNumber:          "xxx",
			SwiftCode:           "xxx",
			CompanyName:         os.Getenv("COMPANY_NAME"),
			CompanyAddressLine1: os.Getenv("COMPANY_ADDRESS_LINE_1"),
			CompanyAddressLine2: os.Getenv("COMPANY_ADDRESS_LINE_2"),
			CompanyEmail:        os.Getenv("COMPANY_EMAIL"),
			ServiceDescription:  "Software Engineering Services for the period of xx/xx/xx - xx/xx/xx",
			NetDays:             30,
			TotalDue:            30,
		}
	)
	if err := t.Execute(w, data); err != nil {
		return nil, err
	}
	return w, nil
}
