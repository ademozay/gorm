package gorm

import (
	"bytes"
	"strings"
	"sync"
)

// Namer is a function which is given a string and return a string
type Namer func(string) string

// NamingStrategy represents naming strategy functions for specific properties
type NamingStrategy struct {
	Default Namer
	DB      Namer
	Table   Namer
	Column  Namer
}

// the more question the more answer!

// namingStrategy is the NamingStrategy.
var namingStrategy = &NamingStrategy{
	Default: defaultNamer,
}

// once to avoid modifying naming strategy once it's been set.
var once sync.Once

type property int

const (
	database property = iota
	table
	column
)

var cache = newSafeMap()

// AddNamingStrategy sets NamingStrategy for once.
func AddNamingStrategy(ns *NamingStrategy) {
	once.Do(func() {
		if ns.Default == nil {
			ns.Default = defaultNamer
		}
		namingStrategy = ns
	})
}

// ToDBName renames given name with database namer.
func ToDBName(name string) string {
	return namingStrategy.rename(name, database)
}

// ToTableName renames given name with table namer.
func ToTableName(name string) string {
	return namingStrategy.rename(name, table)
}

// ToColumnName renames given name with column namer.
func ToColumnName(name string) string {
	return namingStrategy.rename(name, column)
}

// rename renames given name with given property namer.
// if no namer found for given property, default namer will be used.
func (ns *NamingStrategy) rename(name string, property property) string {
	if renamed := cache.Get(name); renamed != "" {
		return renamed
	}

	var renamed string

	namer := ns.namer(property)

	renamed = namer(name)

	cache.Set(name, renamed)

	return renamed
}

func (ns *NamingStrategy) namer(property property) Namer {
	switch property {
	case database:
		return ns.DB
	case table:
		return ns.Table
	case column:
		return ns.Column
	default:
		return ns.Default
	}
}

// defaultNamer converts given string to snake_case
func defaultNamer(name string) string {
	const (
		lower = false
		upper = true
	)

	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}
