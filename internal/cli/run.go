package cli

import (
	"github.com/spf13/cobra"
	"github.com/swellaby/captain-githook/captaingithook"
)

var runGitHook = captaingithook.RunHook
var hookName string

var runCmd = &cobra.Command{
	Use:  "run",
	RunE: runHook,
}

func init() {
	runCmd.Flags().StringVarP(&hookName, "hook", "n", "", "The git hook to run")
	rootCmd.AddCommand(runCmd)
}

func runHook(*cobra.Command, []string) error {
	output, err := runGitHook(hookName)
	if len(output) > 0 {
		logf("Running hook: '%s'...\n", hookName)
		log(output)
	}
	return err
}
