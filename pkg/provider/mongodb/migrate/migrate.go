package migrate

import (
	"github.com/guodongq/jigsaw/pkg/provider"
	"github.com/guodongq/jigsaw/pkg/provider/mongodb"
	"github.com/guodongq/jigsaw/pkg/provider/settings"
)

type Migrate struct {
	provider.DefaultProvider
	Settings *settings.Settings
	mongodb  *mongodb.MongoDB
}

func (m *Migrate) Init() error {
	return nil
}
