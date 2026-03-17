package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
	"github.com/gl0bal01/dorkhound/internal/interactive"
	"github.com/gl0bal01/dorkhound/internal/output"
	"github.com/gl0bal01/dorkhound/web"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "dorkhound",
	Short:   "OSINT missing person finder",
	Long:    "dorkhound is an OSINT tool that generates search engine dork queries to help locate missing persons.",
	Version: version,
	RunE:    run,
}

func run(cmd *cobra.Command, args []string) error {
	// Read all flags.
	flagName, _ := cmd.Flags().GetString("name")
	flagLocation, _ := cmd.Flags().GetString("location")
	flagAge, _ := cmd.Flags().GetInt("age")
	flagDOB, _ := cmd.Flags().GetString("dob")
	flagAKA, _ := cmd.Flags().GetString("aka")
	flagAssociates, _ := cmd.Flags().GetString("associates")
	flagDescription, _ := cmd.Flags().GetString("description")
	flagCase, _ := cmd.Flags().GetString("case")

	flagOpen, _ := cmd.Flags().GetBool("open")
	flagDashboard, _ := cmd.Flags().GetBool("dashboard")
	flagExport, _ := cmd.Flags().GetString("export")
	flagOutput, _ := cmd.Flags().GetString("output")

	flagCategory, _ := cmd.Flags().GetString("category")
	flagRegion, _ := cmd.Flags().GetString("region")
	flagEngine, _ := cmd.Flags().GetString("engine")
	flagDelay, _ := cmd.Flags().GetInt("delay")
	flagInteractive, _ := cmd.Flags().GetBool("interactive")

	var c *caseinfo.Case
	if flagInteractive {
		result, err := interactive.Run()
		if err != nil {
			return fmt.Errorf("interactive mode: %w", err)
		}
		c = result.Case
		flagEngine = result.Engine
		flagRegion = result.Region
		flagCategory = result.Category
		flagOpen = result.OpenBrowser
	} else {
		if flagCase != "" {
			loaded, err := caseinfo.LoadFromFile(flagCase)
			if err != nil {
				return err
			}
			c = loaded
		}

		// Build CLI overrides case.
		cliCase := &caseinfo.Case{
			Name:        flagName,
			Location:    flagLocation,
			Age:         flagAge,
			DOB:         flagDOB,
			Description: flagDescription,
		}
		if flagAKA != "" {
			cliCase.Aliases = caseinfo.SplitTrim(flagAKA)
		}
		if flagAssociates != "" {
			cliCase.Associates = caseinfo.SplitTrim(flagAssociates)
		}

		// Merge: if case loaded from file, merge CLI overrides; otherwise use CLI case.
		if c != nil {
			c.Merge(cliCase)
		} else {
			c = cliCase
			c.FirstName, c.LastName = caseinfo.ParseName(c.Name)
		}

		// Validate: name is required and must not be whitespace-only.
		if strings.TrimSpace(c.Name) == "" {
			return fmt.Errorf("name is required: use --name or --case")
		}
	}

	allDorks := dork.Generate(c)

	categories := caseinfo.SplitTrim(flagCategory)
	regions := caseinfo.SplitTrim(flagRegion)
	engine := strings.ToLower(flagEngine)

	filtered := dork.Filter(allDorks, categories, regions)
	sorted := dork.Sort(filtered)

	if flagDashboard {
		return output.ServeDashboard(c, sorted, engine, web.DashboardHTML)
	}

	// Output writer: default to os.Stdout; if --output is set, create file.
	var w io.Writer = os.Stdout
	if flagOutput != "" {
		f, err := os.Create(flagOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	// Open URLs in browser if requested.
	if flagOpen {
		output.OpenInBrowser(sorted, engine, time.Duration(flagDelay)*time.Millisecond)
	}

	// Export in requested format.
	exportFormat := strings.ToLower(flagExport)
	switch exportFormat {
	case "discord":
		output.Discord(w, c, sorted, engine)
	case "json":
		if err := output.JSON(w, c, sorted, engine); err != nil {
			return fmt.Errorf("writing JSON: %w", err)
		}
	case "csv":
		if err := output.CSV(w, sorted, engine); err != nil {
			return fmt.Errorf("writing CSV: %w", err)
		}
	case "clipboard":
		if err := output.Clipboard(c, sorted, engine); err != nil {
			return fmt.Errorf("copying to clipboard: %w", err)
		}
		fmt.Fprintln(os.Stderr, "Copied to clipboard.")
	case "":
		// Default: print to stdout (unless --open was the only action).
		if !flagOpen {
			output.Stdout(w, sorted, engine)
		}
	default:
		return fmt.Errorf("unknown export format %q: use discord, json, csv, or clipboard", exportFormat)
	}

	return nil
}

func init() {
	// Input flags
	rootCmd.Flags().StringP("name", "n", "", `Full name as "First Last"`)
	rootCmd.Flags().StringP("location", "l", "", "Last known location")
	rootCmd.Flags().Int("age", 0, "Approximate age")
	rootCmd.Flags().String("dob", "", "Date of birth")
	rootCmd.Flags().String("aka", "", "Aliases/nicknames, comma-separated")
	rootCmd.Flags().String("associates", "", "Known associates, comma-separated")
	rootCmd.Flags().String("description", "", "Physical description")
	rootCmd.Flags().String("case", "", "Path to YAML/JSON case file")

	// Output flags
	rootCmd.Flags().Bool("open", false, "Open all URLs in default browser")
	rootCmd.Flags().Bool("dashboard", false, "Serve local web dashboard")
	rootCmd.Flags().String("export", "", "Export format: discord, json, csv, clipboard")
	rootCmd.Flags().StringP("output", "o", "", "Write export to file instead of stdout")

	// Filter flags
	rootCmd.Flags().String("category", "all", "Filter categories")
	rootCmd.Flags().String("region", "global", "Region filter")
	rootCmd.Flags().String("engine", "google", "Search engine")
	rootCmd.Flags().Int("delay", 100, "Delay in ms between opening browser tabs")

	// Other flags
	rootCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")

	// Shell completions for enum flags.
	// Errors here indicate a programming mistake (wrong flag name), so we panic.
	for _, reg := range []struct {
		flag   string
		values []string
	}{
		{"engine", []string{"google", "bing", "duckduckgo", "yandex"}},
		{"region", []string{"global", "all", "us", "ca", "uk", "au", "ru", "fr", "de", "at", "nl"}},
		{"category", []string{"all", "social", "records", "financial", "location", "forums", "people-db"}},
		{"export", []string{"discord", "json", "csv", "clipboard"}},
	} {
		values := reg.values
		if err := rootCmd.RegisterFlagCompletionFunc(reg.flag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return values, cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
			panic(fmt.Sprintf("registering completion for --%s: %v", reg.flag, err))
		}
	}
}
