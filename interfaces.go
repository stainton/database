package database

import "github.com/spf13/pflag"

type Database interface {
	Insert(interface{}) (interface{}, error)
	Query(interface{}) (interface{}, error)
	AddFlags(*pflag.FlagSet)
	Validate()
}
