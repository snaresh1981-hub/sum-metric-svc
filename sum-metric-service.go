package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	keyPath = "key"
	contentTypeHeader = "content-type"
)

type SumMetricService struct {
}

type dataPayload struct {
	Value int `json:"value"`
}

var payload dataPayload

type responsePayload struct {
	Value int `json:"value"`
}

func (s *SumMetricService) getInfo(w http.ResponseWriter, r *http.Request) {
	message := "SumMetricService is alive!"
	log.Println(message)
	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprint(w, message); err != nil {
		log.Printf("getInfo error: %v", err)
	}
}

func (s *SumMetricService) getMetricsSum(w http.ResponseWriter, r *http.Request){
	response := responsePayload{}
    var sum = 0
	key, err := s.fetchIDs(r)
	if err != nil {
		s.badRequest(w, err)
		return
	}
	data, found := cache.get(key,1*time.Hour)
	if(!found){
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _,value := range data{
		sum = sum + value
	}
	response.Value = sum
	responseMessage, err := json.Marshal(response)
	if err != nil{
		log.Printf("json marshalling error: %s", err)
		s.badRequest(w, err)
		return
	}
	log.Printf("getMetrics response: %s", string(responseMessage))
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, string(responseMessage)); err != nil {
		log.Printf("getInfo error: %v", err)
	}
}

func (s *SumMetricService) postMetricData(w http.ResponseWriter, r *http.Request) {
	log.Println("POST data invoked")
	key, err := s.fetchIDs(r)
	if err != nil {
		s.badRequest(w, err)
		return
	}
	if err := s.decodeBody(r, &payload); err != nil {
		s.badRequest(w, err)
		return
	}
	log.Printf("Received data for key = '%s' ", key)
	log.Printf(" - Payload: %d ", payload.Value)
	cache.Add(key,payload.Value)
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "%s", string("{}")); err != nil {
		log.Printf("postData error: %v", err)
	}
}

func (s *SumMetricService) fetchIDs(r *http.Request) (string,error) {
	vars := mux.Vars(r)

	if _, found := vars[keyPath]; !found {
		return "", fmt.Errorf("%s not  found", keyPath)
	}

	return vars[keyPath], nil
}

func (s *SumMetricService) badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	log.Println(err.Error())

	if _, err := fmt.Fprint(w, err.Error()); err != nil {
		log.Printf("postData error: %v", err)
	}

}

func (s *SumMetricService) decodeBody(r *http.Request, payload *dataPayload) error {
	bodyData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err
	} else if len(bodyData) == 0 {
		return errors.New("body is empty")
	}

	fmt.Println("Data: ", string(bodyData))

	if err := json.Unmarshal(bodyData, payload); err != nil {
		return fmt.Errorf("could not decode JSON body: %v", err)
	}

	return nil
}

func (s *SumMetricService) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	message := fmt.Sprintf("Path '%s' not found\n", r.RequestURI)
	log.Printf(message)
	if _, err := fmt.Fprint(w, message); err != nil {
		log.Printf("postData error: %v", err)
	}
}
