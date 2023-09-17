package lib

type CreateTableProps struct {
	Type     string      `json:"type" validate:"required,alpha,oneof=text number date"`
	Nullable interface{} `json:"nullable" validate:"omitempty,boolean"`
	Default  interface{} `json:"default" validate:"omitempty,alpha" `
	Unique   interface{} `json:"unique" validate:"omitempty,boolean"`
}

type UpdateTableProps struct {
	Operation struct {
		Type string      `json:"type" validate:"required,alpha,oneof=add modify delete" enums:"add,modify,delete"`
		Data interface{} `json:"data" validate:"required,notEmpty,updateTableData"`
	} `json:"operation" validate:"required"`
}
