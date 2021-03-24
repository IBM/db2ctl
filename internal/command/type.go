package command

import (
	"time"

	"github.com/IBM/db2ctl/internal/bash"
	"github.com/IBM/db2ctl/internal/config"
	"github.com/spf13/pflag"
)

//Instance is the main struct for command configs
type Instance struct {
	CombinedConfig *config.Combined
	Error          error
	Flags          *pflag.FlagSet
	StartTime      time.Time
	ConfigDir      string
	bash.Instance
	S3ConfigInstance *config.S3ConfigStruct
}

//SetAutoYesEnabled is a setter for autoYesEnabled field
func (i *Instance) SetAutoYesEnabled(autoYesEnabled bool) *Instance {
	i.AutoYesEnabled = autoYesEnabled
	return i
}

//SetRunOnLocal is a setter for runOnLocal field
func (i *Instance) SetRunOnLocal(runOnLocal bool) *Instance {
	i.RunOnLocal = runOnLocal
	return i
}
