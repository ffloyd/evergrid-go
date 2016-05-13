package cmd

import (
	"fmt"
	"os"

	"github.com/ffloyd/evergrid-go/simulation"
	"github.com/spf13/cobra"
)

// simulationCmd represents the simulation command
var simulationCmd = &cobra.Command{
	Use:       "simulation SIMDATA",
	Aliases:   []string{"sim", "s"},
	ValidArgs: []string{"INFRASTRUCTURE"},
	Short:     "Starts simulation with given options",
	Long: `Starts simulation with given infracstructure.

SIMDATA argument is a name of YAML file with simdata config.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("INFRASTRUCTURE argument missing")
			os.Exit(1)
		}
		sim := simulation.New(args[0])
		sim.Run()
	},
}

func init() {
	RootCmd.AddCommand(simulationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// simulationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// simulationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
