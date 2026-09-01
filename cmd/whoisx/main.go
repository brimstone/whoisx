// Copyright (c) 2026 Matt Robinson brimstone@the.narro.ws

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "whoisx",
	Short: "Look up owner info for IPv4 and IPv6 addresses",
	Long: `whoisx queries the appropriate RIR whois server (ARIN, RIPE, APNIC,
AFRINIC, LACNIC) for an IP address, recursing to the correct registry,
and prints the result as JSON.

IP addresses are read from stdin, one per line.`,
}

func main() {
	var programLevel = new(slog.LevelVar) // Info by default

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

	// check for environment variable now
	if strings.ToLower(os.Getenv("LOG_LEVEL")) == "debug" || slices.Contains(os.Args, "--debug") {
		programLevel.Set(slog.LevelDebug)
	}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().Bool("debug", false, "debug mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(Cmd())
	rootCmd.AddCommand(Update())

	ctx := context.Background()

	cmd, _, _ := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	if (len(os.Args) == 1) || (len(os.Args) > 1 &&
		os.Args[1] != "completion" &&
		os.Args[1] != "__complete") &&
		cmd.Use == rootCmd.Use &&
		!errors.Is(cmd.Flags().Parse(os.Args[1:]), pflag.ErrHelp) {
		slog.Debug("Defaulting to whoisx command")

		args := []string{"whoisx"}
		args = append(args, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	err := rootCmd.ExecuteContext(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
