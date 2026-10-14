package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers"
)

// Set at build time by GoReleaser (see .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("kickbucket %s (%s, %s)\n", version, commit, date)
		return
	}
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		printUsage(os.Stdout)
		return
	}
	if len(os.Args) < 2 || os.Args[1] != "scan" {
		printUsage(os.Stderr)
		os.Exit(2)
	}
	if err := runScan(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, "kickbucket:", err)
		os.Exit(2)
	}
}

func printUsage(output *os.File) {
	fmt.Fprintln(output, "Usage: kickbucket scan <keyword> [flags]")
	fmt.Fprintln(output, "\nFlags:")
	fmt.Fprintln(output, "  --wordlist <path>       wordlist file (default: wordlists/default.txt)")
	fmt.Fprintln(output, "  --providers <list>      comma-separated providers: s3,gcs,azure (default: all)")
	fmt.Fprintln(output, "  --json                  emit JSON results")
}

func runScan(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage(os.Stdout)
		return nil
	}
	if strings.HasPrefix(args[0], "-") {
		return fmt.Errorf("keyword is required")
	}
	keyword := args[0]
	wordlistPath := "wordlists/default.txt"
	providerList := "s3,gcs,azure"
	jsonOutput := false
	for i := 1; i < len(args); i++ {
		switch args[i] {
		case "--json":
			jsonOutput = true
		case "--wordlist", "--providers":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires a value", args[i])
			}
			i++
			if args[i-1] == "--wordlist" {
				wordlistPath = args[i]
			} else {
				providerList = args[i]
			}
		default:
			return fmt.Errorf("unknown flag %q", args[i])
		}
	}

	words, err := readWordlist(wordlistPath)
	if err != nil {
		return err
	}
	wanted := make(map[string]bool)
	for _, name := range strings.Split(providerList, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name != "s3" && name != "gcs" && name != "azure" {
			return fmt.Errorf("unknown provider %q", name)
		}
		wanted[name] = true
	}

	fmt.Fprintln(os.Stderr, "Only scan cloud storage you own or are explicitly authorized to test.")
	all := providers.NewRegistry(providers.Config{})
	results := make([]engine.Result, 0)
	for _, provider := range all {
		if !wanted[provider.Name()] {
			continue
		}
		limiter := engine.NewLimiter(10)
		providerResults, scanErr := engine.Scan(context.Background(), provider, provider.Candidates(keyword, words), 8, limiter)
		limiter.Stop()
		if scanErr != nil {
			return scanErr
		}
		results = append(results, providerResults...)
	}
	if jsonOutput {
		return json.NewEncoder(os.Stdout).Encode(results)
	}
	for _, result := range results {
		status := strings.ToUpper(result.Access.String())
		if result.Error != "" {
			status = "ERROR: " + result.Error
		}
		fmt.Printf("%-5s %-30s %s\n", result.Provider, result.Name, status)
	}
	fmt.Printf("%d buckets found, %d exposed\n", len(results), countExposed(results))
	return nil
}

func readWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}
	defer file.Close()
	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" && !strings.HasPrefix(word, "#") {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}
	return words, nil
}

func countExposed(results []engine.Result) int {
	count := 0
	for _, result := range results {
		if result.Access == engine.AccessListable || result.Access == engine.AccessPublicRead {
			count++
		}
	}
	return count
}
