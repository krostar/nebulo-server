package orm

import (
	"github.com/jinzhu/gorm"
	"github.com/krostar/nebulo/log"
)

// Logger implements the orm logger interface
type Logger struct {
	gorm.Logger
}

// Print show informations about the SQL queries
func (l *Logger) Print(values ...interface{}) {
	values = values[2:]
	values = append([]interface{}{"SQL"}, values...)
	log.Logln(log.DEBUG, -1, values...)
}
