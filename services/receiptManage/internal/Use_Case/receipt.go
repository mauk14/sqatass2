package Use_Case

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	validator2 "messanger/pkg/validator"
	"messanger/services/receiptManage/internal/Domain"
	"messanger/services/receiptManage/internal/Repository"
	"unicode/utf8"
)

type ReceiptUseCase interface {
	Create(context.Context, *Domain.Receipt) error
	Delete(context.Context, int64) error
	Update(context.Context, int64, *Domain.Receipt) error
	Get(context.Context, int64) (*Domain.Receipt, error)
	GetAll(context.Context) ([]*Domain.Receipt, error)
}

func ValidateReceipt(v *validator2.Validator, receipt *Domain.Receipt) {
	v.Check(receipt.Title != "", "title", "must be provided")
	v.Check(utf8.RuneCountInString(receipt.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(receipt.Author != "", "author", "must be provided")
	v.Check(utf8.RuneCountInString(receipt.Author) <= 500, "author", "must not be more than 500 bytes long")
	v.Check(receipt.Description != "", "description", "must be provided")
}

type receiptUseCase struct {
	repository Repository.ReceiptRepository
}

func (r *receiptUseCase) Create(ctx context.Context, receipt *Domain.Receipt) error {
	v := validator2.New()

	if ValidateReceipt(v, receipt); !v.Valid() {
		return errors.New("Error validation response")
	}

	return r.repository.Create(ctx, receipt)
}

func (r *receiptUseCase) Delete(ctx context.Context, id int64) error {
	if id < 1 {
		return errors.New("Record not found")
	}

	return r.repository.Delete(ctx, id)

}

func (r *receiptUseCase) Update(ctx context.Context, id int64, receipt *Domain.Receipt) error {
	if id < 1 {
		return errors.New("Record not found")
	}

	rec, err := r.Get(ctx, id)
	//fmt.Println(err)

	if err != nil {
		return err
	}

	if receipt.Title != "" {
		rec.Title = receipt.Title
	}

	if receipt.Author != "" {
		rec.Author = receipt.Author
	}
	if receipt.Description != "" {
		rec.Description = receipt.Description
	}

	return r.repository.Update(ctx, rec)

}
func (r *receiptUseCase) Get(ctx context.Context, id int64) (*Domain.Receipt, error) {
	if id < 1 {
		return nil, errors.New("Record not found")
	}
	receipt, err := r.repository.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("Record not found")
		}
		return nil, err
	}
	return receipt, nil
}
func (r *receiptUseCase) GetAll(ctx context.Context) ([]*Domain.Receipt, error) {
	receipts, err := r.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func NewReceiptUseCase(rep Repository.ReceiptRepository) ReceiptUseCase {
	return &receiptUseCase{
		repository: rep,
	}
}
