package cmd

import (
	"fmt"
	"os"

	"github.com/ffloyd/evergrid-go/simulator/gendata"
	"github.com/spf13/cobra"
)

// simulationCmd represents the simulation command
var gendataCmd = &cobra.Command{
	Use:       "gendata NAME DESTDIR",
	Aliases:   []string{"gen", "g"},
	ValidArgs: []string{"NAME", "DESTDIR"},
	Short:     "Generates scenarion and data for simulation in given directory",
	Long: `Generates scenarion and data for simulation in given directory

DISTDIR argument is a name of directory for generated dimdata. NAME - id a name for main YAML of generated data.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("NAME or DESTDIR argument missing")
			os.Exit(1)
		}

		config := gendata.Config{
			Name:    args[0],
			DestDir: args[1],

			DatsetsCount:   50,
			MinDatasetSize: 1,
			MaxDatasetSize: 20,

			CalculatorsCount: 20,
			MinSpeed:         2000,
			MaxSpeed:         20000,

			CalculatorRuns: 100,
			RunProbability: 0.07,

			NetworkSegments:   3,
			MinNodesInSegment: 5,
			MaxNodesInSegment: 15,
			MinNodeSpeed:      5000,
			MaxNodeSpeed:      20000,
			MinPricePerTick:   1.0 / 60.0,
			MaxPricePerTick:   10.0 / 60.0,
			MinDiskSize:       10,
			MaxDiskSize:       2000,
		}
		gendata.GenData(config)
	},
}

func init() {
	RootCmd.AddCommand(gendataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// simulationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// simulationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
