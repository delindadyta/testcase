package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type CostModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	Category      string
	Cost          float64
}

func NewCostModel() *CostModel {
	m := new(CostModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *CostModel) RecordID() interface{} {
	return e.Id
}

func (m *CostModel) TableName() string {
	return "Cost"
}

type MonthlyModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId ` bson:"_id" , json:"_id" `
	Description   string
	Category      string
	ProjectCost   string
	ActualCost    string
}

func NewMonthlyModel() *MonthlyModel {
	m := new(MonthlyModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *MonthlyModel) RecordID() interface{} {
	return e.Id
}

func (m *MonthlyModel) TableName() string {
	return "MonthlyEpenses"
}
