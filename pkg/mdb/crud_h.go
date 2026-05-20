package mdb

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ICalculateFunc interface {
	CalculateFunc(ctx Context)
}

func ToBsonD(data interface{}) (bson.D, error) {
	d, err := bson.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m bson.D
	err = bson.Unmarshal(d, &m)
	return m, err
}
func DeleteMany(ctx Context, dbName string, coll Collection, filter interface{}) (int, error) {
	r, err := Coll(dbName, coll).DeleteMany(ctx.Ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(r.DeletedCount), nil
}

func create(ctx Context, dbName string, coll Collection, data interface{}) (*bson.ObjectID, error) {
	m := bson.D{{"CU", CU{
		CUTime: CUTime{
			CTime: CTime{CTime: Ref(time.Now())},
		},
		CUWho: CUWho{CWho: CWho{CWho: ctx.LoginId}},
	}}}
	d, err := ToBsonD(data)
	if err != nil {
		return nil, err
	}
	for _, v := range d {
		switch v.Key {
		case "CU":
		default:
			m = append(m, v)
		}
	}
	id, err := Coll(dbName, coll).InsertOne(ctx.Ctx, m)
	if err != nil {
		return nil, err
	}
	newId := id.InsertedID.(bson.ObjectID)
	return &newId, nil
}
func Aggregate(ctx Context, dbName string, collName Collection, pipeline interface{}, result interface{}) error {
	err := aggregate(ctx, dbName, collName, pipeline, result)
	calcFields(ctx, result)
	return err
}
func aggregate(ctx Context, dbName string, collName Collection, pipeline interface{}, result interface{}) error {
	c, err := Coll(dbName, collName).Aggregate(ctx.Ctx, pipeline)
	if err != nil {
		return err
	}
	return c.All(ctx.Ctx, result)
}
func calcFields(ctx Context, cf interface{}) {
	v := reflect.Indirect(reflect.ValueOf(cf))
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if f.CanInterface() && !f.IsZero() {
				if f.Kind() != reflect.Pointer && f.CanAddr() {
					f = f.Addr()
				}
				calcFields(ctx, f.Interface())
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			f := v.Index(i)
			if f.CanInterface() && !f.IsZero() {
				if f.Kind() != reflect.Pointer && f.CanAddr() {
					f = f.Addr()
				}
				calcFields(ctx, f.Interface())
			}
		}
	default:
	}
	if c, ok := cf.(ICalculateFunc); ok {
		c.CalculateFunc(ctx)
	}
}
