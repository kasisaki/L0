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
)

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("error writing response. Err: %v", err)
	}
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

	if _, err = w.Write(jsonResp); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (s *Server) HandleGetOrderById(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		id := req.URL.Query().Get("order_uid")

		if order, ok := db.InMemory.Get(id); ok {
			HandleNormalResponse(w, order)
			return
		}
		HandleGetError(w, sql.ErrNoRows)
		return
	}
}

func (s *Server) HandlePostOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var order mod.Order
		var buf bytes.Buffer

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &order); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, present := db.InMemory.Get(order.OrderUID); present {
			HandleError(w, http.StatusConflict, errors.New("заказ с таким uid есть в базе"))
			return
		}

		if !db.InMemory.Save(order) {
			HandleError(w, http.StatusConflict, errors.New("не удалось сохранить заказ"))
			return
		}

		err = db.InsertOrder(db.Db(), order)
		if err != nil {
			log.Println("DB is unavailable. Data backed up...")
			db.InMemory.SaveToBackup(order)
			log.Println(err.Error())
			HandleNormalResponse(w, "DB is unavailable. Data backed up...")
			go db.InMemory.MigrateBackUpToDB()
			return
		}
		HandleNormalResponse(w, "Data saved")
	}
}

func (s *Server) HandleDeleteOrderById(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
