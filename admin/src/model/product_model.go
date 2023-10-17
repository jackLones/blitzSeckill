package model

import "github.com/jinzhu/gorm"

var ProductNsp Product

type Product struct {
	ProductId   uint32 `gorm:"column:product_id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'商品Id'"`
	ProductName string `gorm:"column:product_name;default:;NOT NULL;comment:'商品名称'"`
	Total       uint32 `gorm:"column:total;default:0;NOT NULL;comment:'商品数量'"`
	Status      uint8  `gorm:"column:status;default:0;NOT NULL;comment:'商品状态'"`
}

func (p *Product) TableName() string {
	return "product"
}

// Find 查找记录
func (p *Product) Find(db *gorm.DB, product_id string) (*Product, error) {
	var data = &Product{}

	err := db.Table(p.TableName()).Where("product_id = ?", product_id).First(data).Error
	return data, err
}

// Create 创建记录
func (p *Product) Create(db *gorm.DB, product *Product) error {
	err := db.Table(p.TableName()).Create(product).Error
	return err
}

// Save 保存记录
func (p *Product) Save(db *gorm.DB, product *Product) error {
	err := db.Table(p.TableName()).Save(product).Error
	return err
}

func (p *Product) UpdateStockNum(db *gorm.DB, product *Product, num int) error {

	result := db.Model(&product).UpdateColumn("total", gorm.Expr("total - ?", num))
	// UPDATE "products" SET "quantity" = quantity - 1 WHERE "id" = 3;
	return result.Error
}
