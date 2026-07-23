package repository

import (
	"context"
	"fmt"
	"mini-crm-back/internal/models"
	"mini-crm-back/pkg/postgres"
)

// 1. Интерфейс ClientRepository со всеми методами для работы с таблицей клиентов.
type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetAll(ctx context.Context, status string, email string) ([]*models.Client, error)
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
}

// 2. Структура clientRepo с изначально пустым db-полем.
type clientRepo struct {
	db *postgres.DB
}

// 2.1 Функция-конструктор NewClientRepository, которая возвращает экземпляр clientRepo с указателем на экземпляр адаптера.
func NewClientRepository(db *postgres.DB) ClientRepository {
	return &clientRepo{db: db}
}

// 3. Метод Create для создания клиентов (POST /clients)
func (r *clientRepo) Create(ctx context.Context, client *models.Client) error {
	// 3.1. Пишем SQL-запрос для последующей отправки.
	query := `
		INSERT INTO clients (name, email, phone, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	// 3.2. Отправляем SQL-запрос через QueryRow-метод.
	row := r.db.QueryRow(ctx, query, client.Name, client.Email, client.Phone, client.Status)

	// 3.3. Записываем с помощью Scan-метода в структуру Client возвращенные базой данных ID, CreatedAt и UpdatedAt
	err := row.Scan(&client.ID, &client.CreatedAt, &client.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// 4. Метод GetAll для получения списка клиентов с учетом фильтров и сортировки (GET /clients)
func (r *clientRepo) GetAll(ctx context.Context, status string, email string) ([]*models.Client, error) {
	// 4.1. Пишем SQL-запрос для последующей отправки.
	query := `
		SELECT id, name, email, phone, status, created_at, updated_at 
		FROM clients 
		WHERE 1=1
	`

	var args []any
	argID := 1

	// 4.2. Фильтрация по статусу клиентов.
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argID)
		args = append(args, status)
		argID++
	}

	// 4.3. Фильтрация по email клиентов.
	if email != "" {
		query += fmt.Sprintf(" AND email = $%d", argID)
		args = append(args, email)
		argID++
	}

	// 4.4. Добавляем условие сортировки в конец готового SQL-запроса.
	query += " ORDER BY created_at DESC"

	// 4.5. Отправляем SQL-запрос в Postgres с помощью Query-метода.
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 4.6. Создаём clients-слайс и через цикл добавляем туда строки из БД
	var clients []*models.Client
	for rows.Next() {
		var c models.Client
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, &c)
	}

	// 4.7. Проверяем, не прервался ли цикл из-за ошибки в процессе чтения
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

// 5. Метод GetByID для получения клиента по ID (GET /clients/:id)
func (r *clientRepo) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	var client models.Client

	// 5.1. Пишем SQL-запрос для последующей отправки.
	query := `
		SELECT id, name, email, phone, status, created_at, updated_at
		FROM clients
		WHERE id = $1
	`

	// 5.2. Отправляем SQL-запрос в Postgres с помощью QueryRow-метода.
	row := r.db.QueryRow(ctx, query, id)

	// 5.3. Записываем с помощью Scan-метода в структуру Client возвращенные базой данных ID, Name, Email..., UpdatedAt
	err := row.Scan(&client.ID, &client.Name, &client.Email, &client.Phone, &client.Status, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// 5.4. Возвращаем указатель на заполненную структуру models.Client и nil для ошибки
	return &client, nil
}

// 6. Метод Update для изменения данных клиента (PUT /clients/:id)
func (r *clientRepo) Update(ctx context.Context, client *models.Client) error {
	// 6.1. Пишем SQL-запрос для последующей отправки.
	query := `
		UPDATE clients
		SET name = $1, email = $2, phone = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
		RETURNING updated_at
	`

	// 6.2. Отправляем SQL-запрос в Postgres с помощью QueryRow-метода.
	row := r.db.QueryRow(ctx, query, client.Name, client.Email, client.Phone, client.Status, client.ID)

	// 6.3. Записываем с помощью Scan-метода в структуру Client возвращенные базой данных UpdatedAt
	err := row.Scan(&client.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// 7. Метод Delete для удаления клиента по его id (DELETE /clients/:id)
func (r *clientRepo) Delete(ctx context.Context, id int64) error {
	// 7.1. Пишем SQL-запрос для последующей отправки.
	query := `DELETE FROM clients WHERE id = $1`

	// 7.2. Отправляем SQL-запрос в Postgres с помощью Exec-метода.
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
