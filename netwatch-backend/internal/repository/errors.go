package repository

import "errors"

var (
	ErrNotFound   = errors.New("registro não encontrado")
	ErrDuplicate  = errors.New("registro duplicado")
	ErrForeignKey = errors.New("violação de chave estrangeira")
)
