package utils

type Config interface {
	DSN() string
}

type SQLiteConfig struct {
	filePath string
}

func NewSQLiteConfig(filePath string) *SQLiteConfig {
	return &SQLiteConfig{filePath}
}

func (sc *SQLiteConfig) DSN() string {
	return sc.filePath
}
