package internal

import "fmt"

func (c *RootCmd) verbosef(format string, args ...interface{}) {
	if !c.Verbose {
		return
	}
	c.verbose(fmt.Sprintf(format, args...))
}

func (c *RootCmd) verbose(msg string) {
	if !c.Verbose {
		return
	}

	fmt.Fprintln(c.Writer, "[verbose] "+msg)
}

func (c *RootCmd) writef(format string, args ...interface{}) {
	c.write(fmt.Sprintf(format, args...))
}

func (c *RootCmd) write(msg string) {
	fmt.Fprintln(c.Writer, msg)
}
