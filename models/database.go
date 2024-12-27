package models

// @description A list of all available databases.
type ListDatabasesResp struct {
	Databases []*string `json:"databases" example:"test,prod,main"`
}

type AttachSqliteDBPayload struct {
	Name string `json:"name" validate:"required,alpha" binding:"required"`
	File string `json:"file" validate:"required,filepath"`
}

type CreatePgMySqlDBPayload struct {
	Name string `json:"name" validate:"required,alphanum" example:"Users"` // Database name.
}
