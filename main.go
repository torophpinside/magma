package main

import (
	"magma/cmd"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "pdbmonitor-agent",
	}

	cmd.NewDorksListRunCmd(rootCmd)

	cobra.CheckErr(rootCmd.Execute())
}
