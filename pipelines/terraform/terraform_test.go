package terraform

import (
	"testing"
)

func TestWithPath(t *testing.T) {
	testPath := "/test/path"
	opt := WithPath(testPath)

	cfg := &TerraformConfig{}
	opt(cfg)

	if cfg.Path != testPath {
		t.Errorf("Expected Path in config to be %s, but got %s", testPath, cfg.Path)
	}
}

func TestWithRecursive(t *testing.T) {
	opt := WithRecursive(true)

	cfg := &TerraformConfig{}
	opt(cfg)

	if cfg.Recursive != true {
		t.Errorf("Expected Recursive in config to be true, but got false")
	}
}
