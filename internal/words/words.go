package words

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

// NewWordGenerator создает новый экземпляр WordGenerator
func NewWordGenerator(datasetPath string) ([]string, error) {
	// Читаем слова из датасета и сохраняем их в список
	wordList, err := readWordsFromFile(datasetPath)
	if err != nil {
		return nil, err
	}

	return wordList, nil
}

// Функция для чтения слов из файла и записи их в слайс
func readWordsFromFile(filename string) ([]string, error) {
	// Читаем содержимое файла
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Разбиваем содержимое файла на слова
	words := strings.Fields(string(content))

	return words, nil
}

// Функция для выбора трех случайных слов из слайса
func ChooseRandomWords(words []string, n int) []string {
	rand.NewSource(time.Now().UnixNano()) // Инициализируем генератор случайных чисел

	// Если в слайсе меньше слов, чем запрашивается, вернуть весь слайс
	if len(words) <= n {
		return words
	}

	// Создаем слайс для хранения выбранных слов
	randomWords := make([]string, n)

	// Выбираем случайные слова из исходного слайса
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(words))
		randomWords[i] = words[randomIndex]
	}

	return randomWords
}
