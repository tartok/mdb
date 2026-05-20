package mdb

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Cond(If, Then, Else interface{}) bson.D {
	return bson.D{{"$cond", bson.D{{"if", If}, {"then", Then}, {"else", Else}}}}
}
func If(cond bool, f func() interface{}) interface{} {
	if cond {
		return f()
	}
	return nil
}
func Archived(a bool) bson.D {
	if a {
		return bson.D{{"state.archived", a}}
	}
	return Ne("state.archived", true)
}
func Invisible(a bool) bson.D {
	if a {
		return bson.D{{"state.invisible", a}}
	}
	return Ne("state.invisible", true)
}
func Disabled(a bool) bson.D {
	if a {
		return bson.D{{"state.disabled", a}}
	}
	return Ne("state.disabled", true)
}
func IfNull(res ...interface{}) bson.D {
	return bson.D{{"$ifNull", res}}
}
func Or(m ...interface{}) bson.D {
	//var or bson.A
	//or = append(or, m...)
	return bson.D{{"$or", m}}
}
func NeAggregation(m ...interface{}) bson.D {
	var ne bson.A
	ne = append(ne, m...)
	return bson.D{{"$ne", ne}}
}
func Size(field string) bson.D {
	return bson.D{{"$size", field}}
}
func Ne(field string, m interface{}) bson.D {
	return bson.D{{field, bson.D{{"$ne", m}}}}
}
func Not(m interface{}) bson.D {
	return bson.D{{"$not", m}}
}
func And(m ...interface{}) bson.D {
	//var and bson.A
	//and = append(and, m...)
	return bson.D{{"$and", m}}
}
func Expr(e bson.D) bson.D {
	return bson.D{{"$expr", e}}
}
func MinAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$min", eq}}
}
func Min(m interface{}) bson.D {
	return bson.D{{"$min", m}}
}
func MaxAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$max", eq}}
}
func Max(m interface{}) bson.D {
	return bson.D{{"$max", m}}
}
func EqAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$eq", eq}}
}
func Eq(field string, m interface{}) bson.D {
	return bson.D{{field, bson.D{{"$eq", m}}}}
}
func Exists(field string, v bool) bson.D {
	return bson.D{{field, bson.D{{"$exists", v}}}}
}
func ExistsAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$exists", eq}}
}
func Slice(m ...interface{}) bson.D {
	var slice bson.A
	slice = append(slice, m...)
	return bson.D{{"$slice", slice}}
}
func Between(field string, from, to interface{}) bson.D {
	return bson.D{{field, bson.D{{"$gte", from}, {"$lte", to}}}}
}
func Intersect(fFrom, fTo string, from, to interface{}) bson.D {
	return And(
		Or(
			bson.D{{fFrom, nil}},
			bson.D{{fFrom, bson.D{{"$lte", to}}}},
		),
		Or(
			bson.D{{fTo, nil}},
			bson.D{{fTo, bson.D{{"$gte", from}}}},
		),
	)
}
func Gte(field string, value interface{}) bson.D {
	return bson.D{{field, bson.D{{"$gte", value}}}}
}
func Gt(field string, value interface{}) bson.D {
	return bson.D{{field, bson.D{{"$gt", value}}}}
}
func GteAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$gte", eq}}
}
func GtAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$gt", eq}}
}
func Lte(field string, value interface{}) bson.D {
	return bson.D{{field, bson.D{{"$lte", value}}}}
}
func Lt(field string, value interface{}) bson.D {
	return bson.D{{field, bson.D{{"$lt", value}}}}
}
func LteAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$lte", eq}}
}
func LtAggregation(m ...interface{}) bson.D {
	var eq bson.A
	eq = append(eq, m...)
	return bson.D{{"$lt", eq}}
}
func InArray[T any](field string, a []T) (res bson.D) {
	if len(a) == 0 {
		return
	}
	return bson.D{{field, bson.D{{"$in", a}}}}
}
func AllArray[T any](field string, a []T) (res bson.D) {
	if len(a) == 0 {
		return
	}
	return bson.D{{field, bson.D{{"$all", a}}}}
}
func Tags(field string, a []string) (res interface{}) {
	if len(a) == 0 {
		return
	}
	var or []interface{}
	b := false
	for _, s := range a {
		ss := strings.Split(strings.Trim(s, "#"), "#")
		if len(ss) == 1 {
			or = append(or, bson.M{field: ss[0]})
		} else {
			or = append(or, AllArray(field, ss))
			b = true
		}
	}
	if b {
		if len(a) == 1 {
			return or[0]
		}
		return Or(or...)
	}
	return InArray(field, a)
}
func Substr(field string, start, length int) (res bson.D) {
	return bson.D{{"$substr", bson.A{field, start, length}}}
}
func NinArray[T any](field string, a []T) (res bson.D) {
	if len(a) == 0 {
		return
	}
	return bson.D{{field, bson.D{{"$nin", a}}}}
}
func Trim(value string) (res bson.D) {
	return bson.D{{"$trim", bson.D{{"input", value}}}}
}
func Concat(values ...interface{}) (res bson.D) {
	return bson.D{{"$concat", values}}
}

