package mdb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type GetContext interface {
	GetContext() Context
}
type Context struct {
	Ctx     context.Context
	LoginId *bson.ObjectID
}

func (c Context) GetContext() Context {
	return c
}

func NewContext(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{
		Ctx:     ctx,
		LoginId: nil,
	}
}
