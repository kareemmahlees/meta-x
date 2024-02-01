package lib

type UpdateTableProps struct {
	Operation struct {
		Type string      `json:"type" validate:"required,alpha,oneof=add modify delete" enums:"add,modify,delete"`
		Data interface{} `json:"data" validate:"required,notEmpty,updateTableData"`
	} `json:"operation" validate:"required"`
}
