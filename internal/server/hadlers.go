package server

import (
	db "L0/internal/database"
	mod "L0/internal/models"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
)

var wg *sync.WaitGroup
var inMemory = db.NewInMemory()

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if s.Db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		log.Println("Error: Database connection is nil")
		return
	}

	jsonResp, err := json.Marshal(s.Db.Health())
	if err != nil {
		http.Error(w, "Error marshaling JSON response", http.StatusInternalServerError)
		log.Printf("Error marshaling JSON response: %v", err)
		return
	}

	_, err = w.Write(jsonResp)
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *Server) HandleGetOrderById(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		id := req.URL.Query().Get("order_uid")
		wg.Add(1)
		if order, ok := inMemory.Get(id, wg); ok {
			HandleNormalResponse(w, order)
			return
		}
		HandleGetError(w, sql.ErrNoRows)
		return
	}
}

func (s *Server) HandlePostOrderById(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var order mod.Order
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// десериализуем JSON в Order
		if err = json.Unmarshal(buf.Bytes(), &order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, present := inMemory.Get(order.OrderUID, wg); present {
			HandleError(w, http.StatusConflict, errors.New("заказ с таким uid есть в базе"))
			return
		}
		wg.Add(1)
		go inMemory.Save(order, wg)
		err = db.InsertOrder(s.Db.Db(), order)

		if err != nil {
			wg.Add(1)
			go inMemory.Remove(order.OrderUID, wg)
			log.Println(err.Error())
			HandleGetError(w, err)
			return
		}
		HandleNormalResponse(w, "Data saved")
	}
}

func (s *Server) HandleDeleteOrderById(writer http.ResponseWriter, request *http.Request) {

}
