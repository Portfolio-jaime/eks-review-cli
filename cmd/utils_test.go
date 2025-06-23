package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetEffectiveNamespace_AllNamespacesFlag(t *testing.T) {
	if ns := GetEffectiveNamespace("", true, "", false); ns != "" {
		t.Errorf("expected empty namespace, got %q", ns)
	}
}

func TestGetEffectiveNamespace_NamespaceFlag(t *testing.T) {
	if ns := GetEffectiveNamespace("myns", false, "", false); ns != "myns" {
		t.Errorf("expected 'myns', got %q", ns)
	}
}

func TestGetEffectiveNamespace_NamespaceFlagAllAllowed(t *testing.T) {
	if ns := GetEffectiveNamespace("All", false, "", true); ns != "" {
		t.Errorf("expected empty namespace when flag is 'all', got %q", ns)
	}
}

func TestGetEffectiveNamespace_FromKubeconfig(t *testing.T) {
	tmp := t.TempDir()
	kubeDir := filepath.Join(tmp, ".kube")
	if err := os.MkdirAll(kubeDir, 0o755); err != nil {
		t.Fatalf("failed to create kube dir: %v", err)
	}
	configPath := filepath.Join(kubeDir, "config")
	kubeconfig := []byte(`apiVersion: v1
kind: Config
current-context: test-context
contexts:
- name: test-context
  context:
    cluster: test
    namespace: kube-ns
    user: test-user
clusters:
- name: test
  cluster:
    server: https://example.com
users:
- name: test-user
  user:
    token: fake
`)
	if err := os.WriteFile(configPath, kubeconfig, 0o644); err != nil {
		t.Fatalf("failed to write kubeconfig: %v", err)
	}

	t.Setenv("HOME", tmp)
	if ns := GetEffectiveNamespace("", false, "", false); ns != "kube-ns" {
		t.Errorf("expected namespace from kubeconfig 'kube-ns', got %q", ns)
	}
}

func TestGetEffectiveNamespace_Default(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	if ns := GetEffectiveNamespace("", false, "custom", false); ns != "custom" {
		t.Errorf("expected default namespace 'custom', got %q", ns)
	}
}
