package types

const (
	//EmptyStr representa uma string vazia
	EmptyStr = ""
	//MaxNumber1 Máximo de números 1
	MaxNumber1 = 10
	//MaxNumber2 Máximo de números 2
	MaxNumber2 = 100
	//MaxNumber3 máximo de números 3
	MaxNumber3 = 1000
	//MaxNumber4 máximo de números 4
	MaxNumber4 = 10000
	//MaxNumber5 máximo de números 5
	MaxNumber5 = 100000
	//MaxNumber6 máximo de números 6
	MaxNumber6 = 1000000
	//CodeErrorValidate é o código do erro da validação

	//Parametros
	//Limit para o gets dos selects
	Limit = "_limit"
	//OffSet para o gets dos selects
	OffSet = "_offSet"

	///errros
	CodeErrValidate uint16 = 2
	//CodeErrJSONDecode é o erro quando o decode do json falha
	CodeErrJSONDecode uint16 = 3
	//CodeErrPoolDB é o erro quando não foi possível capturar uma pool do banco de dados
	CodeErrPoolDB uint16 = 4
	//CodeErrInsertDB é um erro de falha ao inserir algo no banco de dados
	CodeErrInsertDB uint16 = 5
	//CodeErrUpdateDB é um erro de falha ao alterar um registro no banco de dados
	CodeErrUpdateDB uint16 = 6
	//CodeErrDeleteDB é um erro de falha ao remover um item do banco de dados
	CodeErrDeleteDB uint16 = 7
	//CodeErrQueryDB falha ao realizar um select no banco de dados
	CodeErrQueryDB uint16 = 8
	//CodeErrReadParam parametro da requisição inválido
	CodeErrReadParam uint16 = 9
	//CodeErrActionNotAllowed ação inválida
	CodeErrActionNotAllowed uint16 = 10

	//MsgNotUpdate é a mensagem que nenhum registro foi alterado
	MsgNotUpdate = "Nenhum registro foi alterado ou deletado"
	//MsgNotFound
	MsgNotFound = "Nenhum registro encontrado"
)
