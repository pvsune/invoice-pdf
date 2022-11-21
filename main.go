package main

import (
	"bytes"
	"flag"
	"fmt"
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
			TotalDue                                                            string
			ServiceDescription, InvoiceNumber, Now                              string
			NetDays                                                             int
		}{
			Name:                os.Getenv("NAME"),
			AddressLine1:        os.Getenv("ADDRESS_LINE_1"),
			AddressLine2:        os.Getenv("ADDRESS_LINE_2"),
			Email:               os.Getenv("EMAIL"),
			IBANNumber:          os.Getenv("ACCOUNT_NUMBER"),
			SwiftCode:           os.Getenv("SWIFT_CODE"),
			CompanyName:         os.Getenv("COMPANY_NAME"),
			CompanyAddressLine1: os.Getenv("COMPANY_ADDRESS_LINE_1"),
			CompanyAddressLine2: os.Getenv("COMPANY_ADDRESS_LINE_2"),
			CompanyEmail:        os.Getenv("COMPANY_EMAIL"),
			TotalDue:            os.Getenv("TOTAL_DUE"),
			ServiceDescription:  "Software Engineering Services for the period of " + invoicePeriod(),
			InvoiceNumber:       fmt.Sprintf("%03d", invoiceNum()),
			Now:                 time.Now().Format("01/02/06"),
			NetDays:             netDays(),
		}
	)
	if err := t.Execute(w, data); err != nil {
		return nil, err
	}
	return w, nil
}

func invoiceNum() int {
	d := time.Since(*startDate.Date)
	return int(d.Hours()/24/30) + 1
}

func invoicePeriod() string {
	now := time.Now()
	return fmt.Sprintf(
		"%s - %s",
		now.AddDate(0, 0, -now.Day()+1).Format("01/02/06"),
		now.AddDate(0, 1, -now.Day()).Format("01/02/06"),
	)
}

func netDays() int {
	var (
		now   = time.Now()
		start = now.AddDate(0, 0, -now.Day()+1)
		end   = now.AddDate(0, 1, -now.Day())

		date = start
		net  int
	)
	for !date.After(end) {
		if date.Weekday() != time.Saturday && date.Weekday() != time.Sunday {
			net += 1
		}
		date = date.AddDate(0, 0, 1)
	}
	return net
}
