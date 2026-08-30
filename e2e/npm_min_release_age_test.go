package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestNPMMinReleaseAge(t *testing.T) {
	tests := []struct {
		name       string
		scenario   string
		args       []string
		wantOutput []string
		wantErr    bool
	}{
		{
			name:     "latest version is safe",
			scenario: "safe",
			wantOutput: []string{
				"✅ express: current=4.18.0 latest=5.2.1",
				"SAFE",
			},
		},
		{
			name:     "finds latest safe version",
			scenario: "fallback",
			args:     []string{"14"},
			wantOutput: []string{
				"⏳ lodash: current=4.17.20 latest=4.18.1",
				"NOT SAFE | latest safe: 4.17.21",
			},
		},
		{
			name:     "reports no safe version",
			scenario: "no-safe",
			args:     []string{"14"},
			wantOutput: []string{
				"⏳ react: current=18.2.0 latest=19.0.0",
				"NOT SAFE | no safe version found",
			},
		},
		{
			name:       "reports up to date",
			scenario:   "empty",
			wantOutput: []string{"All packages are up to date."},
		},
		{
			name:     "reports npm failure",
			scenario: "failure",
			wantErr:  true,
		},
	}

	binary := buildBinary(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeNPM := writeFakeNPM(t)
			env := append(os.Environ(), "NPM_E2E_SCENARIO="+tt.scenario)
			env = append(env, "PATH="+filepath.Dir(fakeNPM)+string(os.PathListSeparator)+os.Getenv("PATH"))

			cmd := exec.Command(binary, tt.args...)
			cmd.Dir = t.TempDir()
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected command to fail, output: %s", output)
				}
				if !strings.Contains(string(output), "Error checking outdated packages:") {
					t.Fatalf("unexpected error output: %s", output)
				}
				return
			}
			if err != nil {
				t.Fatalf("command failed: %v\noutput: %s", err, output)
			}
			for _, want := range tt.wantOutput {
				if !strings.Contains(string(output), want) {
					t.Errorf("output does not contain %q:\n%s", want, output)
				}
			}
		})
	}
}

func buildBinary(t *testing.T) string {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine test source path")
	}
	root := filepath.Dir(filepath.Dir(sourceFile))
	binary := filepath.Join(t.TempDir(), "npm-min-release-age")
	cmd := exec.Command("go", "build", "-o", binary, "./cmd/npm-min-release-age")
	cmd.Dir = root
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("building binary: %v\n%s", err, output)
	}
	return binary
}

func writeFakeNPM(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "npm")
	script := `#!/bin/sh
set -eu

case "${NPM_E2E_SCENARIO:-}" in
failure)
  printf '%s\n' 'not json'
  exit 1
  ;;
empty)
  printf '%s\n' '{}'
  ;;
safe)
  case " $* " in
  *" outdated --json "*) printf '%s\n' '{"express":{"current":"4.18.0","latest":"5.2.1"}}' ;;
  *) printf '%s\n' '{"5.2.1":"2020-01-01T00:00:00Z"}' ;;
  esac
  ;;
fallback)
  case " $* " in
  *" outdated --json "*) printf '%s\n' '{"lodash":{"current":"4.17.20","latest":"4.18.1"}}' ;;
  *) printf '%s\n' '{"4.18.1":"2099-01-01T00:00:00Z","4.17.21":"2020-01-01T00:00:00Z"}' ;;
  esac
  ;;
no-safe)
  case " $* " in
  *" outdated --json "*) printf '%s\n' '{"react":{"current":"18.2.0","latest":"19.0.0"}}' ;;
  *) printf '%s\n' '{"19.0.0":"2099-01-01T00:00:00Z"}' ;;
  esac
  ;;
esac
`
	if err := os.WriteFile(path, []byte(script), 0755); err != nil {
		t.Fatalf("writing fake npm: %v", err)
	}
	return path
}
