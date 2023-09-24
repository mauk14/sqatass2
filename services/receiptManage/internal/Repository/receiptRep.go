package Repository

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"messanger/services/receiptManage/internal/Domain"
)

type ReceiptRepository interface {
	Create(context.Context, *Domain.Receipt) error
	Delete(context.Context, int64) error
	Update(context.Context, *Domain.Receipt) error
	Get(context.Context, int64) (*Domain.Receipt, error)
	GetAll(context.Context) ([]*Domain.Receipt, error)
}

type receiptRepository struct {
	db *pgxpool.Pool
}

func (r *receiptRepository) Create(ctx context.Context, receipt *Domain.Receipt) error {
	query :=
		`INSERT INTO receipts (title, author, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`
	args := []interface{}{receipt.Title, receipt.Author, receipt.Description}
	return r.db.QueryRow(ctx, query, args...).Scan(&receipt.Id, &receipt.Created_at)
}

func (r *receiptRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM receipts
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("Record not found")
	}

	return nil

}

func (r *receiptRepository) Get(ctx context.Context, id int64) (*Domain.Receipt, error) {
	query :=
		`Select id, title, author, description, created_at, upload_at from receipts
		WHERE id = $1`

	var receipt Domain.Receipt

	err := r.db.QueryRow(ctx, query, id).Scan(
		&receipt.Id,
		&receipt.Title,
		&receipt.Author,
		&receipt.Description,
		&receipt.Created_at,
		&receipt.Updated_at,
	)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *receiptRepository) Update(ctx context.Context, receipt *Domain.Receipt) error {

	query := `
		UPDATE receipts
		SET title = $1, author = $2, description = $3, upload_at = now()
		WHERE id = $4
		RETURNING id`

	args := []interface{}{
		receipt.Title,
		receipt.Author,
		receipt.Description,
		receipt.Id,
	}
	return r.db.QueryRow(ctx, query, args...).Scan(&receipt.Id)
}

func (r *receiptRepository) GetAll(ctx context.Context) ([]*Domain.Receipt, error) {
	query :=
		`SELECT id, title, author, description, created_at, upload_at
		FROM receipts
		ORDER BY id`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := []*Domain.Receipt{}

	for rows.Next() {
		var receipt Domain.Receipt

		err := rows.Scan(
			&receipt.Id,
			&receipt.Title,
			&receipt.Author,
			&receipt.Description,
			&receipt.Created_at,
			&receipt.Updated_at,
		)
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func NewReceiptRepository(db *pgxpool.Pool) ReceiptRepository {
	return &receiptRepository{db: db}
}
