package builderfile

import (
	"github.com/wsxiaoys/terminal/color"
)

const sleepTime = 30

var versionZeroWarningMessage = color.Sprintf(`@{y!}WARNING, PLEASE READ:@{|y}

Uh oh, it looks like you are using a Bobfile with @{!y}a deprecated format@{|y}.

You may want to @{!y}update@{|y} before proceeding.

You should consider @{!y}pressing ^C@{|y} and @{!y}updating your Bobfile@{|y} before proceeding.

For more information, see @{!y}https://github.com/modcloth/docker-builder@{|y}

Sleeping for %d seconds before building...
@{|}`, sleepTime)

var includedStanzaWarningMessage = color.Sprintf(`@{y!}WARNING:@{|y}

Uh oh, it looks like one of the container sections in your Bobfile contains the "included" stanza.

This stanza has been deprecated and will be removed.  The good news is this probably won't affect your build.

Carry on!
@{|}`)

var excludedStanzaWarningMessage = color.Sprintf(`@{r!}ERROR:@{|r}

Uh oh, it looks like one of the container sections in your Bobfile contains the @{r!}"excluded" stanza@{|r}.

This stanza @{r!}has been deprecated@{|r}, as the functionality can be replicated with the .dockerignore file.

This will affect your build, so @{r!}your build is being terminated@{|r}.

Please remove the "excluded" section from your Bobfile and consider introducing a .dockerignore file.
@{|}`)
