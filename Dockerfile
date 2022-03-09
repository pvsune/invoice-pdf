FROM golang:1.14 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o invoice-pdf main.go


FROM icalialabs/wkhtmltopdf:alpine

WORKDIR /

COPY --from=builder /workspace/invoice-pdf .
COPY --from=icalialabs/wkhtmltopdf:alpine /bin/wkhtmltopdf /bin/wkhtmltopdf

ENTRYPOINT ["/invoice-pdf"]
