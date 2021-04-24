package migrate

import (
	"github.com/guodongq/jigsaw/pkg/module"
	"github.com/guodongq/jigsaw/pkg/module/mongodb"
)

type Migrate struct {
	module.DefaultProvider
	mongodb  *mongodb.MongoDB
}

func (m *Migrate) Init() error {
	return nil
}