func Function(body string, args ...interface{}) bson.D {
	return bson.D{{"$function", bson.D{
		{"body", body},
		{"args", args},
		{"lang", "js"},
	}}}
}
func Filter(input, as string, cond interface{}) bson.D {
	return bson.D{
		{"$filter", bson.D{
			{"input", input},
			{"as", as},
			{"cond", cond},
		}}}
}
func In(a ...interface{}) bson.D {
	return bson.D{
		{"$in", a},
	}
}
func Multiply(a ...interface{}) bson.D {
	return bson.D{
		{"$multiply", a},
	}
}
func Divide(a ...interface{}) bson.D {
	return bson.D{
		{"$divide", a},
	}
}
func Add(a ...interface{}) bson.D {
	return bson.D{
		{"$add", a},
	}
}
func Subtract(a ...interface{}) bson.D {
	return bson.D{
		{"$subtract", a},
	}
}
func Sum(a interface{}) bson.D {
	return bson.D{
		{"$sum", a},
	}
}
func Round(a ...interface{}) bson.D {
	return bson.D{
		{"$round", a},
	}
}
func Convert(input interface{}, to string) bson.D {
	return bson.D{
		{"$convert", bson.D{{"input", input}, {"to", to}}},
	}
}
func Map(input, as string, in interface{}) bson.D {
	return bson.D{
		{"$map", bson.D{{"input", input}, {"as", as}, {"in", in}}},
	}
}
func Mod(a ...interface{}) bson.D {
	return bson.D{
		{"$mod", a},
	}
}
func DayOfWeek(field interface{}) bson.D {
	return bson.D{
		{"$dayOfWeek", field},
	}
}
func DateToString(date, format, timezone string, onNull interface{}) bson.D {

	r := bson.D{{"$dateToString", bson.D{
		{"date", date},
		{"format", format},
	}}}
	if timezone != "" {
		r = append(r, bson.E{"timezone", timezone})
	}
	if onNull != nil {
		r = append(r, bson.E{"onNull", onNull})
	}
	return r
}
func DateFromParts(year, month, day interface{}) bson.D {
	return bson.D{
		{"$dateFromParts", bson.D{
			{"year", year},
			{"month", month},
			{"day", day},
		}},
	}
}
func DateTrunc(date, unit string) bson.D {

	r := bson.D{{"$dateTrunc", bson.D{
		{"date", date},
		{"unit", unit},
	}}}
	return r
}

func DateSubtract(startDate interface{}, unit string, amount interface{}) bson.D {
	return bson.D{
		{"$dateSubtract",
			bson.D{
				{"startDate", startDate},
				{"unit", unit},
				{"amount", amount},
			},
		},
	}
}
func DateAdd(startDate interface{}, unit string, amount interface{}) bson.D {
	return bson.D{
		{"$dateAdd",
			bson.D{
				{"startDate", startDate},
				{"unit", unit},
				{"amount", amount},
			},
		},
	}
}
func FirstDayOfWeek(dateField string) bson.D {
	return Subtract(dateField,
		Multiply(
			Subtract(
				bson.D{
					{"$dayOfWeek", DateSubtract(dateField, "day", 1)},
				}, 1), 24*60*60*1000))
}
func WeeksCount(dateField string) bson.D {
	return bson.D{
		{"$dateDiff", bson.D{
			{"startDate", time.Date(1970, 1, 1, 0, 0, 0, 0, nil)},
			{"endDate", dateField},
			{"unit", "week"},
			{"startOfWeek", "mon"},
		}},
	}
}

func ObjectToArray(o string) bson.D {
	return bson.D{{"$objectToArray", o}}
}
func FacetField(name string, pipeline ...bson.A) bson.E {
	var res bson.A
	for _, a := range pipeline {
		res = append(res, a...)
	}
	return bson.E{name, res}
}
func Regex(pattern, options string) bson.D {
	return bson.D{{"$regex", pattern}, {"$options", options}}
}
func MergeObjects(objects ...interface{}) bson.D {
	return bson.D{{"$mergeObjects", objects}}
}
func Reduce(input, initialValue, in interface{}) bson.D {
	return bson.D{
		{"$reduce", bson.D{
			{"input", input},
			{"initialValue", initialValue},
			{"in", in},
		}}}
}
func SetUnion(v ...interface{}) bson.D {
	return bson.D{
		{"$setUnion", v},
	}
}
func ToObjectId(d interface{}) bson.D {
	return bson.D{{"$toObjectId", d}}
}
