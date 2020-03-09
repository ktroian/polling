package main

import (
	"fmt"
	"sync"
	"errors"
)

type Company struct {
	Name        string `json:"name"`
	TaxNumber   string `json:"tax number"`
	PhoneNumber string `json:"phone"`
	Address 	string `json:"address"`
	CEO 		string `json:"CEO"`
}

type Companies []Company

type AllCompanies struct {
	companies *Companies
	mtx sync.Mutex
}

func (c *AllCompanies) add(company Company) {
	c.mtx.Lock()
	*c.companies = append(*c.companies, company)
	c.mtx.Unlock()
}

func (c *AllCompanies) delete(identifier string) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for i, ent := range *c.companies {
		if ent.Name == identifier || ent.TaxNumber == identifier {
			*c.companies = append((*c.companies)[:i], (*c.companies)[i+1:]...)
			return nil
		}
	}
	errMsg := fmt.Sprintf("No company with identifier: '%s'", identifier)
	return errors.New(errMsg)
}

func (c *AllCompanies) get() []Company {
	return *c.companies
}

func (c *Company) dump() []string {
	return []string{
		c.Name,
		c.TaxNumber,
		c.PhoneNumber,
		c.Address,
		c.CEO,
	}
}

func load(fields []string) Company {
	var loadedCompany Company

	if len(fields) == 5 {
		loadedCompany.Name 		  = fields[0]
		loadedCompany.TaxNumber   = fields[1]
		loadedCompany.PhoneNumber = fields[2]
		loadedCompany.Address     = fields[3]
		loadedCompany.CEO 		  = fields[4]
	}

	return loadedCompany
}
