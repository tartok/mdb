package mdb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Context struct {
	Ctx     context.Context
	LoginId *bson.ObjectID
}
