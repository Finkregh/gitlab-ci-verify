package checks

import (
	"fmt"
	"github.com/timo-reymann/gitlab-ci-verify/internal/logging"
	"github.com/timo-reymann/gitlab-ci-verify/internal/shellcheck"
	"github.com/timo-reymann/gitlab-ci-verify/pkg/gitlab/ci-yaml"
	"strconv"
	"sync"
)

type ShellScriptCheck struct {
}

func (s ShellScriptCheck) shellcheckLevelToSeverity(level string) int {
	switch level {
	case "error":
		return SeverityError
	case "warning":
		return SeverityWarning
	case "info":
		return SeverityInfo
	case "style":
		return SeverityStyle
	}
	return -1
}

func (s ShellScriptCheck) Run(i *CheckInput) ([]CheckFinding, error) {
	findingsChan := make(chan CheckFinding)

	shellChecker, err := shellcheck.NewShellChecker()
	if err != nil {
		return nil, err
	}
	defer shellChecker.Close()

	var wg sync.WaitGroup
	for jobWithScripts := range ci_yaml.ExtractScripts(i.CiYaml.ParsedYamlDoc) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key, parts := range jobWithScripts.ScriptParts {
				lines, joinedScript := ci_yaml.Concat(parts)
				result, err := shellChecker.AnalyzeSnippet(joinedScript, i.Configuration.ShellcheckFlags)
				if err != nil {
					logging.Warn("Failed to analyze snippet in job", jobWithScripts.JobName)
					continue
				}

				for _, f := range result.Findings {
					findingsChan <- CheckFinding{
						Severity: s.shellcheckLevelToSeverity(f.Level),
						Code:     fmt.Sprintf("SC-%s", strconv.Itoa(f.Code)),
						Line:     lines[f.Line-1].Node.Line,
						Message:  fmt.Sprintf("[%s:%s:%d] %s", jobWithScripts.JobName, key, f.Line, f.Message),
						Link:     fmt.Sprintf("https://www.shellcheck.net/wiki/SC%d", f.Code),
						File:     i.Configuration.GitLabCiFile,
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(findingsChan)
	}()

	findings := make([]CheckFinding, 0)
	for f := range findingsChan {
		findings = append(findings, f)
	}

	return findings, nil
}