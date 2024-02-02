package models

type ListTablesResp struct {
	Tables []*string `json:"tables"`
}

type TableInfoResp struct {
	Name     string `db:"name" json:"name"`
	Type     string `db:"type" json:"type"`
	Nullable string `db:"nullable" json:"nullable"`
	Key      any    `db:"key" json:"key"`
	Default  any    `db:"default" json:"default"`
}

type CreateTablePayload struct {
	ColName  string      `json:"column_name" validate:"required,alphanum"`
	Type     string      `json:"type" validate:"required,alphanum"`
	Nullable interface{} `json:"nullable" validate:"omitempty,boolean"`
	Default  interface{} `json:"default" validate:"omitempty,alphanum" `
	Unique   interface{} `json:"unique" validate:"omitempty,boolean"`
}

type CreateTableResp struct {
	Created string
}

type AddUpdateColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`
	Type    string `json:"type" validate:"required,ascii"`
}

type DeleteColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`
}

type SuccessResp struct {
	Success bool `json:"success"`
}
