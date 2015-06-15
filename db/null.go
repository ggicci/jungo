package db

import (
	"database/sql"
	"encoding/json"
)

type NullBool struct {
	sql.NullBool
}

func (nb NullBool) MarshalJSON() ([]byte, error) { return json.Marshal(nb.Bool) }

type NullFloat64 struct {
	sql.NullFloat64
}

func (nf NullFloat64) MarshalJSON() ([]byte, error) { return json.Marshal(nf.Float64) }

type NullInt64 struct {
	sql.NullInt64
}

func (ni NullInt64) MarshalJSON() ([]byte, error) { return json.Marshal(ni.Int64) }

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) { return json.Marshal(ns.String) }
