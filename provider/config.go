package provider

// SQLiteConfig is the configuration needed for running an SQLite database
type SQLiteConfig struct {
	Filepath                string `validate:"file=readable:createifmissing"`
	CreateTablesIfNotExists bool   `validate:"-"`
	DropTablesIfExists      bool   `validate:"-"`
}
