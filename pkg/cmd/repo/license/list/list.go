package list

import (
	"fmt"
	"net/http"

	"github.com/cli/cli/v2/api"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/internal/tableprinter"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type ListOptions struct {
	IO         *iostreams.IOStreams
	HTTPClient func() (*http.Client, error)
	Config     func() (gh.Config, error)
}

func NewCmdList(f *cmdutil.Factory, runF func(*ListOptions) error) *cobra.Command {
	opts := &ListOptions{
		IO:         f.IOStreams,
		HTTPClient: f.HttpClient,
		Config:     f.Config,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List available repository license templates",
		Aliases: []string{"ls"},
		Args:    cmdutil.ExactArgs(0, "gh repo license list takes no arguments"),
		RunE: func(cmd *cobra.Command, args []string) error {

			if runF != nil {
				return runF(opts)
			}
			return listRun(opts)
		},
	}
	return cmd
}

func listRun(opts *ListOptions) error {
	client, err := opts.HTTPClient()
	if err != nil {
		return err
	}

	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	if err := opts.IO.StartPager(); err == nil {
		defer opts.IO.StopPager()
	} else {
		fmt.Fprintf(opts.IO.ErrOut, "failed to start pager: %v\n", err)
	}

	hostname, _ := cfg.Authentication().DefaultHost()
	licenseTemplates, err := api.ListLicenseTemplates(client, hostname)
	if err != nil {
		return err
	}

	if len(licenseTemplates) == 0 {
		return cmdutil.NewNoResultsError("No repository license templates found")
	}

	return renderLicenseTemplatesTable(licenseTemplates, opts)
}

func renderLicenseTemplatesTable(licenseTemplates []api.License, opts *ListOptions) error {
	t := tableprinter.New(opts.IO, tableprinter.WithHeader("LICENSE KEY", "LICENSE NAME"))
	for _, l := range licenseTemplates {
		t.AddField(l.Key)
		t.AddField(l.Name)
		t.EndRow()
	}

	return t.Render()
}
