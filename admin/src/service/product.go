package service

import (
	"blitzSeckill/admin/src/db"
	"blitzSeckill/admin/src/model"
	"log"
)

type ProductService struct {
}

func NewProductServer() *ProductService {
	return &ProductService{}
}

func (p *ProductService) UpdateStockNum(product *model.Product, num int) error {
	//更新库存
	err := model.ProductNsp.UpdateStockNum(db.DB, product, num)
	if err != nil {
		log.Printf("ProductNsp.UpdateStockNum, err : %v", err)
		return err
	}
	return nil
}

//func (p *ProductServer) CreateProduct(product *model.Product) error {
//	productEntity := model.NewProductModel()
//	err := productEntity.CreateProduct(product)
//	if err != nil {
//		log.Printf("ProductEntity.CreateProduct, err : %v", err)
//		return err
//	}
//	return nil
//}

//func (p *ProductServer) GetProductList() ([]map[string]interface{}, error) {
//	productEntity := model.NewProductModel()
//	productList, err := productEntity.GetProductList()
//	if err != nil {
//		log.Printf("ProductEntity.CreateProduct, err : %v", err)
//		return nil, err
//	}
//	return productList, nil
//}
