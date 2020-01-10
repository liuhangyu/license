package public

import (
	"bytes"
)

var (
	buildName     = "switch-license"
	buildVersion  = "development"
	buildBranch   = "development"
	buildCommitID = "development"
	buildTime     = "development"
)

func GetAppInfo() string {
	b := bytes.NewBufferString("\nAppInfo:")
	b.WriteString("\n	Name: 	   " + buildName)
	b.WriteString("\n	Version:   " + buildVersion)
	b.WriteString("\n	Branch:    " + buildBranch)
	b.WriteString("\n	CommitID:  " + buildCommitID)
	b.WriteString("\n	BuildTime: " + buildTime)
	b.WriteString("\n")
	return b.String()
}
