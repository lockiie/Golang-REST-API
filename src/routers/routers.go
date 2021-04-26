package routers

import (
	ctrls "eco/src/controllers"
	"eco/src/db"
	"eco/src/types"

	"github.com/gofiber/fiber/v2"
)

const (
	brands         = "/brands"
	categorys      = "/categorys"
	nicks          = "/nicks"
	specifications = "/specifications"
	products       = "/products"
	stock          = "/stock"
	activate       = "/activate"
	deactivate     = "/deactivate"
	marketplace    = "/marketplace"
	price          = "/price"
	prices         = "/prices"
	images         = "/images"
	variations     = "/variations"
)

func init() {
	app := fiber.New()

	v0 := app.Group("/v0")
	defer db.Pool.Close()
	// v0.All("", func(c *fiber.Ctx) error {
	// 	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	// 	return nil
	// })

	v0.Post(brands, ctrls.CreateBrands)                 //brands
	v0.Put(brands+types.ParamID, ctrls.UpdateBrands)    //brands/:id
	v0.Delete(brands+types.ParamID, ctrls.DeleteBrands) //brands/:id
	v0.Get(brands, ctrls.QueryBrands)                   //brands?id&status&name
	v0.Get(brands+types.ParamID, ctrls.QueryBrandsByID) ///brands/:id

	v0.Post(categorys, ctrls.CreateCategorys)                 //categorys
	v0.Put(categorys+types.ParamID, ctrls.UpdateCategorys)    //categorys/:id
	v0.Delete(categorys+types.ParamID, ctrls.DeleteCategorys) //categorys/:id
	v0.Get(categorys+nicks, ctrls.QueryCategorysNiks)         //categorys?id&name
	v0.Get(categorys, ctrls.QueryCategorys)                   //categorys?id&name
	//v0.Get(categorys+types.ParamID, ctrls.QueryCategorysByID) ///categorys/:id

	v0.Post(specifications, ctrls.CreateSpecifications)                 //specifications
	v0.Put(specifications+types.ParamID, ctrls.UpdateSpecifications)    //specifications/:id
	v0.Delete(specifications+types.ParamID, ctrls.DeleteSpecifications) //specifications/:id
	v0.Get(specifications, ctrls.QuerySpecifications)                   //specifications?id&title&status
	v0.Get(specifications+types.ParamID, ctrls.QuerySpecificationsByID) ///specifications/:id

	/////////////////////////////////////////////////////////////////////////////////////////NADA FEITO AINDA
	//Produtos
	//pro := v0.Group(products)

	v0.Post(products, ctrls.CreateProducts)                  //products
	v0.Put(products+types.ParamSKU, ctrls.UpdateProducts)    //products/:sku
	v0.Delete(products+types.ParamSKU, ctrls.DeleteProducts) //products/:sku
	// v0.Get(products, ctrls.QueryProducts)                    //products
	// v0.Get(products+types.ParamSKU, ctrls.QueryProductsByID) ///products/:sku

	v0.Put(products+types.ParamSKU+prices, ctrls.UpdateProductsPrice) //products/:sku/price
	v0.Get(products+types.ParamSKU+prices, ctrls.QueryProductsPrices) //products/:sku/price

	v0.Put(products+types.ParamSKU+stock, ctrls.UpdateProductsStock) //products/:sku/stock
	v0.Get(products+types.ParamSKU+stock, ctrls.QueryProductsStocks) //products/:sku/stock

	v0.Put(products+types.ParamSKU+activate, ctrls.UpdateProductsActivate)     //products/:sku/activate
	v0.Put(products+types.ParamSKU+deactivate, ctrls.UpdateProductsDeactivate) //products/:sku/deactivate

	//Especificação dos produtos
	v0.Post(products+types.ParamSKU+specifications, ctrls.CreateProductsSpecifications)                 //products/:sku/specifications
	v0.Put(products+types.ParamSKU+specifications+types.ParamID, ctrls.UpdateProductsSpecifications)    //products/:sku/specifications/:id
	v0.Delete(products+types.ParamSKU+specifications+types.ParamID, ctrls.DeleteProductsSpecifications) //products/:sku/specifications/:id
	v0.Get(products+types.ParamSKU+specifications, ctrls.QueryProductsSpecifications)                   //products/:sku/specifications
	//v0.Get(products+types.ParamSKU+specifications+types.ParamID, ctrls.QuerySpecificationsByID) //products/:sku/specifications/:id

	//Products Categorias
	v0.Post(products+types.ParamSKU+categorys+types.ParamID, ctrls.CreateProductsCategorys)   //products/:sku/categorys
	v0.Delete(products+types.ParamSKU+categorys+types.ParamID, ctrls.DeleteProductsCategorys) //products/:sku/categorys/:id
	// v0.Get(products+types.ParamSKU+categorys, ctrls.QueryProductsCategorys)                   //products/:sku/categorys
	// v0.Get(products+types.ParamSKU+categorys+types.ParamID, ctrls.QuerySpecificationsByID) //products/:sku/categorys/:id

	// //Produtos com imagem
	v0.Post(products+types.ParamSKU+images, ctrls.CreateProductsImages)                 //products/:sku/images
	v0.Put(products+types.ParamSKU+images+types.ParamID, ctrls.UpdateProductsImages)    //products/:sku/images/:id
	v0.Delete(products+types.ParamSKU+images+types.ParamID, ctrls.DeleteProductsImages) //products/:sku/images/:id
	v0.Get(products+types.ParamSKU+images, ctrls.QueryProductsImages)                   //products/:sku/images
	// v0.Get(products+types.ParamSKU+images+types.ParamID, ctrls.QuerySpecificationsByID) //products/:sku/images/:id
	v0.Post(products+types.ParamSKU+variations, ctrls.CreateProductsVariations)                 //products/:sku/variations
	v0.Put(products+types.ParamSKU+variations+types.ParamID, ctrls.UpdateProductsVariations)    //products/:sku/variations/:id
	v0.Delete(products+types.ParamSKU+variations+types.ParamID, ctrls.DeleteProductsVariations) //products/:sku/variations/:id
	v0.Get(products+types.ParamSKU+variations, ctrls.QueryProductsVariations)                   //products/:sku/variations

	// data, _ := json.MarshalIndent(app.Stack(), "", "  ")
	// fmt.Println(string(data))

	app.Listen(":4000")
}
