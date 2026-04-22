package service

import (
	jwt "pet-care/internal/JWT"
	Store "pet-care/store"
)

type Services struct {
	StoreDB   Store.Storage
	TokenUtil jwt.TokenUtil
}
