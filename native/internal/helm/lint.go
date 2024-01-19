package helm

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/lint/support"
)

type LintOptions struct {
	Path   string
	Strict bool
	Quiet  bool
}

func Lint(options *LintOptions) ([]string, bool) {
	lintClient := action.NewLint()
	lintClient.Strict = options.Strict
	lintClient.Quiet = options.Quiet
	result := lintClient.Run([]string{options.Path}, make(map[string]interface{}))
	var messages []string
	// From upstream https://github.com/helm/helm/blob/33ab3519849a90549f734fbbbc0aecb7f37f7570/cmd/helm/lint.go#L108-L111
	// All Errors that are generated by a chart with failed lint will be included in the result.Messages too
	// Only consider them if Messages is empty
	if len(result.Messages) == 0 {
		for _, err := range result.Errors {
			messages = append(messages, fmt.Sprintf("Error %s", err.Error()))
		}
	}
	for _, msg := range result.Messages {
		// In strict mode, msg with low severity will be ignored
		if !options.Quiet || msg.Severity > support.InfoSev {
			messages = append(messages, fmt.Sprintf("%s", msg))
		}
	}
	return messages, len(result.Errors) > 0
}