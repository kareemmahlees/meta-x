package models

type ListTablesResp struct {
	Tables []string
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
