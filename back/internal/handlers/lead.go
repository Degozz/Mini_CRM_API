package handler

import (
	"net/http"
	"strconv"

	"mini-crm-back/internal/models"
	"mini-crm-back/internal/service"

	"github.com/labstack/echo/v4"
)

// 1. Структура LeadHandler с изначально пустым svc-полем.
type LeadHandler struct {
	svc service.LeadService
}

// 1.1 Функция-конструктор NewLeadHandler, которая возвращает экземпляр LeadHandler
func NewLeadHandler(svc service.LeadService) *LeadHandler {
	return &LeadHandler{svc: svc}
}

// 2. Метод Create для создания заявки (POST /clients/:id/leads)
func (lh *LeadHandler) Create(c echo.Context) error {
	idStr := c.Param("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid client id format"})
	}

	var lead models.Lead
	if err := c.Bind(&lead); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json"})
	}
	lead.ClientID = clientID

	ctx := c.Request().Context()

	err = lh.svc.Create(ctx, &lead)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}
	return c.JSON(http.StatusCreated, lead)
}

// 3. Метод GetByClientID для получения заявок по ID клиента с возможностью фильтрации по статусу заявок (GET /clients/:id/leads)
func (lh *LeadHandler) GetByClientID(c echo.Context) error {
	idStr := c.Param("id")
	clientID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid client id format"})
	}

	// Читаем статус из параметров запроса (?status=new)
	status := c.QueryParam("status")

	ctx := c.Request().Context()

	leads, err := lh.svc.GetByClientID(ctx, clientID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}
	return c.JSON(http.StatusOK, leads)
}

// 4. Метод GetByID для получения заявки по ID самой заявки (GET /leads/:id)
func (lh *LeadHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid id format"})
	}

	ctx := c.Request().Context()
	lead, err := lh.svc.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "lead not found"})
	}
	return c.JSON(http.StatusOK, lead)
}

// 5. Метод Update для обновления списка заявок (PUT /leads/:id)
func (lh *LeadHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid id format"})
	}

	var leadUpdate models.Lead
	if err := c.Bind(&leadUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid json"})
	}
	leadUpdate.ID = id

	ctx := c.Request().Context()

	err = lh.svc.Update(ctx, &leadUpdate)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}

	return c.JSON(http.StatusOK, leadUpdate)
}

// 6. Метод Delete для удаления заявки (DELETE /leads/:id)
func (lh *LeadHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "invalid id format"})
	}

	ctx := c.Request().Context()

	err = lh.svc.Delete(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "error"})
	}

	return c.JSON(http.StatusOK, ErrorResponse{Message: "lead deleted"})
}
