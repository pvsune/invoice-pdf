# invoice-pdf
For software contractors, automate the creation of monthly invoices to client.

![sample invoice](https://github.com/pvsune/invoice-pdf/blob/master/sample.png?raw=true)

## Quickstart
1. Write the `.env` file used to extract values from:
    ```
    NAME=Your name
    ADDRESS_LINE_1=Your address 1
    ADDRESS_LINE_2=Your address 2
    EMAIL=your@email.com
    COMPANY_NAME=Client name
    COMPANY_ADDRESS_LINE_1=Client address 1
    COMPANY_ADDRESS_LINE_2=Client address 2
    COMPANY_EMAIL=client@email.com
    ```
1. Build the docker image:
    ```
    $ docker build -t invoice-pdf:latest .
    ```
1. Create the PDF:
    ``` 
    $ docker run --rm -v $(pwd):/root --env-file .env invoice-pdf:latest
    ```
