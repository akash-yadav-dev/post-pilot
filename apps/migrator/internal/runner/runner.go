package runner

import (
	"errors"
	"fmt"
	"log"

	appconfig "post-pilot/apps/migrator/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Command string

const (
	CommandUp      Command = "up"
	CommandDown    Command = "down"
	CommandSteps   Command = "steps"
	CommandGoto    Command = "goto"
	CommandForce   Command = "force"
	CommandVersion Command = "version"
)

type Options struct {
	Command Command
	Steps   int
	Version int
	Yes     bool
}

func Execute(cfg *appconfig.Config, opts Options) error {
	m, err := migrate.New(cfg.SourceURL(), cfg.DatabaseURL())
	if err != nil {
		return fmt.Errorf("initialize migrator: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("migrator source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("migrator db close error: %v", dbErr)
		}
	}()

	switch opts.Command {
	case CommandUp:
		return ignoreNoChange(m.Up())
	case CommandDown:
		if !cfg.AllowDestructive && !opts.Yes {
			return errors.New("down migrations are disabled; set ALLOW_DESTRUCTIVE_MIGRATIONS=true or pass -yes")
		}
		return ignoreNoChange(m.Down())
	case CommandSteps:
		if opts.Steps == 0 {
			return errors.New("steps command requires non-zero -steps")
		}
		if opts.Steps < 0 && !cfg.AllowDestructive && !opts.Yes {
			return errors.New("negative steps are disabled; set ALLOW_DESTRUCTIVE_MIGRATIONS=true or pass -yes")
		}
		return ignoreNoChange(m.Steps(opts.Steps))
	case CommandGoto:
		if opts.Version < 0 {
			return errors.New("goto command requires non-negative -version")
		}
		return ignoreNoChange(m.Migrate(uint(opts.Version)))
	case CommandForce:
		if !cfg.AllowDestructive && !opts.Yes {
			return errors.New("force is disabled; set ALLOW_DESTRUCTIVE_MIGRATIONS=true or pass -yes")
		}
		if opts.Version < 0 {
			return errors.New("force command requires non-negative -version")
		}
		return m.Force(opts.Version)
	case CommandVersion:
		version, dirty, err := m.Version()
		if err != nil {
			if errors.Is(err, migrate.ErrNilVersion) {
				log.Printf("version: none, dirty: false")
				return nil
			}
			return err
		}
		log.Printf("version: %d, dirty: %t", version, dirty)
		return nil
	default:
		return fmt.Errorf("unsupported command: %s", opts.Command)
	}
}

func ignoreNoChange(err error) error {
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}
