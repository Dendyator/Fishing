package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Dendyator/Fishing/internal/database"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func SetupRoutes(r *mux.Router, db *database.DB) {
	r.HandleFunc("/api/fishing-data", CreateFishingData(db)).Methods("POST")
	r.HandleFunc("/api/fishing-data/{id}", GetFishingDataByID(db)).Methods("GET")
	r.HandleFunc("/api/fishing-data/{id}", UpdateFishingData(db)).Methods("PUT")
	r.HandleFunc("/api/fishing-data/all", GetAllFishingData(db)).Methods("GET")

	r.HandleFunc("/api/fishing-methods", GetAllFishingMethods(db)).Methods("GET")
	r.HandleFunc("/api/fish-species", GetAllFishSpecies(db)).Methods("GET")
	r.HandleFunc("/api/baits", GetAllBaits(db)).Methods("GET")
}

func CreateFishingData(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data database.FishingData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logrus.Errorf("Ошибка декодирования тела запроса: %v", err)
			http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
			return
		}

		data.Date = time.Now()
		if err := db.CreateFishingData(data); err != nil {
			logrus.Errorf("Ошибка создания данных о рыбалке: %v", err)
			http.Error(w, "Не удалось создать данные о рыбалке", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	}
}

func GetFishingDataByID(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logrus.Errorf("Ошибка парсинга ID: %v", err)
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		data, err := db.GetFishingDataByID(id)
		if err != nil {
			logrus.Errorf("Ошибка получения данных о рыбалке: %v", err)
			http.Error(w, "Нет данных о рыбалке для указанного ID", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}

func UpdateFishingData(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			logrus.Errorf("Ошибка парсинга ID: %v", err)
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var data database.FishingData
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			logrus.Errorf("Ошибка декодирования тела запроса: %v", err)
			http.Error(w, "Неверное тело запроса", http.StatusBadRequest)
			return
		}

		if err := db.UpdateFishingData(id, data); err != nil {
			logrus.Errorf("Ошибка обновления данных о рыбалке: %v", err)
			http.Error(w, "Не удалось обновить данные о рыбалке", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}

func GetAllFishingData(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := db.GetAllFishingData()
		if err != nil {
			logrus.Errorf("Ошибка получения данных о рыбалке: %v", err)
			http.Error(w, "Не удалось получить данные о рыбалке", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(data)
	}
}

func GetAllFishingMethods(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		methods, err := db.GetAllFishingMethods()
		if err != nil {
			logrus.Errorf("Ошибка получения способов ловли: %v", err)
			http.Error(w, "Не удалось получить способы ловли", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(methods)
	}
}

func GetAllFishSpecies(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		species, err := db.GetAllFishSpecies()
		if err != nil {
			logrus.Errorf("Ошибка получения видов рыб: %v", err)
			http.Error(w, "Не удалось получить виды рыб", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(species)
	}
}

func GetAllBaits(db *database.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		baits, err := db.GetAllBaits()
		if err != nil {
			log.Printf("Ошибка получения приманок: %v", err)
			http.Error(w, "Не удалось получить приманки", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(baits)
	}
}
