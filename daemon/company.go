package main

import (
	"unsafe"
	
    "github.com/ktroian/polling/daemon/action"
)

type Company struct {
	Name        string `json:"name" pg:",pk,unique"`
	TaxNumber   string `json:"tax number" pg:",unique"`
	PhoneNumber string `json:"phone" pg:",unique"`
	Address 	string `json:"address"`
	CEO 		string `json:"CEO"`
}

func newAction(actionType int, c *Company) {
    action.Add(actionType, unsafe.Pointer(c))
}

func (c *Company) load() {
	newAction(SelectOne, c)
}

func (c *Company) save() {
    newAction(InsertOne, c)
}

func (c *Company) delete() {
    newAction(DeleteOne, c)
}
