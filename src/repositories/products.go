package repositories

import (
	"context"
	"database/sql"
	"eco/src/functions"
	"eco/src/models"
	"eco/src/types"
	"errors"
)

//RepoProducts é uma estrutura para inserir um usuário ao banco de dados
type RepoProducts struct {
	tx   *sql.Tx
	ctx  *context.Context
	conn *sql.Conn
}

//NewRepoProducts Inicia um novo repositorio de produtos
func NewRepoProducts(tx *sql.Tx, ctx *context.Context, conn *sql.Conn) *RepoProducts {
	return &RepoProducts{tx, ctx, conn}
}

//Insert é para inserir um novo produto ao banco de dados
func (db RepoProducts) Insert(p *models.Products) error {
	_, err := db.tx.ExecContext(
		*db.ctx,
		`INSERT INTO TB_PRODUCTS (PRO_SKU, PRO_TITLE, PRO_STATUS, COM_ID, BND_ID, PRO_DESCRIPTION, PRO_WEIGTH, PRO_SIZE, PRO_WIDTH, PRO_HEIGTH, 
			PRO_BARCODE, PRO_VARIATION, PRO_UPDATE, PRO_CROSSDOCKING, PRO_PRICE, PRO_PRICEPROMO, USR_ID, PDT_ID, PRO_REGISTER) 
		 VALUES (:0, :1, :2, :3, (SELECT BND_ID FROM BRANDS WHERE BND_CODE = :4 AND COM_ID = :5),
		         :6, :7, :8, :9 ,:10,
				:11, :12, SYSDATE, :13, 0, 0, :16, :17, SYSDATE) RETURNING PRO_ID INTO :18`,
		p.SKU, p.Title, functions.BoolToByte(p.Status), p.ComID, p.BndCode, p.ComID, p.Description, p.Weigth, p.Size, p.Width, p.Heigth,
		p.Barcode, p.Variation, p.Crossdocking, p.UsrID, p.PdtID, sql.Out{Dest: &p.ID},
	)
	return err
}

