package handlers

import (
	"myapp/internal/words"
	"myapp/pkg/io"
	"net/http"
)

// WordsHandler обрабатывает запросы на получение случайных слов
func WordsHandler(w http.ResponseWriter, r *http.Request) {
	// Создаем новый генератор слов
	wordGen, err := words.NewWordGenerator("../../russianwords.txt")
	if err != nil {
		io.SendError(w, "Failed to initialize word generator", http.StatusInternalServerError)
		return
	}

	// Генерируем 3 случайных слова
	randomWords := words.ChooseRandomWords(wordGen, 3)

	// Отправляем случайные слова в ответ
	io.WriteJSON(w, http.StatusOK, map[string][]string{"words": randomWords})
}
