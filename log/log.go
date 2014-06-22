package log

import (
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wsxiaoys/terminal/color"
)

/*
An OutWriter is responsible for for implementing the io.Writer interface.
*/
type OutWriter struct {
	*logrus.Logger
	fmtString string
}

/*
NewOutWriter accepts a logger and a format string and returns an OutWriter.
When written to, the OutWriter will take the input, split it into lines, and
print it to the logger using the provided format string.  The intended use case
of this functionality is for printing nice, colorful messages
*/
func NewOutWriter(logger *logrus.Logger, fmtString string) *OutWriter {
	return &OutWriter{
		Logger:    logger,
		fmtString: fmtString,
	}
}

/*
Write writes the provided bytes, one line at a time, after interpolating them
into the provided format string, to the provided logger.
*/
func (ow *OutWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for _, line := range lines {
		if logrus.IsTerminal() {
			ow.Debug(color.Sprintf(ow.fmtString, line))
		} else {
			ow.Debug(ow.fmtString, line)
		}
	}

	return len(p), nil
}
