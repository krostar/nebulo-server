package provider

// DefaultConfig is the configuration needed for running every database
type DefaultConfig struct {
	CreateTablesIfNotExists bool `validate:"-"`
	DropTablesIfExists      bool `validate:"-"`
}

// SQLiteConfig is the configuration needed for running an SQLite database
type SQLiteConfig struct {
	DefaultConfig
	Filepath string `validate:"file=readable:createifmissing"`
}