//Update altera um produto
func (db RepoProducts) Update(p *models.Products) error {
	res, err := db.tx.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS SET PRO_TITLE = :0, PRO_STATUS = :1, PRO_REGISTER = SYSDATE,
		BND_ID = (SELECT BND_ID FROM BRANDS WHERE BND_CODE = :2 AND COM_ID = :3), PRO_DESCRIPTION = :4, 
		PRO_WEIGTH = :5, PRO_SIZE = :6, PRO_WIDTH = :7, PRO_HEIGTH = :8,  PRO_BARCODE = :9, 
		PRO_UPDATE = SYSDATE, PRO_CROSSDOCKING = :11,
		USR_ID = :12, PDT_ID = :13 WHERE PRO_SKU = :14 AND COM_ID = :15 RETURNING PRO_ID INTO :16`,
		p.Title, functions.BoolToByte(p.Status), p.BndCode, p.ComID, p.Description, p.Weigth, p.Size,
		p.Width, p.Heigth, p.Barcode, p.Crossdocking, p.UsrID, p.PdtID, p.SKU, p.ComID, sql.Out{Dest: &p.ID},
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//ProUpdateDeactivate desativa um produto
func (db RepoProducts) ProUpdateDeactivate(SKU string, usrID uint32, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS SET PRO_STATUS = 0, PRO_UPDATE = SYSDATE, USR_ID = :0 WHERE PRO_SKU = :1 AND COM_ID = :2`,
		usrID, SKU, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//ProUpdateActivate ativa um produto
func (db RepoProducts) ProUpdateActivate(SKU string, usrID uint32, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`UPDATE TB_PRODUCTS SET PRO_STATUS = 1, PRO_UPDATE = SYSDATE, USR_ID = :0 WHERE PRO_SKU = :1 AND COM_ID = :2`,
		usrID, SKU, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

//Delete deleta um produto
func (db RepoProducts) Delete(SKU string, comID uint32) error {
	res, err := db.conn.ExecContext(
		*db.ctx,
		`DELETE FROM TB_PRODUCTS WHERE PRO_SKU = :0 AND COM_ID = :1 AND PRO_ID = PRO_VARIATION`,
		SKU, comID,
	)
	affected, err := res.RowsAffected()
	if affected == 0 {
		return errors.New(types.MsgNotUpdate)
	}
	return err
}

// //Query retornar um slice de produtos
// func (db RepoProducts) Query(title string, status string, SKU string, proType string, comID uint32, pag Pagination) (*[]models.Products, error) {
// 	var args []interface{}
// 	args = append(args, comID)
// 	var where string
// 	if title != "" {
// 		where += " AND REGEXP_LIKE(PRO_TITLE, :PRO_TITLE, 'i')"
// 		args = append(args, title)
// 	}

// 	if bStatus := functions.BoolStrToByte(status); bStatus != 10 {
// 		where += " AND PRO_STATUS = " + strconv.Itoa(int(bStatus))
// 	}

// 	if SKU != "" {
// 		where += " AND REGEXP_LIKE(PRO_SKU, :PRO_SKU, 'i')"
// 		args = append(args, SKU)
// 	}
// 	if proType != "" {
// 		where += " AND PDT_ID = :PDT_ID"
// 		args = append(args, proType)
// 	}
// 	strpag := pag.pag()
// 	rows, err := db.conn.QueryContext(
// 		*db.ctx,
// 		`SELECT P.PRO_SKU, P.PRO_TITLE, P.PRO_STATUS, P.PRO_DESCRIPTION, P.PRO_WEIGTH, P.PRO_SIZE, P.PRO_WIDTH, P.PRO_HEIGTH,
// 		        P.PRO_BARCODE, P.PRO_CROSSDOCKING, B.BND_CODE, P.PDT_ID, S.PRO_QTY, S.PRO_QTYBOOKING, CT.CTT_TITLE,
// 	     CURSOR(SELECT SP.SPE_VALUE, ST.SPT_CODE, ST.SPT_TITLE
// 			    FROM PRODUCTS_SPECIFICATIONS SP, SPECIFICATIONS_TITLES ST
// 			    WHERE SP.SPT_ID = ST.SPT_ID
// 			      AND SP.PRO_ID = P.PRO_ID) SPECIFICATIONS,
// 	     CURSOR(SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
// 			    FROM PRODUCTS_PRICES PC
// 			    WHERE PC.PRO_ID = P.PRO_ID) PRICES,
// 	     CURSOR(SELECT I.PRI_ORDER, I.PRI_URI
// 			    FROM PRODUCTS_IMAGES I
// 			    WHERE I.PRO_ID = P.PRO_ID) IMAGES,
// 	     CURSOR(SELECT PK.PRO_SKU, K.KIT_QTY
// 			    FROM PRODUCTS PK, KITS K
// 			    WHERE PK.PRO_ID = K.KIT_PRO_ID
// 			      AND K.PRO_ID = P.PRO_ID) KITS,
// 	     CURSOR(SELECT V.PRO_SKU, V.PRO_TITLE, V.PRO_STATUS, V.PRO_DESCRIPTION, V.PRO_CROSSDOCKING,
// 			    VS.PRO_QTY, VS.PRO_QTYBOOKING,
// 			    CURSOR(SELECT SP.SPE_VALUE, ST.SPT_CODE
// 				       FROM PRODUCTS_SPECIFICATIONS SP, SPECIFICATIONS_TITLES ST
// 				       WHERE SP.SPT_ID = ST.SPT_ID
// 					     AND SP.PRO_ID = V.PRO_ID) SPECIFICATIONS,
// 			    CURSOR(SELECT I.PRI_ORDER, I.PRI_URI
// 				       FROM PRODUCTS_IMAGES I
// 				       WHERE I.PRO_ID = V.PRO_ID) IMAGES,
// 			    CURSOR(SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
// 				       FROM PRODUCTS_PRICES PC
// 				       WHERE PC.PRO_ID = V.PRO_ID) PRICES
// 	            FROM PRODUCTS V, PRODUCTS_STOCKS VS
// 	            WHERE V.PRO_ID = VS.PRO_ID
// 	              AND V.PRO_VARIATION = P.PRO_ID
// 	              AND V.PRO_ID <> V.PRO_VARIATION)
//          FROM (SELECT P.PRO_SKU, P.PRO_TITLE, P.PRO_STATUS, P.PRO_DESCRIPTION, P.PRO_WEIGTH, P.PRO_SIZE, P.PRO_WIDTH, P.PRO_HEIGTH, P.PRO_BARCODE, P.PRO_CROSSDOCKING,
// 		            P.PDT_ID, P.PRO_VARIATION, P.PRO_ID, P.BND_ID, P.COM_ID
//                     FROM PRODUCTS P `+strpag+`) P, BRANDS B, PRODUCTS_STOCKS S
//         WHERE P.BND_ID = B.BND_ID
//           AND P.PRO_ID = S.PRO_ID
//           AND P.PRO_ID = P.PRO_VARIATION
// 		  AND P.COM_ID = :COM_ID`+where+" ORDER BY P.PRO_TITLE", args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var products []models.Products
// 	for rows.Next() {
// 		var pro models.Products
// 		var cat *sql.Rows
// 		var spe *sql.Rows
// 		var pri *sql.Rows
// 		var ima *sql.Rows
// 		var kits *sql.Rows
// 		var variations *sql.Rows
// 		err = rows.Scan(&pro.SKU, &pro.Title, &pro.Status, &pro.Description, &pro.Weigth, &pro.Size, &pro.Width, &pro.Heigth,
// 			&pro.Barcode, &pro.Crossdocking, &pro.BndCode, &pro.PdtID, &pro.Stock.Qty, &pro.Stock.QtyBooking,
// 			&cat, &spe, &pri, &ima, &kits, &variations)
// 		for cat.Next() {
// 			var category string
// 			cat.Scan(&category)
// 			pro.Categorys = append(pro.Categorys, category)
// 		}
// 		cat.Close()

// 		for spe.Next() {
// 			var specification models.ProductsSpecifications
// 			spe.Scan(&specification.Value, &specification.Code, &specification.Title)
// 			pro.Specifications = append(pro.Specifications, specification)
// 		}
// 		spe.Close()
// 		for pri.Next() {
// 			var price models.ProductsPrices
// 			pri.Scan(&price.Price, &price.PricePromo, &price.Status, &price.MkcID)
// 			pro.Prices = append(pro.Prices, price)
// 		}
// 		pri.Close()

// 		for ima.Next() {
// 			var image models.ProductsImages
// 			ima.Scan(&image.Order, &image.URI)
// 			pro.Images = append(pro.Images, image)
// 		}
// 		ima.Close()
// 		for kits.Next() {
// 			var kit models.Kits
// 			kits.Scan(&kit.KitSKU, &kit.Qty)
// 			pro.Kits = append(pro.Kits, kit)
// 		}
// 		kits.Close()

// 		for variations.Next() {
// 			var vari models.ProductsVariations
// 			var variSpe *sql.Rows
// 			var variImages *sql.Rows
// 			var variPrices *sql.Rows
// 			variations.Scan(&vari.SKU, &vari.Title, &vari.Status, &vari.Description, &vari.Crossdocking, &vari.Stock.Qty,
// 				&vari.Stock.QtyBooking, &variSpe, &variImages, &variPrices)

// 			for variSpe.Next() {
// 				var specification models.ProductsSpecifications
// 				variSpe.Scan(&specification.Value, &specification.Code)
// 				vari.Specifications = append(vari.Specifications, specification)
// 			}
// 			variSpe.Close()

// 			for variImages.Next() {
// 				var image models.ProductsImages
// 				variImages.Scan(&image.Order, &image.URI)
// 				vari.Images = append(vari.Images, image)
// 			}
// 			variImages.Close()

// 			for variPrices.Next() {
// 				var price models.ProductsPrices
// 				variPrices.Scan(&price.Price, &price.PricePromo, &price.Status, &price.MkcID)
// 				vari.Prices = append(vari.Prices, price)
// 			}
// 			variPrices.Close()
// 			pro.Variations = append(pro.Variations, vari)
// 		}
// 		variations.Close()

// 		products = append(products, pro)
// 	}
// 	return &products, nil
// }

// //QueryByID traz apenas um registro no select
// func (db RepoProducts) QueryByID(SKU string, comID uint32) (*models.Products, error) {

// 	row, err := db.conn.QueryContext(
// 		*db.ctx,
// 		`SELECT P.PRO_SKU, P.PRO_TITLE, P.PRO_STATUS, P.PRO_DESCRIPTION, P.PRO_WEIGTH, P.PRO_SIZE, P.PRO_WIDTH, P.PRO_HEIGTH,
// 	       	    P.PRO_BARCODE, P.PRO_CROSSDOCKING, B.BND_CODE, P.PDT_ID, S.PRO_QTY, S.PRO_QTYBOOKING, CT.CTT_TITLE,
// 	            CURSOR(SELECT SP.SPE_VALUE, ST.SPT_CODE, ST.SPT_TITLE
// 			           FROM PRODUCTS_SPECIFICATIONS SP, SPECIFICATIONS_TITLES ST
// 			           WHERE SP.SPT_ID = ST.SPT_ID
// 			             AND SP.PRO_ID = P.PRO_ID) SPECIFICATIONS,
// 	            CURSOR(SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
// 			           FROM PRODUCTS_PRICES PC
// 			           WHERE PC.PRO_ID = P.PRO_ID) PRICES,
// 	            CURSOR(SELECT I.PRI_ORDER, I.PRI_URI
// 			           FROM PRODUCTS_IMAGES I
// 			           WHERE I.PRO_ID = P.PRO_ID) IMAGES,
// 	            CURSOR(SELECT PK.PRO_SKU, K.KIT_QTY
// 			           FROM PRODUCTS PK, KITS K
// 			           WHERE PK.PRO_ID = K.KIT_PRO_ID
// 			             AND K.PRO_ID = P.PRO_ID) KITS,
// 	            CURSOR(SELECT V.PRO_SKU, V.PRO_TITLE, V.PRO_STATUS, V.PRO_DESCRIPTION, V.PRO_CROSSDOCKING,
// 			           VS.PRO_QTY, VS.PRO_QTYBOOKING,
// 			           CURSOR(SELECT SP.SPE_VALUE, ST.SPT_CODE
// 				              FROM PRODUCTS_SPECIFICATIONS SP, SPECIFICATIONS_TITLES ST
// 				              WHERE SP.SPT_ID = ST.SPT_ID
// 					            AND SP.PRO_ID = V.PRO_ID) SPECIFICATIONS,
// 			           CURSOR(SELECT I.PRI_ORDER, I.PRI_URI
// 				              FROM PRODUCTS_IMAGES I
// 					          WHERE I.PRO_ID = V.PRO_ID) IMAGES,
// 			           CURSOR(SELECT PC.PRP_PRICE, PC.PRP_PRICEPROMO, PC.PRP_STATUS, PC.MKC_ID
// 				              FROM PRODUCTS_PRICES PC
// 				              WHERE PC.PRO_ID = V.PRO_ID) PRICES
// 	            FROM PRODUCTS V, PRODUCTS_STOCKS VS
// 	            WHERE V.PRO_ID = VS.PRO_ID
// 		          AND V.PRO_VARIATION = P.PRO_ID
// 		          AND V.PRO_ID <> V.PRO_VARIATION)
//         FROM PRODUCTS P, BRANDS B, PRODUCTS_STOCKS S, CATEGORYS_TITLES CT, CATEGORYS_NICKNAMES N
//         WHERE P.BND_ID = B.BND_ID
//           AND P.PRO_ID = S.PRO_ID
//           AND P.PRO_ID = P.PRO_VARIATION
//           AND P.CNN_ID = N.CNN_ID
//           AND CT.CTT_ID = N.CTT_ID
// 		  AND P.COM_ID = :COM_ID
// 		  AND P.PRO_SKU = :PRO_SKU`, comID, SKU)

// 	var pro models.Products
// 	var spe *sql.Rows
// 	var pri *sql.Rows
// 	var ima *sql.Rows
// 	var kits *sql.Rows
// 	var variations *sql.Rows

// 	if row.Next() {
// 		err = row.Scan(&pro.SKU, &pro.Title, &pro.Status, &pro.Description, &pro.Weigth, &pro.Size, &pro.Width, &pro.Heigth,
// 			&pro.Barcode, &pro.Crossdocking, &pro.BndCode, &pro.PdtID, &pro.Stock.Qty, &pro.Stock.QtyBooking,
// 			&pro.Category, &spe, &pri, &ima, &kits, &variations)
// 	} else {
// 		return nil, errors.New(types.MsgNotFound)
// 	}
// 	defer row.Close()
// 	if err != nil {
// 		return nil, err
// 	}
// 	for spe.Next() {
// 		var specification models.ProductsSpecifications
// 		spe.Scan(&specification.Value, &specification.Code, &specification.Title)
// 		pro.Specifications = append(pro.Specifications, specification)
// 	}
// 	spe.Close()
// 	for pri.Next() {
// 		var price models.ProductsPrices
// 		pri.Scan(&price.Price, &price.PricePromo, &price.Status, &price.MkcID)
// 		pro.Prices = append(pro.Prices, price)
// 	}
// 	pri.Close()

// 	for ima.Next() {
// 		var image models.ProductsImages
// 		ima.Scan(&image.Order, &image.URI)
// 		pro.Images = append(pro.Images, image)
// 	}
// 	ima.Close()
// 	for kits.Next() {
// 		var kit models.Kits
// 		kits.Scan(&kit.KitSKU, &kit.Qty)
// 		pro.Kits = append(pro.Kits, kit)
// 	}
// 	kits.Close()

// 	for variations.Next() {
// 		var vari models.ProductsVariations
// 		var variSpe *sql.Rows
// 		var variImages *sql.Rows
// 		var variPrices *sql.Rows
// 		variations.Scan(&vari.SKU, &vari.Title, &vari.Status, &vari.Description, &vari.Crossdocking, &vari.Stock.Qty,
// 			&vari.Stock.QtyBooking, &variSpe, &variImages, &variPrices)

// 		for variSpe.Next() {
// 			var specification models.ProductsSpecifications
// 			variSpe.Scan(&specification.Value, &specification.Code)
// 			vari.Specifications = append(vari.Specifications, specification)
// 		}
// 		variSpe.Close()

// 		for variImages.Next() {
// 			var image models.ProductsImages
// 			variImages.Scan(&image.Order, &image.URI)
// 			vari.Images = append(vari.Images, image)
// 		}
// 		variImages.Close()

// 		for variPrices.Next() {
// 			var price models.ProductsPrices
// 			variPrices.Scan(&price.Price, &price.PricePromo, &price.Status, &price.MkcID)
// 			vari.Prices = append(vari.Prices, price)
// 		}
// 		variPrices.Close()
// 		pro.Variations = append(pro.Variations, vari)
// 	}
// 	variations.Close()

// 	return &pro, nil
// }

// // SELECT V.PRO_SKU, V.PRO_TITLE, V.PRO_STATUS, V.PRO_DESCRIPTION, V.PRO_CROSSDOCKING,
// //VS.PRO_QTY, VS.PRO_QTYBOOKING
