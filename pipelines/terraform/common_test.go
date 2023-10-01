package terraform

import (
	"context"
	"testing"

	"dagger.io/dagger"
)

func TestWithCommonPath(t *testing.T) {
	tf := &Terraform{
		commonConfig: &CommonConfig{},
	}
	opt := WithCommonPath("/test/path")
	opt(tf)

	if tf.commonConfig.Path != "/test/path" {
		t.Errorf("Expected path to be '/test/path', got '%s'", tf.commonConfig.Path)
	}
}

func TestWithFmtRecursive(t *testing.T) {
	tf := &Terraform{
		fmtConfig: &FmtConfig{},
	}
	opt := WithFmtRecursive(true)
	opt(tf)

	if !tf.fmtConfig.Recursive {
		t.Errorf("Expected Recursive to be true, got false")
	}
}

func TestNew(t *testing.T) {
	ctx := context.Background()
	client := &dagger.Client{}
	tf := New(ctx, client)

	if tf.ctx != ctx {
		t.Errorf("Expected context to be equal, got different context")
	}
	if tf.client != client {
		t.Errorf("Expected client to be equal, got different client")
	}
}

// Continue writing tests for other functions...
