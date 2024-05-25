package main

import (
	"fmt"
	"log"
	"net/http"

	"myapp/internal/handlers"
	"myapp/internal/logic"
	"myapp/internal/voice"
	"myapp/pkg/dbconnect"

	"github.com/gorilla/mux"
)

func main() {
	// Инициализация репозитория и логики
	mysqlDB, err := dbconnect.ConnectToMySQL()
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
		return
	}
	defer mysqlDB.Close()

	voiceRepo := voice.NewMySQLRepository(mysqlDB)
	voiceLogic := logic.NewVoiceLogic(voiceRepo)

	handlers := handlers.NewHandlers(voiceLogic)

	// Инициализация маршрутизатора
	router := mux.NewRouter()

	// Определение маршрутов
	router.HandleFunc("/api/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/voice/{user-id}", handlers.VoiceCreationHandler).Methods("POST")
	router.HandleFunc("/api/voice/{user-id}", handlers.VoiceUpdateHandler).Methods("PUT")

	// Начало прослушивания сервера
	fmt.Println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
