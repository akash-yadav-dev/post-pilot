package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	appconfig "post-pilot/apps/migrator/internal/config"
	"post-pilot/apps/migrator/internal/runner"
)

func main() {
	log.SetFlags(log.LstdFlags | log.LUTC | log.Lshortfile)

	cmd := flag.String("command", "up", "migration command: up|down|steps|goto|force|version")
	steps := flag.Int("steps", 0, "number of migration steps for command=steps (positive for up, negative for down)")
	version := flag.Int("version", 0, "target version for command=goto or command=force")
	yes := flag.Bool("yes", false, "confirm destructive operations for down/negative-steps/force")
	flag.Parse()

	cfg, err := appconfig.Load()
	if err != nil {
		fatalf("load config", err)
	}

	opts := runner.Options{
		Command: runner.Command(strings.ToLower(strings.TrimSpace(*cmd))),
		Steps:   *steps,
		Version: *version,
		Yes:     *yes,
	}

	if err := runner.Execute(cfg, opts); err != nil {
		fatalf("run migration", err)
	}

	log.Printf("migration command %q completed successfully", opts.Command)
}

func fatalf(action string, err error) {
	log.Printf("%s failed: %v", action, err)
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
