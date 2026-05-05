package repository

import (
	"context"
	"database/sql"
	"golangrest/internal/models"
	"log/slog"

	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(dbPath string) (*SQLiteProductRepository, error) {
	// Use WAL mode and busy timeout to handle concurrent writes from the worker pool
	dsn := dbPath + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	repo := &SQLiteProductRepository{db: db}
	repo.seedData() // Seed with initial data if empty

	return repo, nil
}

func (r *SQLiteProductRepository) seedData() {
	var count int
	if err := r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count); err != nil {
		slog.Warn("Failed to count products during seeding", "error", err)
		return
	}
	if count == 0 {
		if _, err := r.db.Exec("INSERT INTO products (name, price) VALUES ('Espresso', 3.50), ('Latte', 4.50), ('Cappuccino', 4.00)"); err != nil {
			slog.Warn("Failed to seed initial product data", "error", err)
		}
	}
}

func (r *SQLiteProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *SQLiteProductRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO products (name, price) VALUES (?, ?)", product.Name, product.Price)
	if err != nil {
		return models.Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Product{}, err
	}

	product.ID = int(id)
	return product, nil
}

func (r *SQLiteProductRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
