package repository

import (
	"context"
	"fmt"
	"mini-crm-back/internal/models"
	"mini-crm-back/pkg/postgres"
)

// 1. Интерфейс LeadRepository со всеми методами для работы с таблицей заявок.
type LeadRepository interface {
	Create(ctx context.Context, lead *models.Lead) error
	GetByClientID(ctx context.Context, clientID int64, status string) ([]*models.Lead, error)
	GetByID(ctx context.Context, id int64) (*models.Lead, error)
	Update(ctx context.Context, lead *models.Lead) error
	Delete(ctx context.Context, id int64) error
}

// 2. Структура leadRepo с пустым db-полем.
type leadRepo struct {
	db *postgres.DB
}

// 2.1 Функция-конструктор NewLeadRepository, которая возвращает экземпляр leadRepo с указателем на экземпляр адаптера.
func NewLeadRepository(db *postgres.DB) LeadRepository {
	return &leadRepo{
		db: db,
	}
}

// 3. Метод Create для создания заявки конкретному клиенту (POST /clients/:id/leads)
func (r *leadRepo) Create(ctx context.Context, lead *models.Lead) error {
	query := `
		INSERT INTO leads (client_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	row := r.db.QueryRow(ctx, query, lead.ClientID, lead.Title, lead.Description, lead.Status)

	err := row.Scan(&lead.ID, &lead.CreatedAt, &lead.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// 4. Метод GetByClientID для получения списка заявок конкретного клиента (GET /clients/:id/leads)
func (r *leadRepo) GetByClientID(ctx context.Context, clientID int64, status string) ([]*models.Lead, error) {
	query := `
		SELECT id, client_id, title, description, status, created_at, updated_at 
		FROM leads 
		WHERE client_id = $1
	`

	var args []any
	args = append(args, clientID)
	argID := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, status)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leads []*models.Lead
	for rows.Next() {
		var l models.Lead
		err := rows.Scan(&l.ID, &l.ClientID, &l.Title, &l.Description, &l.Status, &l.CreatedAt, &l.UpdatedAt)
		if err != nil {
			return nil, err
		}
		leads = append(leads, &l)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return leads, nil

}

// 5. Метод GetByID для получения конкретной заявки клиента по ID заявки (GET /leads/:id)
func (r *leadRepo) GetByID(ctx context.Context, id int64) (*models.Lead, error) {
	var lead models.Lead
	query := `
		SELECT id, client_id, title, description, status, created_at, updated_at
		FROM leads 
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&lead.ID, &lead.ClientID, &lead.Title, &lead.Description, &lead.Status, &lead.CreatedAt, &lead.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &lead, nil
}

// 6. Метод Update для изменения статуса или данных заявки (PUT /leads/:id)
func (r *leadRepo) Update(ctx context.Context, lead *models.Lead) error {
	query := `
		UPDATE leads 
		SET title = $1, description = $2, status = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	row := r.db.QueryRow(ctx, query, lead.Title, lead.Description, lead.Status, lead.ID)

	err := row.Scan(&lead.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// 7. Метод Delete для удаления заявки по его id (DELETE /leads/:id)
func (r *leadRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM leads WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
