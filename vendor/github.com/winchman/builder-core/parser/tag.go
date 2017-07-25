package parser

import (
	"bytes"
	"github.com/rafecolton/go-gitutils"
	"text/template"
	"time"
)

// Tag is for tagging
type Tag struct {
	value string
}

/*
NewTag returns a Tag instance.  See function implementation for details on what
args to pass.
*/
func NewTag(value string) Tag {
	return Tag{value: value}
}

// Evaluate evaluates any git-based tags
func (t Tag) Evaluate(top string) string {
	switch t.value {
	case "git:branch":
		return git.Branch(top)
	case "git:rev", "git:sha":
		return git.Sha(top)
	case "git:short", "git:tag":
		return git.Tag(top)
	}
	funcMap := template.FuncMap{
		"branch": func() string { return git.Branch(top) },
		"sha":    func() string { return git.Sha(top) },
		"tag":    func() string { return git.Tag(top) },
		"date":   func(format string) string { return time.Now().Format(format) },
	}
	templ, err := template.New("tagParser").Funcs(funcMap).Parse(t.value)
	if err != nil {
		return ""
	}

	var b bytes.Buffer
	err = templ.Execute(&b, t.value)
	if err != nil {
		return ""
	}

	return b.String()
}
