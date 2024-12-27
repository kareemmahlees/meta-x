package models

// List of Tables.
//
//	@description	List of tables.
type ListTablesResp struct {
	Tables []*string `json:"tables" example:"table1,table2,table3"`
}

// Table Info.
//
//	@description	Info about a table column.
type TableColumnInfo struct {
	Key      any    `db:"key" json:"key" example:"PRI"`            // Constraint name on the column.
	Default  any    `db:"default" json:"default" example:"123"`    // Default value of the column.
	Name     string `db:"name" json:"name" example:"email"`        // Name of the column.
	Type     string `db:"type" json:"type" example:"varchar"`      // Data type of the column.
	Nullable string `db:"nullable" json:"nullable" example:"true"` // If the column accepts null values or not.
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
