package service

import (
	"context"
	"mini-crm-back/internal/models"
	"mini-crm-back/internal/repository"
)

// 1. Интерфейс LeadService со всеми методами, которые наше приложение должно выполнять с заявками.
type LeadService interface {
	Create(ctx context.Context, lead *models.Lead) error
	GetByClientID(ctx context.Context, clientID int64, status string) ([]*models.Lead, error)
	GetByID(ctx context.Context, id int64) (*models.Lead, error)
	Update(ctx context.Context, lead *models.Lead) error
	Delete(ctx context.Context, id int64) error
}

// 2. Структура leadService с пустым repo-полем.
type leadService struct {
	repo repository.LeadRepository
}

// 2.1 Функция-конструктор NewLeadService, которая возвращает экземпляр leadService
func NewLeadService(repo repository.LeadRepository) LeadService {
	return &leadService{repo: repo}
}

func (ls *leadService) Create(ctx context.Context, lead *models.Lead) error {
	return ls.repo.Create(ctx, lead)
}

func (ls *leadService) GetByClientID(ctx context.Context, clientID int64, status string) ([]*models.Lead, error) {
	return ls.repo.GetByClientID(ctx, clientID, status)
}

func (ls *leadService) GetByID(ctx context.Context, id int64) (*models.Lead, error) {
	return ls.repo.GetByID(ctx, id)
}

func (ls *leadService) Update(ctx context.Context, lead *models.Lead) error {
	return ls.repo.Update(ctx, lead)
}

func (ls *leadService) Delete(ctx context.Context, id int64) error {
	return ls.repo.Delete(ctx, id)
}
