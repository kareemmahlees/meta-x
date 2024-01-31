package models

type ListDatabasesResult struct {
	Databases []string
}

type CreateDatabaseResult struct {
	Created int
}

type AttachSqliteDBPayload struct {
	Name string `json:"name" validate:"required,alpha" binding:"required"`
	File string `json:"file" validate:"required,filepath"`
}

type CreatePgMySqlDBPayload struct {
	Name string `json:"name" validate:"required,alphanum"`
}
