package repositories

import (
	"strconv"
)

//Pagination é a estrutura para a paginação da querys
type Pagination struct {
	OffSet string
	Limit  string
}

func (pag Pagination) pag() string {
	offSet, err := strconv.Atoi(pag.OffSet)
	if err != nil {
		offSet = 0
	}
	limit, err := strconv.Atoi(pag.Limit)
	if err != nil {
		limit = 20
	}
	return " OFFSET " + strconv.Itoa(offSet) + " ROWS FETCH NEXT " + strconv.Itoa(limit) + " ROWS ONLY"
}
