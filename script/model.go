// Package provides structs and functions to work with SQL statements from
// files.
package script

import "github.com/jmoiron/sqlx"

var (
	format = "20060102150405"
)

// Rows is a wrapper around sql.Rows which caches costly reflect operations
// during a looped StructScan.
type Rows = sqlx.Rows

// Param is a command parameter for given query.
type Param = interface{}
