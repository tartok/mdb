package mdb

import (
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Crud struct {
	dbName     string
	collection Collection
}
type UpdateStruct struct {
	Set         interface{}
	SetOnInsert interface{}
	Inc         interface{}
	Push        interface{}
	AddToSet    interface{}
	Unset       interface{}
	Raw         interface{}
	Upsert      bool
	NoLog       bool
	After       bool
}

func New(dbName string, coll Collection) *Crud {
	return &Crud{dbName: dbName, collection: coll}
}

func (c Crud) CollName() Collection {
	return c.collection
}
func (c Crud) Create(ctx GetContext, data interface{}) (*bson.ObjectID, error) {
	return create(ctx.GetContext(), c.dbName, c.collection, data)
}
func (c Crud) Update(ctx GetContext, filter interface{}, data UpdateStruct, res interface{}) error {
	t := time.Now()
	var set bson.D
	if !data.NoLog {
		set = append(set, bson.E{Key: "CU.uTime", Value: t}, bson.E{Key: "CU.uWho", Value: ctx.GetContext().LoginId})
	}
	if data.Set != nil {
		d, err := ToBsonD(data.Set)
		if err != nil {
			return err
		}
		for _, v := range d {
			switch v.Key {
			case "CU", "sys":
			default:
				set = append(set, v)
			}
		}
	}
	setOnInsert := bson.D{
		{"CU.cTime", t},
		{"CU.cWho", ctx.GetContext().LoginId},
	}
	if data.SetOnInsert != nil {
		d, err := ToBsonD(data.SetOnInsert)
		if err != nil {
			return err
		}
		for _, v := range d {
			switch v.Key {
			case "CU", "sys":
			default:
				setOnInsert = append(setOnInsert, v)
			}
		}
	}
	m := bson.D{
		{"$setOnInsert", setOnInsert},
	}
	if set != nil {
		m = append(m, bson.E{Key: "$set", Value: set})
	}
	if data.Inc != nil {
		m = append(m, bson.E{Key: "$inc", Value: data.Inc})
	}
	if data.Push != nil {
		m = append(m, bson.E{Key: "$push", Value: data.Push})
	}
	if data.AddToSet != nil {
		m = append(m, bson.E{Key: "$addToSet", Value: data.AddToSet})
	}
	if data.Unset != nil {
		m = append(m, bson.E{Key: "$unset", Value: data.Unset})
	}
	opts := options.FindOneAndUpdate()
	if data.Upsert {
		opts = opts.SetUpsert(true)
	}
	if data.Upsert || data.After {
		opts = opts.SetReturnDocument(options.After)
	}
	r := Coll(c.dbName, c.collection).FindOneAndUpdate(ctx.GetContext().Ctx, filter, m, opts)
	if r.Err() != nil {
		return r.Err()
	}
	if res != nil {
		err := r.Decode(res)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c Crud) UpdateMany(ctx GetContext, filter interface{}, data UpdateStruct) (*mongo.UpdateResult, error) {
	t := time.Now()
	var set bson.D
	if !data.NoLog {
		set = append(set, bson.E{Key: "CU.uTime", Value: t}, bson.E{Key: "CU.uWho", Value: ctx.GetContext().LoginId})
	}
	if data.Set != nil {
		d, err := ToBsonD(data.Set)
		if err != nil {
			return nil, err
		}
		for _, v := range d {
			switch v.Key {
			case "CU", "sys":
			default:
				set = append(set, v)
			}
		}
	}
	m := bson.D{
		{"$setOnInsert", bson.D{
			{"CU.cTime", t},
			{"CU.cWho", ctx.GetContext().LoginId},
		}},
	}
	if set != nil {
		m = append(m, bson.E{Key: "$set", Value: set})
	}
	if data.Inc != nil {
		m = append(m, bson.E{Key: "$inc", Value: data.Inc})
	}
	if data.Push != nil {
		m = append(m, bson.E{Key: "$push", Value: data.Push})
	}
	if data.Unset != nil {
		m = append(m, bson.E{Key: "$unset", Value: data.Unset})
	}
	return Coll(c.dbName, c.collection).UpdateMany(ctx.GetContext().Ctx, filter, m)
}
func (c Crud) Delete(ctx GetContext, filter interface{}, res interface{}) error {
	r := Coll(c.dbName, c.collection).FindOneAndDelete(ctx.GetContext().Ctx, filter)
	if r.Err() != nil {
		return r.Err()
	}
	if res != nil {
		err := r.Decode(res)
		if err != nil {
			return err
		}
	}
	return nil
}
func (c Crud) DeleteMany(ctx GetContext, filter interface{}) (int, error) {
	return DeleteMany(ctx.GetContext(), c.dbName, c.collection, filter)
}
func (c Crud) ReadMany(ctx GetContext, result interface{}, pipeline ...bson.A) error {
	var agg bson.A
	for _, a := range pipeline {
		agg = append(agg, a...)
	}
	return Aggregate(ctx, c.dbName, c.collection, agg, result)
}
func (c Crud) ReadOne(ctx GetContext, result interface{}, pipeline ...bson.A) error {
	var agg bson.A
	for _, a := range pipeline {
		agg = append(agg, a...)
	}
	v := reflect.ValueOf(result)
	a := reflect.New(reflect.SliceOf(v.Type()))
	r := a.Interface()
	err := Aggregate(ctx, c.dbName, c.collection, agg, r)
	if err != nil {
		return err
	}
	if a.Elem().Len() == 0 {
		return mongo.ErrNoDocuments
	}
	e := a.Elem().Index(0)
	v.Elem().Set(e.Elem())
	return nil
}
