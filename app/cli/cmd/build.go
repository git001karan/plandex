package cmd

import (
	"fmt"
	"os"
	"plandex/auth"
	"plandex/lib"
	"plandex/plan_exec"
	"plandex/term"

	"github.com/plandex/plandex/shared"
	"github.com/spf13/cobra"
)

var buildBg bool

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build pending changes",
	// Long:  ``,
	Args: cobra.NoArgs,
	Run:  build,
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVar(&buildBg, "bg", false, "Execute autonomously in the background")
}

func build(cmd *cobra.Command, args []string) {
	if os.Getenv("OPENAI_API_KEY") == "" {
		term.OutputNoApiKeyMsg()
		os.Exit(1)
	}

	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		fmt.Fprintln(os.Stderr, "No current plan")
		return
	}

	didBuild, err := plan_exec.Build(plan_exec.ExecParams{
		CurrentPlanId: lib.CurrentPlanId,
		CurrentBranch: lib.CurrentBranch,
		CheckOutdatedContext: func(cancelOpt bool, maybeContexts []*shared.Context) (bool, bool, bool) {
			return lib.MustCheckOutdatedContext(cancelOpt, false, maybeContexts)
		},
	}, buildBg)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if !didBuild {
		fmt.Println()
		term.PrintCmds("", "log", "tell", "continue")
		return
	}

	if buildBg {
		fmt.Println("🏗️ Building plan in the background")
		fmt.Println()
		term.PrintCmds("", "ps", "connect", "stop")
	} else {
		fmt.Println()
		term.PrintCmds("", "changes", "apply", "log")
	}
}
