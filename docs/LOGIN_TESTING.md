# Login / Bruteforce Attack Test Scenario

This guide describes how to verify the `LOGIN` attack method using a local mock server. This ensures that the tool is correctly reading credentials and sending properly formatted requests without attacking real targets.

## Prerequisites

You need Go installed to run the mock server.

## 1. Create a Mock Server

Save the following code as `test_server.go`. This server listens on port 8090 and logs the body of every POST request it receives.

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			body, _ := io.ReadAll(r.Body)
			fmt.Printf("[RECEIVED] Method: %s | Body: %s\n", r.Method, string(body))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Mock Login Server started on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
```

## 2. Prepare Credentials

Create a file named `test_creds.txt` with username:password pairs (one per line):

```text
admin:123456
root:toor
user:password
guest:guest
test:test
```

## 3. Run the Test

1.  **Start the Server**:
    Open a terminal and run:

    ```bash
    go run test_server.go
    ```

2.  **Run the Attack**:
    Open a **new** terminal window and run the tool targeting the local server.
    _Note: We use `-proxy-file=` (empty value) to disable proxies and force a direct connection._

    ```bash
    go run main.go LOGIN http://localhost:8090 -data test_creds.txt -duration 10 -proxy-file=
    ```

## 4. Verify Results

Check the terminal window running `test_server.go`. You should see a flood of incoming requests with your credentials formatted in JSON:

```text
[RECEIVED] Method: POST | Body: {"username": "admin", "password": "123456"}
[RECEIVED] Method: POST | Body: {"username": "root", "password": "toor"}
[RECEIVED] Method: POST | Body: {"username": "user", "password": "password"}
...
```

If you see these logs, the `LOGIN` attack is working correctly!
