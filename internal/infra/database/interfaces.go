package database

import "github.com/Guilherme-daCosta/goexpert-api/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
type ProductInteface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}