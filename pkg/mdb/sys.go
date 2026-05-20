package mdb

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	GetterId interface {
		GetId() bson.ObjectID
	}
	GetterCompanyId interface {
		GetCompanyId() bson.ObjectID
	}
	GetterUserId interface {
		GetUserId() bson.ObjectID
	}
	GetterPatientId interface {
		GetPatientId() bson.ObjectID
	}
	ICTime interface {
		InitCTime()
		ClearCTime()
	}
	IUTime interface {
		InitUTime()
		ClearUTime()
	}
	ICWho interface {
		InitCWho(who bson.ObjectID)
		ClearCWho()
	}
	IUWho interface {
		InitUWho(who bson.ObjectID)
		ClearUWho()
	}
)

type (
	Id struct {
		Id *bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty" swaggerignore:"true"`
	}
	CompanyId struct {
		CompanyId *bson.ObjectID `json:"companyId,omitempty" bson:"companyId,omitempty" swaggerignore:"true"`
	}
	UserId struct {
		UserId *bson.ObjectID `json:"userId,omitempty" bson:"userId,omitempty" binding:"required"`
	}
	UserIdP struct {
		UserId *bson.ObjectID `json:"userId" bson:"userId"`
	}
	PatientId struct {
		PatientId *bson.ObjectID `json:"patientId,omitempty" bson:"patientId,omitempty" binding:"required"`
	}
	PatientIdP struct {
		PatientId *bson.ObjectID `json:"patientId" bson:"patientId"`
	}
	CTime struct {
		CTime *time.Time `json:"cTime" bson:"cTime,omitempty" swaggerignore:"true"`
	}
	UTime struct {
		UTime *time.Time `json:"uTime" bson:"uTime,omitempty" swaggerignore:"true"`
	}
	CUTime struct {
		CTime `bson:",inline"`
		UTime `bson:",inline"`
	}
	CWho struct {
		CWho *bson.ObjectID `json:"cWho" bson:"cWho,omitempty"  swaggerignore:"true"`
	}
	UWho struct {
		UWho *bson.ObjectID `json:"uWho" bson:"uWho,omitempty" swaggerignore:"true"`
	}
	CUWho struct {
		CWho `bson:",inline"`
		UWho `bson:",inline"`
	}
	CU struct {
		CUTime `bson:",inline"`
		CUWho  `bson:",inline"`
	}
	Base struct {
		Id       `bson:",inline"`
		CU       `json:"CU" bson:"CU" swaggerignore:"true"`
		Settings map[string]interface{} `json:"-" bson:"_settings,omitempty"`
		Old      []byte                 `json:"-" bson:"-"`
	}
)

func getSettings(patch []string, s interface{}, v any) error {
	if len(patch) > 0 {
		if m, ok := s.(map[string]interface{}); ok {
			return getSettings(patch[1:], m[patch[0]], v)
		}
		return fmt.Errorf("patch not found")
	}
	d, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, v)
}
func (b Base) GetSettings(patch string, v any) error {
	var p []string
	if patch != "" {
		p = strings.Split(patch, ".")
	}
	return getSettings(p, b.Settings, v)
}

func (i *Id) GetId() bson.ObjectID {
	if i != nil && i.Id != nil {
		return *i.Id
	}
	return bson.ObjectID{}
}
func (i *Id) SetId(id bson.ObjectID) {
	i.Id = &id
}
func (i *CompanyId) GetCompanyId() bson.ObjectID {
	if i != nil && i.CompanyId != nil {
		return *i.CompanyId
	}
	return bson.ObjectID{}
}
func (i *CompanyId) SetCompanyId(companyId bson.ObjectID) {
	i.CompanyId = &companyId
}
func (i *UserId) GetUserId() bson.ObjectID {
	if i != nil {
		return *i.UserId
	}
	return bson.ObjectID{}
}
func (i *UserIdP) SetUserId(userId bson.ObjectID) {
	i.UserId = &userId
}
func (i *UserIdP) GetUserId() bson.ObjectID {
	if i != nil {
		return *i.UserId
	}
	return bson.ObjectID{}
}
func (i *UserId) SetUserId(userId bson.ObjectID) {
	i.UserId = &userId
}
func (i *PatientId) GetPatientId() bson.ObjectID {
	if i != nil {
		return *i.PatientId
	}
	return bson.ObjectID{}
}
func (i *PatientIdP) SetPatientId(patientId bson.ObjectID) {
	i.PatientId = &patientId
}
func (i *PatientIdP) GetPatientId() bson.ObjectID {
	if i != nil {
		return *i.PatientId
	}
	return bson.ObjectID{}
}
func (i *PatientId) SetPatientId(patientId bson.ObjectID) {
	i.PatientId = &patientId
}

func (c *CTime) InitCTime() {
	c.CTime = Ref(time.Now())
}
func (c *CTime) ClearCTime() {
	c.CTime = nil
}
func (c *UTime) InitUTime() {
	c.UTime = Ref(time.Now())
}
func (c *UTime) ClearUTime() {
	c.UTime = nil
}
func (u *UWho) InitUWho(who bson.ObjectID) {
	u.UWho = &who
}
func (u *UWho) ClearUWho() {
	u.UWho = nil
}

func (c *CWho) InitCWho(who bson.ObjectID) {
	c.CWho = &who
}
func (c *CWho) ClearCWho() {
	c.CWho = nil
}
func (c CU) LastUnix() int64 {
	if c.UTime.UTime != nil {
		return c.UTime.UTime.Unix()
	}
	if c.CTime.CTime != nil {
		return c.CTime.CTime.Unix()
	}
	return 0
}

func Ref[T any](t T) *T { return &t }
