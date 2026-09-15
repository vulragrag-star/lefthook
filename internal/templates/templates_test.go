package templates

import (
	"strings"
	"testing"
)

func TestHookQuotesPathsWithSpaces(t *testing.T) {
	t.Parallel()

	hook := string(Hook("pre-commit", Args{
		Rc:           "/tmp/my config/lefthook-rc.sh",
		LefthookPath: "/opt/my tools/lefthook",
		Roots:        []string{"apps/web"},
	}))

	checks := []string{
		`[ -f "/tmp/my config/lefthook-rc.sh" ] && . "/tmp/my config/lefthook-rc.sh"`,
		`elif test -n "/opt/my tools/lefthook"`,
		`"/opt/my tools/lefthook" "$@"`,
	}
	for _, want := range checks {
		if !strings.Contains(hook, want) {
			t.Fatalf("generated hook missing quoted path fragment %q\n----\n%s", want, hook)
		}
	}

	// Unquoted path with a space must not appear as a bare word.
	bads := []string{
		`[ -f /tmp/my config/lefthook-rc.sh ]`,
		`. /tmp/my config/lefthook-rc.sh`,
		`elif /opt/my tools/lefthook -h`,
		"\t/opt/my tools/lefthook \"$@\"",
	}
	for _, bad := range bads {
		if strings.Contains(hook, bad) {
			t.Fatalf("generated hook still contains unquoted path fragment %q\n----\n%s", bad, hook)
		}
	}

	// os.Executable() path used for LefthookPathCurrent must be quoted on both check and invoke.
	for _, line := range strings.Split(hook, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "elif ") && strings.Contains(trim, " -h >/dev/null") && !strings.Contains(trim, "lefthook") && !strings.Contains(trim, "go tool") && !strings.Contains(trim, "bundle") && !strings.Contains(trim, "yarn") && !strings.Contains(trim, "pnpm") && !strings.Contains(trim, "swift") && !strings.Contains(trim, "mint") && !strings.Contains(trim, "uv ") && !strings.Contains(trim, "mise") && !strings.Contains(trim, "devbox") && !strings.Contains(trim, "test -n") && !strings.Contains(trim, "test -f") && !strings.Contains(trim, "command -v") {
			if !strings.HasPrefix(trim, `elif "`) {
				t.Fatalf("LefthookPathCurrent check not quoted: %s", trim)
			}
		}
	}
}
