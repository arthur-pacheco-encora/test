package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type Asset struct {
	Asset      string   `json:"asset"`
	Validator  string   `json:"validator"`
	Operations []string `json:"operations"`
	OnChain    bool     `json:"on_chain"`
	Claimable  bool     `json:"claimable"`
	Active     bool     `json:"active"`
}

type UpdateRequestBody = []Asset

var client *storage.Client

func initClient() {
	ctx := context.Background()
	var err error
	client, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
}

func GetAssets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	initClient()

	bucketName := "billingcalc-data"
	queryValues := r.URL.Query()
	filePath := queryValues.Get("filePath")

	if filePath == "" {
		http.Error(w, "filePath parameter is required", http.StatusBadRequest)
		return
	}

	data, err := getDataFromBucket(bucketName, filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func UpdateAssets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	initClient()

	bucketName := "billingcalc-data"
	queryValues := r.URL.Query()
	filePath := queryValues.Get("filePath")

	if filePath == "" {
		http.Error(w, "filePath parameter is required", http.StatusBadRequest)
		return
	}

	modifiedData, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("Error closing request body: %v", closeErr)
		}
	}()

	if !validBody(modifiedData) {
		http.Error(w, "Invalid body request", http.StatusInternalServerError)
		return
	}

	if err := sendDataToBucket(bucketName, filePath, modifiedData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("File updated successfully"))
	if err != nil {
		log.Printf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}
}

func getDataFromBucket(bucketName, filePath string) ([]byte, error) {
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	obj := bucket.Object(filePath)
	rc, err := obj.NewReader(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new reader for object")
	}
	defer func() {
		if closeErr := rc.Close(); closeErr != nil {
			log.Printf("Error closing reader: %v", closeErr)
		}
	}()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read data from object")
	}

	return data, nil
}

func sendDataToBucket(bucketName, filePath string, data []byte) error {
	ctx := context.Background()
	bucket := client.Bucket(bucketName)

	obj := bucket.Object(filePath)
	wc := obj.NewWriter(ctx)
	if _, err := wc.Write(data); err != nil {
		return errors.New("error writing data to bucket: " + err.Error())
	}
	if err := wc.Close(); err != nil {
		return errors.New("error closing writer: " + err.Error())
	}

	return nil
}

func validBody(body []byte) bool {
	err := json.Unmarshal(body, &UpdateRequestBody{})
	if err != nil {
		return false
	}
	return true
}
