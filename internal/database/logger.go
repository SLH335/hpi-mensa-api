package database

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
)

// gooseAdapter wraps zerolog.Logger to implement goose.Logger
type gooseAdapter struct {
	logger zerolog.Logger
}

// SetGooseLogger replaces goose's internal logger
func SetGooseLogger(logger zerolog.Logger) {
	goose.SetLogger(&gooseAdapter{logger: logger})
}

// Print logs general messages
func (g *gooseAdapter) Print(v ...any) {
	msg := fmt.Sprint(v...)
	g.logStructured(msg)
}

// Printf logs formatted messages
func (g *gooseAdapter) Printf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	g.logStructured(msg)
}

// Fatal logs fatal messages and exits
func (g *gooseAdapter) Fatal(v ...any) {
	g.logger.Fatal().Msg(fmt.Sprint(v...))
}

// Fatalf logs fatal formatted messages and exits
func (g *gooseAdapter) Fatalf(format string, v ...any) {
	g.logger.Fatal().Msgf(format, v...)
}

// --- Helper to parse goose messages and extract structured fields ---
func (g *gooseAdapter) logStructured(msg string) {
	// 1. No migrations
	if strings.HasPrefix(msg, "goose: no migrations to run. current version:") {
		version := strings.TrimSpace(strings.Split(msg, ":")[2])
		g.logger.Info().
			Str("component", "goose").
			Str("version", version).
			Msg("No migrations to run")
		return
	}

	// 2. Applied migration
	// Example: "applied migration 00001_init.sql (1)"
	re := regexp.MustCompile(`goose: applied migration (\S+) \((\d+)\)`)
	if matches := re.FindStringSubmatch(msg); len(matches) == 3 {
		g.logger.Info().
			Str("component", "goose").
			Str("migration", matches[1]).
			Str("version", matches[2]).
			Msg("Applied migration")
		return
	}

	// 3. Dropped migration (if using goose Down)
	// Example: "rolled back migration 00001_init.sql (1)"
	reDown := regexp.MustCompile(`goose: rolled back migration (\S+) \((\d+)\)`)
	if matches := reDown.FindStringSubmatch(msg); len(matches) == 3 {
		g.logger.Info().
			Str("component", "goose").
			Str("migration", matches[1]).
			Str("version", matches[2]).
			Msg("Rolled back migration")
		return
	}

	// 4. Fallback: log as plain message
	g.logger.Info().Str("component", "goose").Msg(msg)
}
