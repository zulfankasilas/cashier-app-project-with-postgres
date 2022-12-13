package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	data := p.db.Create(&product)

	if data.Error != nil {
		return data.Error
	}

	return nil
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	products := []model.Product{}

	if err := p.db.Find(&products).Error; err != nil {
		return []model.Product{}, err
	}

	return products, nil
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	if err := p.db.Where("id = ?", id).Delete(&model.Product{}).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {

	data := p.db.Where("id = ?", id).Updates(&product)

	if data.Error != nil {
		return data.Error
	}

	return nil
}
