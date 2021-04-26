## Eco


git branch -M main  -->  git remote add origin https://github.com/lockiie/Eco_Api.git  --> git push -u origin main


###### ROTAS

| Marcas | URI |
|:----------:|:-------------:|
| POST   |  v0/brands      |
| PUT    |  v0/brands/:id  |  
| DELETE |  v0/brands/:id  |
| GET    |  v0/brands      |
| GET    |  v0/brands/:id  |

<br />

| Categorias | URI |
|:----------:|:-------------:|
| POST   |    v0/categorys                  |
| PUT    |    v0/categorys/:id              |
| DELETE |    v0/categorys/:id              |
| GET    |    v0/categorys?id&title&status  |
| GET    |    v0/categorys/:id              |

<br />

| Especificações | URI |
|:----------:|:-------------:|
| POST   |  v0/specifications                 |
| PUT    |  v0/specifications/:id             |
| DELETE |  v0/specifications/:id             |
| GET    |  v0/specifications?id&title&status |
| GET    |  v0/specifications/:id             |

<br />

| Produtos | URI |
|:----------:|:-------------:|
| POST  |  v0/products                                    |
| PUT   |  v0/products/:sku                               |
| DELETE|  v0/products/:sku                               |
| GET   |  v0/products?id&title&status&type_product       |
| GET   |  v0/products/:id                                |

<br />

| Alterações dos Produtos | URI |
|:----------:|:-------------:|
| PUT    |    v0/products/:sku/prices     |
| GET    |    v0/products/:sku/prices     |
| PUT    |    v0/products/:sku/stock      |
| GET    |    v0/products/:sku/stock      |
| PUT    |    v0/products/:sku/activate   |
| PUT    |    v0/products/:sku/deactivate |  

<br />

| Espicificações dos produtos | URI |
|:----------:|:-------------:|
| POST   | v0/products/:sku/specifications      |
| PUT    | v0/products/:sku/specifications/:id  |
| DELETE | v0/products/:sku/specifications/:id  |
| GET    | v0/products/:sku/specifications      |

<br />

| Categorias dos produtos | URI |
|:----------:|:-------------:|
| POST   | v0/products/:sku/categorys/:id  |
| DELETE | v0/products/:sku/categorys/:id  |
| GET   | v0/products/:sku/categorys/      |

<br />

| Imagens dos produtos | URI |
|:----------:|:-------------:|
| POST   |  v0/products/:sku/images           |
| PUT    |  v0/products/:sku/images/:order    |
| DELETE |  v0/products/:sku/images/:order    |
| GET    |  v0/products/:sku/images           |

<br />

| Variações dos produtos | URI |
|:----------:|:-------------:|
| POST   | v0/products/:sku/variations       |
| PUT    | v0/products/:sku/variations/:id   |
| DELETE | v0/products/:sku/variations/:id   |
| GET    | v0/products/:sku/variations       |