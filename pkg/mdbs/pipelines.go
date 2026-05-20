package mdbs

import (
	"mdb/pkg/mdb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CollName interface {
	CollName() mdb.Collection
}

func AddFields(fields bson.D) (res bson.A) {
	var ff bson.D
	ff = append(ff, fields...)
	return bson.A{bson.D{
		{"$addFields", ff}},
	}

}
func Project(fields bson.D) (res bson.A) {
	var ff bson.D
	ff = append(ff, fields...)
	return bson.A{bson.D{
		{"$project", ff}},
	}

}

func If(cond bool, f func() bson.A) bson.A {
	if cond {
		return f()
	}
	return nil
}
func MatchId(id bson.ObjectID) (res bson.A) {
	return Match(bson.D{{"_id", id}})
}
func MatchEmail(email string) (res bson.A) {
	return Match(bson.D{{"email", email}})
}
func MatchPairs(m ...interface{}) (res bson.A) {
	i := 0
	var r bson.D
	for i < len(m)-1 {
		if f, ok := m[i].(string); ok {
			i++
			r = append(r, bson.E{f, m[i]})
		}
		i++
	}
	return Match(r)
}
func Match(m ...interface{}) (res bson.A) {
	if len(m) == 0 {
		return
	}
	d := bson.D{}
	var f func(m interface{})
	f = func(m interface{}) {
		if m == nil {
			return
		}
		switch r := m.(type) {
		case bson.E:
			d = append(d, r)
		case bson.D:
			for _, e := range r {
				d = append(d, e)
			}
		case bson.A:
			for _, i2 := range r {
				f(i2)
			}
		case bson.M:
			for s, i := range r {
				d = append(d, bson.E{s, i})
			}
		case []interface{}:
			for _, i2 := range r {
				f(i2)
			}
		default:
			res = append(res, bson.D{{"$match", r}})
		}
	}
	f(m)
	if len(res) > 0 {
		return res
	}
	return append(res, bson.D{{"$match", d}})
}
func InArray[T any](field string, a []T) (res bson.A) {
	if len(a) == 0 {
		return
	}
	return append(res, bson.D{{"$match", bson.D{{field, bson.D{{"$in", a}}}}}})
}
func Sort(m interface{}) (res bson.A) {
	return append(res, bson.D{{"$sort", m}})
}
func ReplaceRoot(r interface{}) (res bson.A) {
	return append(res, bson.D{{"$replaceRoot", bson.D{{"newRoot", r}}}})
}
func Limit(limit int) (res bson.A) {
	return append(res, bson.D{{"$limit", limit}})
}
func Skip(skip int) (res bson.A) {
	return append(res, bson.D{{"$skip", skip}})
}

func SearchName(search string) bson.A {
	stage := bson.A{
		bson.D{{"$addFields",
			bson.D{
				{"search", mdb.Concat(mdb.Trim("$name"), " ", "$surname")},
			},
		}},
	}
	if search != "" {
		stage = append(stage, bson.D{{"$match",
			bson.D{
				{"search",
					bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}},
		)
	}
	return stage
}
func Search(field, search string) (stage bson.A) {
	if search != "" {
		stage = append(stage, bson.D{{"$match",
			bson.D{
				{field,
					bson.D{
						{"$regex", search},
						{"$options", "i"},
					},
				},
			},
		}},
		)
	}
	return stage
}
func Group(id interface{}, fields ...bson.E) (res bson.A) {
	g := bson.D{{"_id", id}}
	g = append(g, fields...)
	p := bson.D{{"$group", g}}
	return bson.A{p}
}
func Unwind(field string, preserveNullAndEmptyArrays bool) (res bson.A) {

	return bson.A{bson.D{{"$unwind",
		bson.D{
			{"path", "$" + field},
			{"preserveNullAndEmptyArrays", preserveNullAndEmptyArrays},
		},
	}}}
}

func Lookup(from CollName, localfield, foreignFiled, as string, preserveNullAndEmptyArrays interface{}, pipeline ...bson.A) (res bson.A) {
	lookup := bson.D{
		{"from", from.CollName()},
		{"localField", localfield},
		{"foreignField", foreignFiled},
		{"as", as},
	}
	var p bson.A
	for _, a := range pipeline {
		p = append(p, a...)
	}
	if len(p) > 0 {
		lookup = append(lookup, bson.E{"pipeline", p})
	}
	res = bson.A{
		bson.D{
			{"$lookup", lookup},
		},
	}
	switch p := preserveNullAndEmptyArrays.(type) {
	case bool:
		res = append(res, bson.D{
			{"$unwind",
				bson.D{
					{"path", "$" + as},
					{"preserveNullAndEmptyArrays", p},
				},
			},
		},
		)
	}
	return res
}
func LookupLet(from CollName, let interface{}, as string, preserveNullAndEmptyArrays interface{}, pipeline ...bson.A) (res bson.A) {
	lookup := bson.D{
		{"from", from.CollName()},
		{"let", let},
		{"as", as},
	}
	var p bson.A
	for _, a := range pipeline {
		p = append(p, a...)
	}
	if len(p) > 0 {
		lookup = append(lookup, bson.E{"pipeline", p})
	}
	res = bson.A{
		bson.D{
			{"$lookup", lookup},
		},
	}
	switch p := preserveNullAndEmptyArrays.(type) {
	case bool:
		res = append(res, bson.D{
			{"$unwind",
				bson.D{
					{"path", "$" + as},
					{"preserveNullAndEmptyArrays", p},
				},
			},
		},
		)
	}
	return res
}
func Merge(pipeline ...bson.A) (res bson.A) {
	for _, a := range pipeline {
		res = append(res, a...)
	}
	return
}
func Facet(fields ...bson.E) bson.A {
	return bson.A{bson.D{{"$facet", fields}}}
}
