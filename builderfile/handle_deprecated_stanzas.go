package builderfile

import (
	"errors"
)

// HandleDeprecatedStanzas returns an error if file contains any deprecated
// stanzas that would cause a failure.  If no failure would be caused, the
// error will be nil, but a warning message will be printed out and file will
// be sanitized to effectively not include the deprecated stanza
func (file *Builderfile) HandleDeprecatedStanzas() error {
	var includedWarningPrinted bool

	for _, container := range file.ContainerArr {
		if len(container.Excluded) != 0 {
			logger.Error(excludedStanzaWarningMessage)
			return errors.New(`container section contains deprecated "excluded" section`)
		}

		if len(container.Included) != 0 {
			container.Included = nil
			if !includedWarningPrinted {
				logger.Warn(includedStanzaWarningMessage)
				includedWarningPrinted = true
			}
		}
	}

	if file.ContainerGlobals != nil {
		if len(file.ContainerGlobals.Excluded) != 0 {
			logger.Error(excludedStanzaWarningMessage)
			return errors.New(`container section contains deprecated "excluded" section`)
		}

		if len(file.ContainerGlobals.Included) != 0 {
			file.ContainerGlobals.Included = nil
			if !includedWarningPrinted {
				logger.Warn(includedStanzaWarningMessage)
				includedWarningPrinted = true
			}
		}
	}
	return nil
}
