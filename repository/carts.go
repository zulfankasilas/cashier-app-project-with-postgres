package repository

import (
	"a21hc3NpZ25tZW50/model"
	"errors"
	"math"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	results := []model.JoinCart{}
	err := c.db.Table("carts").Select("carts.id, carts.product_id, products.name, carts.quantity, carts.total_price").Joins("left join products on products.id = carts.product_id").Scan(&results)

	if err.Error != nil {
		return []model.JoinCart{}, err.Error
	}

	return results, nil
}

func (c *CartRepository) AddCart(product model.Product) error {
	cart := model.Cart{}

	c.db.Transaction(func(tx *gorm.DB) error {
		data := tx.Where("product_id = ?", product.ID).First(&cart)

		if errors.Is(data.Error, gorm.ErrRecordNotFound) {
			if err := tx.Create(&model.Cart{ProductID: product.ID, Quantity: 0, TotalPrice: 0}).Error; err != nil {
				return err
			}

			if err := tx.Where("product_id = ?", product.ID).First(&cart).Error; err != nil {
				return err
			}
		}

		newCart := model.Cart{ProductID: product.ID, Quantity: cart.Quantity + 1, TotalPrice: cart.TotalPrice + (product.Price * ((100 - product.Discount) / 100.0))}

		if err := tx.Model(&model.Cart{}).Where("product_id = ?", product.ID).Updates(&newCart).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", product.ID).Update("stock", product.Stock-1).Error; err != nil {
			return err
		}

		// return nil will commit the whole transaction
		return nil
	})

	return nil
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {
	cart := model.Cart{}
	product := model.Product{}

	c.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&cart).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", id).Delete(&model.Cart{}).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}

		if err := tx.Model(&model.Product{}).Where("id = ?", productID).Update("stock", product.Stock+int(math.Round(cart.Quantity))).Error; err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {
	err := c.db.Where("id = ?", id).Updates(&cart)

	if err.Error != nil {
		return err.Error
	}

	return nil
}
