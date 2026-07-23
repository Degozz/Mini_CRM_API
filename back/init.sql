DROP TABLE IF EXISTS clients CASCADE;
CREATE TABLE clients
(
    id      BIGSERIAL PRIMARY KEY,
    name    VARCHAR(100) NOT NULL,
    email   VARCHAR(255) UNIQUE NOT NULL,
    phone   VARCHAR(50),
    status  VARCHAR(50) DEFAULT 'new',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


DROP TABLE IF EXISTS leads CASCADE;
CREATE TABLE leads
(
    id      BIGSERIAL PRIMARY KEY,
    client_id  BIGINT NOT NULL,
    title    VARCHAR(100) NOT NULL,
    description   VARCHAR(255) NOT NULL,
    status  VARCHAR(50) DEFAULT 'new',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (client_id) REFERENCES clients(id) ON DELETE CASCADE
);


-- Ускоряем поиск лидов конкретного клиента, заменяя перебор всей таблицы на мгновенную выборку по индексу.
CREATE INDEX IF NOT EXISTS idx_leads_client_id ON leads(client_id);
