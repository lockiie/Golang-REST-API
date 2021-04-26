package types

//ProAction é um type para saber se vai inserir ou alterar
type ProAction byte

const (
	//Insert é para inserir um produto
	Insert ProAction = 0
	//Update é para alterar o produto
	Update ProAction = 1
)

//ProType é para identificar o tipo do produto
type ProType byte

const (
	//Normal é um produto normal
	Normal ProType = 1
	//Variation é um produto variado
	Variation ProType = 2
	//Kit é  tipo do produto kit
	Kit ProType = 3
)
