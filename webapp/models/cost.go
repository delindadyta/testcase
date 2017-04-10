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
