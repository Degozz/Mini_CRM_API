package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handler "mini-crm-back/internal/handlers"
	"mini-crm-back/internal/repository"
	"mini-crm-back/internal/service"
	"mini-crm-back/pkg/postgres"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	// 1. Загружаем переменные из файла .env
	_ = godotenv.Load()

	// 2. Подключение к Postgres с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		log.Fatal("Переменная окружения DB_CONN не задана.")
	}

	db, err := postgres.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к PostgreSQL", err)
	}
	defer db.Close()
	fmt.Println("Соединение с PostgreSQL успешно установлено!")

	// 3. Инициализация слоев
	clientRepo := repository.NewClientRepository(db)
	clientSvc := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientSvc)

	leadRepo := repository.NewLeadRepository(db)
	leadSvc := service.NewLeadService(leadRepo)
	leadHandler := handler.NewLeadHandler(leadSvc)

	// 4. Запуск Echo
	e := echo.New()

	// Эндпоинт корневого маршрута
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Привет!")
	})

	// Эндпоинт здоровья
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// 5. Маршруты для клиентов
	e.POST("/clients", clientHandler.Create)
	e.GET("/clients", clientHandler.GetAll)
	e.GET("/clients/:id", clientHandler.GetByID)
	e.PUT("/clients/:id", clientHandler.Update)
	e.DELETE("/clients/:id", clientHandler.Delete)

	// 6. Маршруты для заявок
	e.POST("/clients/:id/leads", leadHandler.Create)
	e.GET("/clients/:id/leads", leadHandler.GetByClientID)
	e.GET("/leads/:id", leadHandler.GetByID)
	e.PUT("/leads/:id", leadHandler.Update)
	e.DELETE("/leads/:id", leadHandler.Delete)

	// 7. Запуск сервера
	fmt.Println("Сервер Mini-CRM запущен на http://127.0.0.1:8080")
	e.Logger.Fatal(e.Start(":8080"))

}
