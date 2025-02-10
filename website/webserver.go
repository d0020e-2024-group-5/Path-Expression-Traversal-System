// för att köra programmet 'go run webserver.go'
// om det inte funkar 'go mod init webserver.go' sen det tidigare commandot

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type RequestData struct {
	Data string `json:"data"`
}

type ResponseData struct {
	Message string `json:"message"`
}

func main() {
	// route handlers
	http.HandleFunc("/", handler) // servers the main HTML file

	http.HandleFunc("/api/data", handleAPI)      // API endpoint to return static message
	http.HandleFunc("/api/submit", handleSubmit) // API endpoint to handle form submission

	fmt.Printf("Server is running on http://localhost:8080")

	// start the TTP server on port 80080
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// return 404 for other paths
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return
	//}

	// reads the html file
	html, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	// set the response content type to HTML and write file content
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", html)
}

// for the Get Server Data button, returns static JSON message
func handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Hello from server"}`)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	log.Print("api subbmit was requested")

	// checking request method is POST
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// parse the request data into requestData struct
	var requestData RequestData
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	fmt.Println("Received from client:", requestData.Data)

	// create a response containging the recieived data
	response := ResponseData{Message: "You entered: " + requestData.Data}

	// set response content type to JSON and send it back
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// using gin framework
//func main() {
//	r := gin.Default()
//	r.GET("/", func(c *gin.Context) {
//		c.String(http.StatusOK, "Hello World")
//	})
//	r.Run(":8080")
//}

//func handler(c *gin.Context) {
//	c.String(http.StatusOK, "Hello World!")
//}

// using fiber framework
//func main() {
//	app := fiber.New()
//
//	app.Get("/", func(c *fiber.Ctx) error {
//		return c.SendString("Hello World!")
//	})
//
//	app.Listen(":8080")
//}
