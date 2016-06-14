package cmd

import (
	"fmt"
	"os"

	"github.com/ffloyd/evergrid-go/simulator/gendata"
	"github.com/spf13/cobra"
)

// flag vars
var datasetsCount int
var minDatasetSize int
var maxDatasetSize int
var calculatorsCount int
var minCalculatorComplexity int
var maxCalculatorComplexity int

var calculatorRuns int
var runProbability float64

var networkSegments int
var minNodesInSegment int
var maxNodesInSegment int
var minNodeSpeed int
var maxNodeSpeed int
var minPricePerTick float64
var maxPricePerTick float64
var minDiskSize int
var maxDiskSize int

var gendataCmd = &cobra.Command{
	Use:       "gendata NAME DESTDIR",
	Aliases:   []string{"gen", "g"},
	ValidArgs: []string{"NAME", "DESTDIR"},
	Short:     "Generates scenarion and data for simulation in given directory",
	Long: `Generates scenarion and data for simulation in given directory

DISTDIR - is a name of directory for generated dimdata.
NAME    - is a name for main YAML of generated data.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("NAME or DESTDIR argument missing")
			os.Exit(1)
		}

		config := gendata.Config{
			Name:    args[0],
			DestDir: args[1],

			DatsetsCount:   datasetsCount,
			MinDatasetSize: minDatasetSize,
			MaxDatasetSize: maxDatasetSize,

			CalculatorsCount:        calculatorsCount,
			MinCalculatorComplexity: minCalculatorComplexity,
			MaxCalculatorComplexity: maxCalculatorComplexity,

			CalculatorRuns: calculatorRuns,
			RunProbability: runProbability,

			NetworkSegments:   networkSegments,
			MinNodesInSegment: minNodesInSegment,
			MaxNodesInSegment: maxNodesInSegment,
			MinNodeSpeed:      minNodeSpeed,
			MaxNodeSpeed:      maxNodeSpeed,
			MinPricePerTick:   minPricePerTick,
			MaxPricePerTick:   maxPricePerTick,
			MinDiskSize:       minDiskSize,
			MaxDiskSize:       maxDiskSize,
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
	gendataCmd.Flags().IntVar(&datasetsCount, "datasets_count", 20, "Datasets count")
	gendataCmd.Flags().IntVar(&minDatasetSize, "min_dataset_size", 1, "Minimal dataset size in gigabytes")
	gendataCmd.Flags().IntVar(&maxDatasetSize, "max_dataset_size", 20, "Maximum dataset size in gigabytes")

	gendataCmd.Flags().IntVar(&calculatorsCount, "calculators_count", 20, "Calculators count")
	gendataCmd.Flags().IntVar(&minCalculatorComplexity, "min_calculator_complexity", 2000, "Minimal calculator complexity in MFlop per megabyte of data")
	gendataCmd.Flags().IntVar(&maxCalculatorComplexity, "max_calculator_complexity", 20000, "Maximal calculator complexity in MFlop per megabyte of data")

	gendataCmd.Flags().IntVar(&calculatorRuns, "calculator_runs", 100, "Count of calculator runs in sumulation")
	gendataCmd.Flags().Float64Var(&runProbability, "run_probability", 0.1, "'Tick includes calculator run' probability")

	gendataCmd.Flags().IntVar(&networkSegments, "network_segments", 10, "Count of network segments")
	gendataCmd.Flags().IntVar(&minNodesInSegment, "min_nodes_in_segment", 5, "Minimal count of nodes inside particular segment")
	gendataCmd.Flags().IntVar(&maxNodesInSegment, "max_nodes_in_segment", 10, "Maximal count of nodes inside particular segment")
	gendataCmd.Flags().IntVar(&minNodeSpeed, "min_node_speed", 5000, "Minimal node performance in MFlops")
	gendataCmd.Flags().IntVar(&maxNodeSpeed, "max_node_speed", 20000, "Maximal node performance in MFlops")
	gendataCmd.Flags().Float64Var(&minPricePerTick, "min_price_per_tick", 0.16, "Minimal price of one minute of node work")
	gendataCmd.Flags().Float64Var(&maxPricePerTick, "max_price_per_tick", 0.16, "Maximal price of one minute of node work")
	gendataCmd.Flags().IntVar(&minDiskSize, "min_disk_size", 10, "Minimal disk size on node in gigabytes")
	gendataCmd.Flags().IntVar(&maxDiskSize, "max_disk_size", 2000, "Maximal disk size on node in gigabytes")
}
