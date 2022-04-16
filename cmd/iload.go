package cmd

import (
	"github.com/sergrom/iload/internal/app"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const NumberOfThreadsDefault = 5

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "iload",
	Short: "iload - simple file downloader",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		var (
			inputFile, outputDirectory string
			threadsNumber              int
			verbose                    bool
			err                        error
		)

		if inputFile, err = cmd.Flags().GetString("input-file"); err != nil {
			log.Fatal(err.Error())
		}
		if outputDirectory, err = cmd.Flags().GetString("output-dir"); err != nil {
			log.Fatal(err.Error())
		}
		if threadsNumber, err = cmd.Flags().GetInt("threads-num"); err != nil {
			log.Fatal(err.Error())
		}
		if verbose, err = cmd.Flags().GetBool("verbose"); err != nil {
			log.Fatal(err.Error())
		}

		err = app.NewLoader().
			WithInputFile(inputFile).
			WithOutputDirectory(outputDirectory).
			Download(threadsNumber, verbose)

		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("input-file", "f", "", "Input file with urls you want to download. Each url must be in separate line")
	rootCmd.Flags().StringP("output-dir", "d", "", "Output directory to which you want to save downloaded files")
	rootCmd.Flags().IntP("threads-num", "t", NumberOfThreadsDefault, "Number of threads")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose")

	_ = rootCmd.MarkFlagRequired("input-file")
	_ = rootCmd.MarkFlagRequired("output-dir")
}
