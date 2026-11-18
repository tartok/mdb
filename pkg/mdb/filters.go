package mdb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

func BsonD(data ...any) (bson.D, error) {
	var res bson.D
	for _, i := range data {
		if i == nil {
			continue
		}
		var f bson.D
		switch d := i.(type) {
		case []byte:
			err := bson.UnmarshalExtJSON(d, true, &f)
			if err != nil {
				return nil, err
			}
		default:
			b, err := bson.Marshal(i)
			if err != nil {
				return nil, err
			}
			err = bson.Unmarshal(b, &f)
			if err != nil {
				return nil, err
			}
		}
		res = append(res, f...)
	}
	return res, nil
}
func FilterId(id bson.ObjectID, data ...any) any {
	if len(data) > 0 {
		var res bson.D
		res = append(res, bson.E{Key: "_id", Value: id})
		for _, i := range data {
			if i == nil {
				continue
			}
			b, err := bson.Marshal(i)
			if err != nil {
				continue
			}
			var f bson.D
			err = bson.Unmarshal(b, &f)
			if err != nil {
				continue
			}
			res = append(res, f...)
		}
		return res
	}
	return bson.M{"_id": id}
}
