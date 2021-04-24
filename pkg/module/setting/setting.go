package setting

import (
	"github.com/guodongq/jigsaw/pkg/util/helper"
	"github.com/spf13/cobra"
)

type Setting struct {
	CfgFile *string
}

func New(cmd *cobra.Command) *Setting {
	var st = &Setting{}
	if flag := cmd.Flag("config"); flag != nil {
		if val := flag.Value.String(); len(val) > 0 {
			if ok, _ := helper.PathExists(val); ok {
				st.CfgFile = &val
			}
		}
	}
	return st
}

func (s *Setting) Enable() bool {
	return s.CfgFile != nil
}
