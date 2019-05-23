package models

type Db struct {
	Sql  string   `json:"sql"`
	Args []string `json:"args"`
}
