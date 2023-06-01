package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/nestoca/cx/src/cmd"
	"github.com/nestoca/cx/src/cmd/meta"
	"github.com/nestoca/cx/src/cmd/meta/export"
	"github.com/nestoca/cx/src/cmd/meta/project"
	"github.com/nestoca/cx/src/cmd/meta/releases"
	"github.com/nestoca/cx/src/cmd/meta/version"
	"github.com/nestoca/cx/src/cmd/meta/version/current"
	"github.com/nestoca/cx/src/cmd/meta/version/next"
	"github.com/nestoca/cx/src/cmd/promote"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
	}()

	rootCmd := cmd.NewRoot()

	metaCmd := meta.New()
	rootCmd.AddCommand(metaCmd)
	metaCmd.AddCommand(project.New())
	metaCmd.AddCommand(releases.New())
	metaCmd.AddCommand(export.New())

	versionCmd := version.New()
	versionCmd.AddCommand(current.New())
	versionCmd.AddCommand(next.New())
	metaCmd.AddCommand(versionCmd)

	promoteCmd := promote.New()
	rootCmd.AddCommand(promoteCmd)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(-1)
	}
}
