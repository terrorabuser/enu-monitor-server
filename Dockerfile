FROM golang:1.23.3-alpine

WORKDIR /app

# Устанавливаем необходимые зависимости
RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/

# Копируем файлы проекта
COPY go.mod ./  
COPY go.sum ./  
RUN go mod download  

COPY . .  

# Компилируем приложение
RUN go build -o main ./cmd/main.go  

# Открываем порт
EXPOSE 8080  

# Копируем файлы миграций
COPY migrations /migrations  

# Команда для запуска миграций и сервера
CMD /usr/local/bin/migrate -path /migrations -database "postgres://postgres:mypass@postgres_db:5432/postgres?sslmode=disable" up && ./main
