package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang_gpt/internal/app"
	
)

func main() {
	// Инициализируем приложение
	application, err := app.NewApp()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения:", err)
	}

	// Запускаем HTTP сервер
	go func() {
		if err := application.Router.Run(":8080"); err != nil {
			log.Fatal("Ошибка запуска HTTP сервера:", err)
		}
	}()

	// Ожидаем сигнал завершения (Ctrl+C или завершение процесса)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Блокируем выполнение, пока не придёт сигнал

	log.Println("Выключение сервера...")
}
