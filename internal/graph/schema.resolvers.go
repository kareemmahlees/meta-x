package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.38

import (
	"context"
	"meta-x/internal/db"
	"meta-x/internal/graph/model"
	"meta-x/lib"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

// CreateDatabase is the resolver for the createDatabase field.
func (r *mutationResolver) CreateDatabase(ctx context.Context, name string) (*model.CreateDatabaseResponse, error) {
	// var rowAffected int64
	// var err error

	// switch r.Provider {
	// case lib.SQLITE3:
	// 	rowAffected,err = db.AttachSqliteDatabase(r.DB,)
	// }

	rowsAffected, err := db.CreatePgMysqlDatabase(r.DB, name)
	if err != nil {
		return nil, err
	}
	return &model.CreateDatabaseResponse{
		Created: int(rowsAffected),
	}, nil
}

// CreateTable is the resolver for the createTable field.
func (r *mutationResolver) CreateTable(ctx context.Context, name string, props []*model.CreateTableData) (*model.CreateTableResponse, error) {
	data := make(map[string]lib.CreateTableProps)
	for _, col := range props {
		data[*col.ColName] = lib.CreateTableProps{
			Type:     *col.Props.Type,
			Nullable: col.Props.Nullable,
			Default:  col.Props.Default,
			Unique:   col.Props.Unique,
		}
	}
	err := db.CreateTable(r.DB, name, data)
	if err != nil {
		return nil, err
	}
	return &model.CreateTableResponse{
		Created: name,
	}, nil
}

// DeleteTable is the resolver for the deleteTable field.
func (r *mutationResolver) DeleteTable(ctx context.Context, name string) (*model.SuccessResponse, error) {
	err := db.DeleteTable(r.DB, name)
	if err != nil {
		return nil, err
	}
	return &model.SuccessResponse{
		Success: true,
	}, nil
}

// UpdateTable is the resolver for the updateTable field.
func (r *mutationResolver) UpdateTable(ctx context.Context, name string, prop *model.UpdateTableData) (*model.SuccessResponse, error) {
	data := lib.UpdateTableProps{}
	data.Operation.Type = string(prop.Operation.Type)
	switch data.Operation.Type {
	case "add":
		if prop.Operation.ColumnsToAdd == nil {
			return nil, gqlerror.Errorf("Operation type 'add' must specifiy 'ColumnsToAdd' field")
		}
		data.Operation.Data = prop.Operation.ColumnsToAdd
	case "modify":
		if prop.Operation.ColumnsToModify == nil {
			return nil, gqlerror.Errorf("Operation type 'modify' must specifiy 'ColumnsToModify' field")
		}
		data.Operation.Data = prop.Operation.ColumnsToModify
	case "delete":
		if prop.Operation.ColumnsToDelete == nil {
			return nil, gqlerror.Errorf("Operation type 'delete' must specifiy 'ColumnsToDelete' field")
		}
		columnstoDetel := []interface{}{}
		for _, col := range prop.Operation.ColumnsToDelete {
			columnstoDetel = append(columnstoDetel, *col)
		}
		data.Operation.Data = columnstoDetel
	}
	err := db.UpdateTable(r.DB, name, data)
	if err != nil {
		return nil, err
	}
	return &model.SuccessResponse{
		Success: true,
	}, nil
}

// Databases is the resolver for the databases field.
func (r *queryResolver) Databases(ctx context.Context) ([]*string, error) {
	var dbs []*string
	var err error

	provider := r.Provider
	switch provider {
	case lib.SQLITE3:
		dbs, err = db.ListDatabasesSqlite(r.DB)
	case lib.PSQL:
		dbs, err = db.ListDatabasesPgMySQL(r.DB, lib.PSQL)
	case lib.MYSQL:
		dbs, err = db.ListDatabasesPgMySQL(r.DB, lib.MYSQL)
	}
	if err != nil {
		return nil, err
	}

	return dbs, nil
}

// Tables is the resolver for the tables field.
func (r *queryResolver) Tables(ctx context.Context) ([]*string, error) {
	tables, err := db.ListTables(r.DB)
	if err != nil {
		return nil, err
	}
	var ps []*string
	for _, v := range tables {
		table := v
		ps = append(ps, &table)
	}

	return ps, nil
}

// Table is the resolver for the table field.
func (r *queryResolver) Table(ctx context.Context, name *string) ([]*model.TableInfo, error) {
	result, err := db.GetTableInfo(r.DB, *name)
	if err != nil {
		return nil, err
	}
	var tableInfo []*model.TableInfo
	for _, info := range result {
		field := info.Field
		typ := info.Type
		null := info.Null
		key := info.Key
		mod := &model.TableInfo{
			Field:   &field,
			Type:    &typ,
			Null:    &null,
			Key:     &key,
			Default: info.Default,
			Extra:   info.Extra,
		}
		tableInfo = append(tableInfo, mod)
	}
	return tableInfo, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
