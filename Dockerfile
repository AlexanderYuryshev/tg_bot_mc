# Образ golang-alpine почти в 3 раза меньше обычного, есть нюансы,
# но для самых простых приложений пойдет
FROM golang:1.24-alpine

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы
COPY . .

# Компилируем приложение
RUN go build -o tg-bot .

# Запускаем приложение
CMD ["./tg-bot"]
