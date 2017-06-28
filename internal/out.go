package internal

import "fmt"

func verbosef(cfg *Config, format string, args ...interface{}) {
	if !cfg.Verbose {
		return
	}
	verbose(cfg, fmt.Sprintf(format, args...))
}

func verbose(cfg *Config, msg string) {
	if !cfg.Verbose {
		return
	}

	fmt.Fprintln(cfg.Writer, "[verbose] "+msg)
}

func writef(cfg *Config, format string, args ...interface{}) {
	write(cfg, fmt.Sprintf(format, args...))
}

func write(cfg *Config, msg string) {
	fmt.Fprintln(cfg.Writer, msg)
}
