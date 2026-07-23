package handler

import (
	"net/http"
	"strconv"

	"mini-crm-back/internal/models"
	"mini-crm-back/internal/service"

	"github.com/labstack/echo/v4"
)

// 1. Структура ClientHandler с изначально пустым svc-полем.
type ClientHandler struct {
	svc service.ClientService
}

// 1.1 Функция-конструктор NewClientHandler, которая возвращает экземпляр ClientHandler
func NewClientHandler(svc service.ClientService) *ClientHandler {
	return &ClientHandler{svc: svc}
}

// 1.2 Структура ErrorResponse для реализации вывода сообщения при ошибке в формате JSON
type ErrorResponse struct {
	Message string `json:"message"`
}

// 2. Метод Create для создания клиентов (POST /clients)
func (ch *ClientHandler) Create(c echo.Context) error {
	var clientData models.Client

	// 2.1. Читаем входящего JSON-данные из HTTP-запроса через функцию Bind()
	if err := c.Bind(&clientData); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json"})
	}

	ctx := c.Request().Context()

	// 2.2. Вызываем метод Create() в сервисном слое
	err := ch.svc.Create(ctx, &clientData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}

	return c.JSON(http.StatusCreated, clientData)
}

// 3. Метод GetAll для получения списка клиентов с возможностью фильтрации по их статусу и email (GET /clients)
func (ch *ClientHandler) GetAll(c echo.Context) error {
	status := c.QueryParam("status")
	email := c.QueryParam("email")

	ctx := c.Request().Context()

	clients, err := ch.svc.GetAll(ctx, status, email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}
	return c.JSON(http.StatusOK, clients)
}

// 4. Метод GetByID для получения клиента по его ID (GET /clients/:id)
func (ch *ClientHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "error: invalid id format"})
	}

	ctx := c.Request().Context()

	client, err := ch.svc.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "client not found"})
	}
	return c.JSON(http.StatusOK, client)
}

// 5. Метод Update для обновления списка клиентов (PUT /clients/:id)
func (ch *ClientHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid id format"})
	}

	var clients models.Client
	if err := c.Bind(&clients); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json"})
	}
	clients.ID = id

	ctx := c.Request().Context()

	err = ch.svc.Update(ctx, &clients)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}
	return c.JSON(http.StatusOK, clients)
}

// 6. Метод Delete для удаления клиента (DELETE /clients/:id)
func (ch *ClientHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "error: invalid id format"})
	}

	ctx := c.Request().Context()

	err = ch.svc.Delete(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}

	return c.JSON(http.StatusOK, ErrorResponse{Message: "client deleted"})
}
