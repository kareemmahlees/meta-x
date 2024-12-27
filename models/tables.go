package models

//	@description	List of tables.
type ListTablesResp struct {
	Tables []*string `json:"tables" example:"table1,table2,table3"`
}

type TableInfoResp struct {
	Key      any    `db:"key" json:"key"`
	Default  any    `db:"default" json:"default"`
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	Nullable string `db:"nullable" json:"nullable"`
}

type CreateTablePayload struct {
	Nullable interface{} `json:"nullable" validate:"omitempty,boolean"`
	Default  interface{} `json:"default" validate:"omitempty,alphanum" `
	Unique   interface{} `json:"unique" validate:"omitempty,boolean"`
	ColName  string      `json:"column_name" validate:"required,alphanum"`
	Type     string      `json:"type" validate:"required,ascii"`
}

type CreateTableResp struct {
	Created string `json:"created"`
}

type AddModifyColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`
	Type    string `json:"type" validate:"required,ascii"`
}

type DeleteColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`
}
