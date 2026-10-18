package main

import (
	"bufio"
	"context"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cyberactivelabz/kickbucket/internal/engine"
	"github.com/cyberactivelabz/kickbucket/internal/engine/providers"
	scanserver "github.com/cyberactivelabz/kickbucket/internal/server"
)

//go:embed wordlists/default.txt
var embeddedDefaultWordlist string

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
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		if err := runServe(os.Args[2:]); err != nil {
			fmt.Fprintln(os.Stderr, "kickbucket:", err)
			os.Exit(2)
		}
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
	fmt.Fprintln(output, "Usage:")
	fmt.Fprintln(output, "  kickbucket scan <keyword> [flags]")
	fmt.Fprintln(output, "  kickbucket serve [flags]")
	fmt.Fprintln(output, "\nFlags:")
	fmt.Fprintln(output, "  --wordlist <path>       wordlist file (default: wordlists/default.txt)")
	fmt.Fprintln(output, "  --providers <list>      comma-separated providers: s3,gcs,azure (default: all)")
	fmt.Fprintln(output, "  --object-path <path>    object path to probe for public read")
	fmt.Fprintln(output, "  --json                  emit JSON results")
	fmt.Fprintln(output, "  --addr <address>        serve API address (default: 127.0.0.1:8080)")
}

type serveOptions struct {
	address      string
	wordlistPath string
	objectPath   string
}

func runServe(args []string) error {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		printUsage(os.Stdout)
		return nil
	}
	options, err := parseServeArgs(args)
	if err != nil {
		return err
	}
	if err := validateServeAddress(options.address); err != nil {
		return err
	}
	words, err := readWordlist(options.wordlistPath)
	if err != nil {
		return err
	}
	handler := scanserver.NewHandler(scanserver.Config{Providers: providers.NewRegistry(providers.Config{ObjectPath: options.objectPath}), Words: words})
	server := &http.Server{Addr: options.address, Handler: handler}
	fmt.Fprintf(os.Stderr, "kickbucket API listening on http://%s\n", options.address)
	shutdown, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-shutdown.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
	}()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func validateServeAddress(address string) error {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("invalid serve address %q: %w", address, err)
	}
	if host == "" || strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip == nil || !ip.IsLoopback() {
		return fmt.Errorf("serve address must use localhost or a loopback IP")
	}
	return nil
}

func parseServeArgs(args []string) (serveOptions, error) {
	options := serveOptions{address: "127.0.0.1:8080"}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--addr":
			if i+1 >= len(args) {
				return serveOptions{}, fmt.Errorf("--addr requires a value")
			}
			i++
			options.address = args[i]
		case "--wordlist":
			if i+1 >= len(args) {
				return serveOptions{}, fmt.Errorf("--wordlist requires a value")
			}
			i++
			options.wordlistPath = args[i]
		case "--object-path":
			if i+1 >= len(args) {
				return serveOptions{}, fmt.Errorf("--object-path requires a value")
			}
			i++
			options.objectPath = args[i]
		case "--help", "-h":
			printUsage(os.Stdout)
			return options, nil
		default:
			return serveOptions{}, fmt.Errorf("unknown flag %q", args[i])
		}
	}
	return options, nil
}

func runScan(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		printUsage(os.Stdout)
		return nil
	}
	options, err := parseScanArgs(args)
	if err != nil {
		return err
	}
	words, err := readWordlist(options.wordlistPath)
	if err != nil {
		return err
	}
	wanted := make(map[string]bool)
	for _, name := range options.providers {
		wanted[name] = true
	}
	fmt.Fprintln(os.Stderr, "Only scan cloud storage you own or are explicitly authorized to test.")
	all := providers.NewRegistry(providers.Config{ObjectPath: options.objectPath})
	results := make([]engine.Result, 0)
	for _, provider := range all {
		if !wanted[provider.Name()] {
			continue
		}
		limiter := engine.NewLimiter(10)
		providerResults, scanErr := engine.Scan(context.Background(), provider, provider.Candidates(options.keyword, words), 8, limiter)
		limiter.Stop()
		if scanErr != nil {
			return scanErr
		}
		results = append(results, providerResults...)
	}
	if options.json {
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

type scanOptions struct {
	keyword      string
	wordlistPath string
	providers    []string
	objectPath   string
	json         bool
}

func parseScanArgs(args []string) (scanOptions, error) {
	flagArgs := make([]string, 0, len(args))
	positionals := make([]string, 0, 1)
	for i := 0; i < len(args); i++ {
		argument := args[i]
		switch argument {
		case "--wordlist", "--providers", "--object-path":
			if i+1 >= len(args) {
				return scanOptions{}, fmt.Errorf("%s requires a value", argument)
			}
			flagArgs = append(flagArgs, argument, args[i+1])
			i++
		case "--json":
			flagArgs = append(flagArgs, argument)
		default:
			if strings.HasPrefix(argument, "-") {
				flagArgs = append(flagArgs, argument)
			} else {
				positionals = append(positionals, argument)
			}
		}
	}
	if len(positionals) != 1 {
		return scanOptions{}, fmt.Errorf("exactly one keyword is required")
	}
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	wordlistPath := flags.String("wordlist", "", "wordlist file")
	providerList := flags.String("providers", "s3,gcs,azure", "comma-separated providers")
	objectPath := flags.String("object-path", "", "object path")
	jsonOutput := flags.Bool("json", false, "emit JSON")
	if err := flags.Parse(append(flagArgs, positionals[0])); err != nil {
		return scanOptions{}, err
	}
	options := scanOptions{keyword: positionals[0], wordlistPath: *wordlistPath, objectPath: *objectPath, json: *jsonOutput}
	for _, name := range strings.Split(*providerList, ",") {
		name = strings.ToLower(strings.TrimSpace(name))
		if name != "s3" && name != "gcs" && name != "azure" {
			return scanOptions{}, fmt.Errorf("unknown provider %q", name)
		}
		options.providers = append(options.providers, name)
	}
	return options, nil
}

func readWordlist(path string) ([]string, error) {
	if path == "" {
		return parseWordlist(strings.NewReader(embeddedDefaultWordlist))
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("read wordlist: %w", err)
	}
	defer file.Close()
	return parseWordlist(file)
}

func parseWordlist(input io.Reader) ([]string, error) {
	var words []string
	scanner := bufio.NewScanner(input)
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
