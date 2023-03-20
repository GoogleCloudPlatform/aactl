package main

import (
	"context"

	"github.com/GoogleCloudPlatform/aactl/pkg/types"
	"github.com/GoogleCloudPlatform/aactl/pkg/vul"
	gha "github.com/sethvargo/go-githubactions"
)

var (
	version = "v0.0.1-default"
	commit  = "none"
	date    = "unknown"
)

func main() {
	a := gha.WithFieldsMap(map[string]string{
		"version": version,
		"commit":  commit,
		"date":    date,
	})

	a.Infof("starting action")
	defer a.Infof("action completed")

	project := a.GetInput("project")
	source := a.GetInput("source")
	file := a.GetInput("file")
	format := a.GetInput("format")

	if project == "" || source == "" || file == "" || format == "" {
		a.Fatalf("required parameter not provided, got project=%s, source=%s, file=%s, format=%s", project, source, file, format)
	}

	f, err := types.ParseSourceFormat(format)
	if err != nil {
		a.Fatalf("error parsing source format, got %s", err)
	}

	opt := &types.ImportOptions{
		Project: project,
		Source:  source,
		File:    file,
		Format:  f,
	}

	if err := vul.Import(context.Background(), opt); err != nil {
		a.Fatalf("error importing vulnerabilities, got %s", err)
	}

	// TODO: add output params with the number of vulnerabilities imported
	a.SetOutput("import_count", "0")
}
