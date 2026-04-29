package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gl0bal01/dorkhound/internal/caseinfo"
	"github.com/gl0bal01/dorkhound/internal/dork"
	"github.com/gl0bal01/dorkhound/internal/integrations/nuclei"
	"github.com/gl0bal01/dorkhound/internal/interactive"
	"github.com/gl0bal01/dorkhound/internal/output"
	"github.com/gl0bal01/dorkhound/internal/preflight"
	"github.com/gl0bal01/dorkhound/web"
	"github.com/spf13/cobra"
)

var version = "dev"

var knownCategories = []string{
	"all", "social", "records", "financial", "location", "forums", "people-db",
	"email", "phone", "username", "cache", "documents", "dating", "marketplace",
	"nuclei", "image", "gravatar", "github", "academic", "direct-profile",
	"twitter", "reddit", "fundraiser", "telegram", "vehicle", "crypto",
}

var knownRegions = []string{"all", "global", "us", "ca", "uk", "au", "ru", "fr", "de", "at", "nl"}

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
	flagEmails, _ := cmd.Flags().GetString("emails")
	flagPhones, _ := cmd.Flags().GetString("phones")
	flagUsernames, _ := cmd.Flags().GetString("usernames")

	flagOpen, _ := cmd.Flags().GetBool("open")
	flagDashboard, _ := cmd.Flags().GetBool("dashboard")
	flagExport, _ := cmd.Flags().GetString("export")
	flagOutput, _ := cmd.Flags().GetString("output")
	flagForceOutput, _ := cmd.Flags().GetBool("force-output")

	flagCategory, _ := cmd.Flags().GetString("category")
	flagRegion, _ := cmd.Flags().GetString("region")
	flagEngine, _ := cmd.Flags().GetString("engine")
	flagDelay, _ := cmd.Flags().GetInt("delay")
	flagBatch, _ := cmd.Flags().GetInt("batch")
	flagBatchPause, _ := cmd.Flags().GetDuration("batch-pause")
	flagInteractive, _ := cmd.Flags().GetBool("interactive")
	flagNoiseFilter, _ := cmd.Flags().GetBool("noise-filter")
	flagNuclei, _ := cmd.Flags().GetBool("nuclei")
	flagNucleiTags, _ := cmd.Flags().GetString("nuclei-tags")
	flagNucleiTimeout, _ := cmd.Flags().GetDuration("nuclei-timeout")
	flagPreflight, _ := cmd.Flags().GetBool("preflight")
	flagPhotoURL, _ := cmd.Flags().GetString("photo-url")
	flagPhotoPath, _ := cmd.Flags().GetString("photo-path")
	flagPreflightTimeout, _ := cmd.Flags().GetDuration("preflight-timeout")
	flagPreflightConcurrency, _ := cmd.Flags().GetInt("preflight-concurrency")
	flagPreflightRate, _ := cmd.Flags().GetDuration("preflight-rate")
	flagStats, _ := cmd.Flags().GetBool("stats")
	flagListCategories, _ := cmd.Flags().GetBool("list-categories")
	flagListRegions, _ := cmd.Flags().GetBool("list-regions")

	if flagListCategories {
		for _, c := range knownCategories {
			fmt.Println(c)
		}
		return nil
	}

	if flagListRegions {
		for _, r := range knownRegions {
			fmt.Println(r)
		}
		return nil
	}

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
		if flagEmails != "" {
			cliCase.Emails = caseinfo.SplitTrim(flagEmails)
		}
		if flagPhones != "" {
			cliCase.Phones = caseinfo.SplitTrim(flagPhones)
		}
		if flagUsernames != "" {
			cliCase.Usernames = caseinfo.SplitTrim(flagUsernames)
		}
		cliCase.PhotoURL = flagPhotoURL
		cliCase.PhotoPath = flagPhotoPath

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
	allDorks = dork.Dedup(allDorks)

	// Run nuclei if requested and usernames are available.
	if flagNuclei && len(c.Usernames) > 0 {
		nucleiDorks, err := nuclei.Run(context.Background(), c.Usernames, flagNucleiTags, flagNucleiTimeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nuclei: %v\n", err)
		}
		allDorks = append(allDorks, nucleiDorks...)
	}

	categories := caseinfo.SplitTrim(flagCategory)
	regions := caseinfo.SplitTrim(flagRegion)
	engine := strings.ToLower(flagEngine)

	filtered := dork.Filter(allDorks, categories, regions)
	sorted := dork.Sort(filtered)

	if flagPreflight {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		survivors, preflightReport := preflight.Run(ctx, sorted, preflight.Options{
			Concurrency: flagPreflightConcurrency,
			Timeout:     flagPreflightTimeout,
			RateLimit:   flagPreflightRate,
		})
		fmt.Fprintf(os.Stderr, "preflight: %d checked, %d alive, %d dead (%d search dorks skipped)\n",
			preflightReport.Checked, preflightReport.Alive, preflightReport.Dead, preflightReport.Skipped)
		sorted = survivors
	}

	if flagNoiseFilter {
		sorted = dork.ApplyNoiseFilter(sorted)
	}

	if flagStats {
		fmt.Print(dork.Stats(sorted))
		return nil
	}

	if flagDashboard {
		return output.ServeDashboard(c, sorted, engine, web.DashboardHTML)
	}

	// Output writer: default to os.Stdout; if --output is set, create file.
	var w io.Writer = os.Stdout
	if flagOutput != "" {
		f, err := createOutputFile(flagOutput, flagForceOutput)
		if err != nil {
			return fmt.Errorf("creating output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	// Open URLs in browser if requested.
	if flagOpen {
		output.OpenInBrowser(sorted, engine, time.Duration(flagDelay)*time.Millisecond, flagBatchPause, flagBatch)
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
	case "tracelabs":
		if err := output.TraceLabs(w, c, sorted, engine); err != nil {
			return fmt.Errorf("writing TraceLabs: %w", err)
		}
	case "":
		// Default: print to stdout (unless --open was the only action).
		if !flagOpen {
			output.Stdout(w, sorted, engine)
		}
	default:
		return fmt.Errorf("unknown export format %q: use discord, json, csv, clipboard, or tracelabs", exportFormat)
	}

	return nil
}

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion script",
	Long: `Generate a shell completion script for dorkhound.

Examples:

  # Bash (current session)
  source <(dorkhound completion bash)

  # Bash (persist)
  dorkhound completion bash > /etc/bash_completion.d/dorkhound

  # Zsh
  dorkhound completion zsh > "${fpath[1]}/_dorkhound"

  # Fish
  dorkhound completion fish > ~/.config/fish/completions/dorkhound.fish

  # PowerShell
  dorkhound completion powershell > dorkhound.ps1`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)

	// Input flags
	rootCmd.Flags().StringP("name", "n", "", `Full name as "First Last"`)
	rootCmd.Flags().StringP("location", "l", "", "Last known location")
	rootCmd.Flags().Int("age", 0, "Approximate age")
	rootCmd.Flags().String("dob", "", "Date of birth (YYYY-MM-DD)")
	rootCmd.Flags().String("aka", "", "Aliases/nicknames, comma-separated")
	rootCmd.Flags().String("associates", "", "Known associates, comma-separated")
	rootCmd.Flags().String("description", "", "Physical description")
	rootCmd.Flags().String("case", "", "Path to YAML/JSON case file")
	rootCmd.Flags().String("emails", "", "Email addresses, comma-separated")
	rootCmd.Flags().String("phones", "", "Phone numbers, comma-separated")
	rootCmd.Flags().String("usernames", "", "Usernames/handles, comma-separated")
	rootCmd.Flags().String("photo-url", "", "URL of the target's photo for reverse image search")
	rootCmd.Flags().String("photo-path", "", "Local path to the target's photo (operator note only)")

	// Output flags
	rootCmd.Flags().Bool("open", false, "Open all URLs in default browser")
	rootCmd.Flags().Bool("dashboard", false, "Serve local web dashboard")
	rootCmd.Flags().String("export", "", "Export format: discord, json, csv, clipboard, tracelabs")
	rootCmd.Flags().StringP("output", "o", "", "Write export to file instead of stdout")
	rootCmd.Flags().Bool("force-output", false, "Overwrite --output if it already exists")

	// Filter flags
	rootCmd.Flags().String("category", "all", "Category filter; use --list-categories for choices")
	rootCmd.Flags().String("region", "global", "Region filter; use --list-regions for choices")
	rootCmd.Flags().String("engine", "google", "Search engine: google, bing, duckduckgo, yandex")
	rootCmd.Flags().Int("delay", 2000, "Delay in ms between opening browser tabs (2000 is safe for Google to avoid CAPTCHA)")
	rootCmd.Flags().Int("batch", 10, "Number of URLs to open per batch (0 = no batching) when --open is set")
	rootCmd.Flags().Duration("batch-pause", 30*time.Second, "Sleep duration between batches when --open is set")

	// Noise filter
	rootCmd.Flags().Bool("noise-filter", false, "Append noise-suppression operators to search queries")

	// Preflight flags
	rootCmd.Flags().Bool("preflight", false, "Probe direct-URL dorks and drop dead links before output")
	rootCmd.Flags().Duration("preflight-timeout", 5*time.Second, "Per-request timeout for preflight probes")
	rootCmd.Flags().Int("preflight-concurrency", 8, "Max parallel preflight HTTP requests")
	rootCmd.Flags().Duration("preflight-rate", 250*time.Millisecond, "Per-worker cooldown between preflight requests")

	// Nuclei flags
	rootCmd.Flags().Bool("nuclei", false, "Run nuclei OSINT templates against usernames (requires nuclei in PATH)")
	rootCmd.Flags().String("nuclei-tags", nuclei.DefaultTags, "Nuclei template tags to run")
	rootCmd.Flags().Duration("nuclei-timeout", 2*time.Minute, "Timeout per username for nuclei")

	// Info / utility flags
	rootCmd.Flags().Bool("stats", false, "Print dork count breakdown by category/region/priority and exit")
	rootCmd.Flags().Bool("list-categories", false, "Print all known category names and exit")
	rootCmd.Flags().Bool("list-regions", false, "Print all known region names and exit")

	// Other flags
	rootCmd.Flags().BoolP("interactive", "i", false, "Interactive mode (guided prompts)")

	// Shell completions for enum flags.
	// Errors here indicate a programming mistake (wrong flag name), so we panic.
	for _, reg := range []struct {
		flag   string
		values []string
	}{
		{"engine", []string{"google", "bing", "duckduckgo", "yandex"}},
		{"region", knownRegions},
		{"category", knownCategories},
		{"export", []string{"discord", "json", "csv", "clipboard", "tracelabs"}},
	} {
		values := reg.values
		if err := rootCmd.RegisterFlagCompletionFunc(reg.flag, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return values, cobra.ShellCompDirectiveNoFileComp
		}); err != nil {
			panic(fmt.Sprintf("registering completion for --%s: %v", reg.flag, err))
		}
	}
}

func createOutputFile(path string, force bool) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE | os.O_EXCL
	if force {
		flags = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	}
	f, err := os.OpenFile(path, flags, 0600) // #nosec G304 -- CLI intentionally writes an operator-supplied output path.
	if err != nil {
		if os.IsExist(err) && !force {
			return nil, fmt.Errorf("%s already exists; pass --force-output to overwrite", path)
		}
		return nil, err
	}
	return f, nil
}
