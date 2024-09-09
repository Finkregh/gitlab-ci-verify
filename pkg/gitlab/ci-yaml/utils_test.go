package ci_yaml

import (
	"errors"
	"github.com/timo-reymann/gitlab-ci-verify/internal/yamlpathutils"
	"testing"
)

func Test_MustPath(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("Expected error to panic")
		}
	}()
	yamlpathutils.MustPath(nil, errors.New("alarm"))
}
