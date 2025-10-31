package models

type ListTablesResp struct {
	Tables []*string `json:"tables" example:"table1,table2,table3"`
}

type TableColumnInfo struct {
	Key      any    `db:"key" json:"key" doc:"Constraint name on the column"`
	Default  any    `db:"default" json:"default" example:"123" doc:"Default value of the column"`
	Name     string `db:"name" json:"name" example:"email" doc:"Name of the column"`
	Type     string `db:"type" json:"type" example:"varchar" doc:"Data type of the column"`
	Nullable string `db:"nullable" json:"nullable" example:"true" doc:"If the column accepts null values or not"`
}

type CreateTablePayload struct {
	Nullable any    `json:"nullable" validate:"omitempty,boolean" doc:"If the column accepts null values or not"`
	Default  any    `json:"default" validate:"omitempty,alphanum" doc:"Default value of the column"`
	Unique   any    `json:"unique" validate:"omitempty,boolean" doc:"Wether to add a unique constraint on the column"`
	ColName  string `json:"column_name" validate:"required,alphanum" doc:"Name of the column"`
	Type     string `json:"type" validate:"required,ascii" doc:"Data type of the column"`
}

type CreateTableResp struct {
	Created string `json:"created" doc:"Created table name"`
}

type AddModifyColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`
	Type    string `json:"type" validate:"required,ascii" example:"Int" doc:"New type"`
}

type DeleteColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"` // Column Name
}
