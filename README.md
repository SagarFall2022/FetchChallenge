
## Prerequisites

1. Install Go (Golang): You can download and install Go from the official website: https://golang.org/dl/. Follow the installation instructions based on your operating system.

## Installing Gorilla Mux

Before running the API, you need to install the Gorilla Mux package. Open a terminal or command prompt and execute the following command:

```bash
go get -u github.com/gorilla/mux
```

## Building and Running the API

1. Copy the provided code into a file named `main.go`.

2. Open a terminal or command prompt and navigate to the directory containing the `main.go` file.

3. Build the API using the following command:

```bash
go build -o myapi
```

This will create an executable file named `myapi` (or `myapi.exe` on Windows) in the same directory.

4. Run the API server using the following command:

You can use `go run` to build and run the API in one step:

```bash
go run main.go
```

The server will start and listen on port 8000.

## Testing the API

### 1. Process Receipt

To process a receipt, you can use tools like `curl`, `Postman`, or any other HTTP client.

#### Using `curl`:

```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },
    {
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },
    {
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },
    {
      "shortDescription": "Klarbrunn 12-PK 12 FL OZ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}' http://localhost:8000/receipts/process
```

You will receive a response with the generated receipt ID.

### 2. Get Receipt Points

To get the points for a processed receipt, use the ID received from the previous step.

#### Using `curl`:

```bash
curl http://localhost:8000/receipts/{RECEIPT_ID}/points
```

Replace `{RECEIPT_ID}` with the actual ID received from the previous step.

You will receive a response with the calculated points for the receipt.

## Note

- This API is for demonstration purposes and does not use a persistent database. Receipt data is stored in memory and will be lost upon server restart.
