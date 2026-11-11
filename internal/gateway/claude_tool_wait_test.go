package gateway

import (
	"os"
	"strings"
	"testing"
)

// A turn's run waits on magpie's MCP calls for as long as the caller's tool
// takes: Claude Code neither moves one still running after 120s to the
// background nor aborts one silent for 30 minutes. A value the user set in
// magpie's environment stays as they set it.
func TestClaudeRunWaitsOnToolCalls(t *testing.T) {
	for _, c := range []struct {
		name string
		env  map[string]string
		want string
	}{
		{name: "nothing set", want: "background 0 idle 0"},
		{name: "magpie's env", env: map[string]string{"CLAUDE_CODE_MCP_AUTO_BACKGROUND_MS": "600000"}, want: "background 600000 idle 0"},
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, kv := range claudeToolWaitVars {
				k, _, _ := strings.Cut(kv, "=")
				t.Setenv(k, "")
				os.Unsetenv(k)
			}
			for k, v := range c.env {
				t.Setenv(k, v)
			}
			h := newScriptHarness(t, "ok", `echo "background ${CLAUDE_CODE_MCP_AUTO_BACKGROUND_MS-unset} idle ${CLAUDE_CODE_MCP_TOOL_IDLE_TIMEOUT-unset}" >> $log`)
			h.ask("rules", "hello", "next")
			runs := h.runs()
			if len(runs) != 1 {
				t.Fatalf("runs: %d", len(runs))
			}
			for _, log := range runs {
				if !strings.Contains(log, c.want) {
					t.Fatalf("run's environment, want %q:\n%s", c.want, log)
				}
			}
		})
	}
}
