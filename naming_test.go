package gorm_test

import (
	"testing"

	"github.com/jinzhu/gorm"
)

func TestTheNamingStrategy(t *testing.T) {

	cases := []struct {
		name     string
		namer    gorm.Namer
		expected string
	}{
		{name: "auth", expected: "auth", namer: gorm.ToDBName},
		{name: "userRestrictions", expected: "user_restrictions", namer: gorm.ToTableName},
		{name: "clientID", expected: "client_id", namer: gorm.ToColumnName},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.namer(c.name)
			if result != c.expected {
				t.Errorf("error in naming strategy. expected: %v got :%v\n", c.expected, result)
			}
		})
	}

}

func TestNamingStrategy(t *testing.T) {

	dbNameNS := func(name string) string {
		return "db_" + name
	}
	tableNameNS := func(name string) string {
		return "tbl_" + name
	}
	columnNameNS := func(name string) string {
		return "col_" + name
	}

	ns := &gorm.NamingStrategy{
		DB:     dbNameNS,
		Table:  tableNameNS,
		Column: columnNameNS,
	}

	gorm.AddNamingStrategy(ns)

	cases := []struct {
		name     string
		namer    gorm.Namer
		expected string
	}{
		{name: "auth", expected: "db_auth", namer: gorm.ToDBName},
		{name: "user", expected: "tbl_user", namer: gorm.ToTableName},
		{name: "password", expected: "col_password", namer: gorm.ToColumnName},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := c.namer(c.name)
			if result != c.expected {
				t.Errorf("error in naming strategy. expected: %v got :%v\n", c.expected, result)
			}
		})
	}

}
