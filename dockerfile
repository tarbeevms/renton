# Используем официальный образ Golang как базовый образ
FROM golang:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /../voice

# Устанавливаем Python и pip для Ubuntu, а также необходимые библиотеки
RUN apt-get update && apt-get install -y python3 python3-pip python3-venv libsndfile1 ffmpeg

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости перед копированием остальных файлов
RUN go mod download

# Копируем остальные файлы вашего приложения внутрь контейнера
COPY . .

# Изменяем права доступа к файлу конфигурации
RUN chmod 777 /voice/config.yml

# Переходим в директорию cmd/redditclone
WORKDIR /voice/cmd/voice

# Создаем и активируем виртуальное окружение Python
RUN python3 -m venv venv
RUN /bin/bash -c "source venv/bin/activate && pip install --no-cache-dir -r req.txt"

# Собираем ваше приложение
RUN go build -o myapp .

# Указываем порт, который будет использоваться вашим приложением
EXPOSE 8081

# Команда для запуска вашего приложения при старте контейнера
CMD ["./myapp"]
