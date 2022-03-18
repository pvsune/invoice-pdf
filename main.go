package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	pdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

var (
	inputTemplate = "/root/invoice.html"
	outputPDF     = "/root/invoice.pdf"
)
var (
	startDate = &DateValue{}
)

func init() {
	flag.Var(startDate, "start", "Your start date. Used to compute inovice number (e.g. Sep 1991)")
	flag.Parse()
	if startDate.Date.After(time.Now()) {
		log.Fatalf("invalid start date")
	}
}

func main() {
	w, err := parseTemplate(inputTemplate)
	if err != nil {
		log.Fatalf("cannot parse template: %s", err)
	}

	if err := writeToPDF(w, outputPDF); err != nil {
		log.Fatalf("cannot write to PDF: %s", err)
	}
	log.Println("Successfully created PDF: " + outputPDF)
}

func writeToPDF(w io.Reader, name string) error {
	pdfg, err := pdf.NewPDFGenerator()
	if err != nil {
		return err
	}
	pdfg.AddPage(pdf.NewPageReader(w))
	if err := pdfg.Create(); err != nil {
		return err
	}
	if err := pdfg.WriteFile("./" + name); err != nil {
		log.Fatalf("cannot write PDF file: %s", err)
	}
	return nil
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

func invoiceNum() int {
	d := time.Since(*startDate.Date)
	if num := int(d.Hours() / 24 / 30); num > 0 {
		return num
	}
	return 1
}
