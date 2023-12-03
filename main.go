package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type Recipient struct {
	FirstName  string `json:"firstName`
	LastName   string `json:"lastName`
	Street     string `json:"street"`
	PostalCode string `json:"postalCode`
	City       string `json:"city"`
	Extra      string `json:"extra, omitempty"`
}

func print(recipient Recipient) ([]byte, error) {
	out, err := exec.Command("echo", "firstname:", recipient.FirstName, "\nlastname:", recipient.LastName, "\nadress: ", recipient.Street, "\ncomplement: ", recipient.Extra, "\npostal: ", recipient.PostalCode, recipient.City).Output()
	return out, err
}

func getRoot(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("got / request\n")
	var recipient Recipient
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Failed to read body: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Failed to read body: %v\n", err.Error()), http.StatusInternalServerError)
		return
	}
	fmt.Printf("Parsed body: '%s'\n", body)

	err = json.Unmarshal(body, &recipient)
	if err != nil {
		fmt.Printf("Failed to parse json body: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Failed to parse json body: %v\n", err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed input: '%#v'\n", recipient)

	out, err := print(recipient)
	if err != nil {
		fmt.Printf("Failed to print: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Failed to print: %v\n", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("Printing sent: \n%s\n", out))
}

func main() {
	http.HandleFunc("/", getRoot)
	err := http.ListenAndServe(":3333", nil)
	fmt.Printf("{%v}\n", err)
}
