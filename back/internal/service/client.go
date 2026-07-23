package service

import (
	"context"
	"errors"
	"mini-crm-back/internal/models"
	"mini-crm-back/internal/repository"
	"net/mail"
)

// 1. Интерфейс ClientService со всеми методами, которые наше приложение должно выполнять с клиентами.
type ClientService interface {
	Create(ctx context.Context, client *models.Client) error
	GetAll(ctx context.Context, status string, email string) ([]*models.Client, error)
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
}

// 2. Структура clientService с изначально пустым repo-полем.
type clientService struct {
	repo repository.ClientRepository
}

// 2.1 Функция-конструктор NewClientService, которая возвращает экземпляр clientService
func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

// 3. Метод Create для валидации данных клиентов
func (cs *clientService) Create(ctx context.Context, client *models.Client) error {
	// 3.1. Проверка имени
	if client.Name == "" {
		return errors.New("Имя клиента не может быть пустым")
	}

	// 3.2. Проверка номера телефона
	if client.Phone == "" {
		return errors.New("Номер телефона не может быть пустым")
	}

	for _, r := range client.Phone {
		if r < '0' || r > '9' {
			return errors.New("номер телефона должен содержать только цифры")
		}
	}

	// 3.3. Проверка почты через функцию ParseAddress
	_, err := mail.ParseAddress(client.Email)
	if err != nil {
		return errors.New("Неправильный формат почты")
	}

	// 3.4. Если все проверки успешны, то передаем данные в слой репозиториев
	return cs.repo.Create(ctx, client)
}

func (cs *clientService) GetAll(ctx context.Context, status string, email string) ([]*models.Client, error) {
	return cs.repo.GetAll(ctx, status, email)
}

func (cs *clientService) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	return cs.repo.GetByID(ctx, id)
}

func (cs *clientService) Update(ctx context.Context, client *models.Client) error {
	return cs.repo.Update(ctx, client)
}

func (cs *clientService) Delete(ctx context.Context, id int64) error {
	return cs.repo.Delete(ctx, id)
}
