package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

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

	// Путь к статическим файлам
	staticDir := "../../static"

	// Определение маршрута для статических файлов
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Определение маршрута для страницы index.html
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "index_page.html"))
	})

	// Определение маршрута для страницы register.html
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "register.html"))
	})

	// Определение маршрута для страницы login.html
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "login.html"))
	})

	// Определение маршрута для страницы write.html
	router.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticDir, "write.html"))
	})

	// Начало прослушивания сервера
	fmt.Println("Server is listening on port 8081...")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
