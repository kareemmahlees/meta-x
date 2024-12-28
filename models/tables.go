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

// Create Table Payload.
//
//	@description	Data about the column to create.
type CreateTablePayload struct {
	Nullable interface{} `json:"nullable" validate:"omitempty,boolean"`    // If the column accepts null values or not.
	Default  interface{} `json:"default" validate:"omitempty,alphanum"`    // Default value of the column.
	Unique   interface{} `json:"unique" validate:"omitempty,boolean"`      // Wether to add a unique constraint on the column.
	ColName  string      `json:"column_name" validate:"required,alphanum"` // Name of the column.
	Type     string      `json:"type" validate:"required,ascii"`           // Data type of the column.
}

// Table Created Successfully.
//
// @description Table Created Successfully.
type CreateTableResp struct {
	Created string `json:"created"` // Created table name.
}

// Add/Modify a column payload.
//
// @description Adding or Modifying a column payload.
type AddModifyColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"`     // Column Name
	Type    string `json:"type" validate:"required,ascii" example:"Int"` // New type
}

// Column deletion payload.
//
// @description Column to delete data.
type DeleteColumnPayload struct {
	ColName string `json:"column_name" validate:"required,alphanum"` // Column Name
}
