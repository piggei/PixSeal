package watermark

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// Exercise the actual disk -> fresh runtime -> Gen4 -> semantic gate path in
// both conditions without requiring private images or a long retained run.
func TestExperimentalV4Build103FreshProcessSmoke(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("live process CPU measurement requires Linux")
	}
	// Match the child runtime when the caller has fewer CPUs or GOMAXPROCS=1.
	previous := runtime.GOMAXPROCS(8)
	defer runtime.GOMAXPROCS(previous)
	plane, pilot, w, h, anchor, root := experimentalV4Build97Fixture(t)
	f := experimentalV4Build102Frozen{plane: plane, pilot: pilot, cw: w, ch: h, anchor: anchor, tasks: []experimentalV4PhoneHypothesis{root, root}, seeds: 2, workers: 2}
	bank, evals, seeds, workers, gp, fp := experimentalV4Build102Replay(f)
	s := experimentalV4Build103Capture(f)
	s.Image, s.BankHash, s.Bank = "synthetic", experimentalV4Build102Hash(bank), len(bank)
	s.EvalsTotal, s.SeedsTotal, s.WorkersTotal = evals, seeds, workers
	s.Expected, s.ExpectedFreeze = experimentalV4Build103Geo(gp.experimentalV4PhoneBuild98GeometryTelemetry), fp
	path := filepath.Join(t.TempDir(), "synthetic.gob")
	if err := experimentalV4Build103Write(path, s); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(b))
	hz, err := exec.Command("getconf", "CLK_TCK").Output()
	if err != nil {
		t.Fatal(err)
	}
	binary, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, variant := range []string{"default", "async_off"} {
		t.Run(variant, func(t *testing.T) {
			debug := "asyncpreemptoff=0"
			if variant == "async_off" {
				debug = "asyncpreemptoff=1"
			}
			values := map[string]string{"GODEBUG": debug, "GOMAXPROCS": "8", "GOGC": "100", "GOMEMLIMIT": "off", "GOTRACEBACK": "single", "PIXSEAL_BUILD103_SNAPSHOT": path, "PIXSEAL_BUILD103_SNAPSHOT_SHA256": hash, "PIXSEAL_BUILD103_VARIANT": variant, "PIXSEAL_BUILD103_PASS": "timing", "PIXSEAL_BUILD103_RUN": "1", "PIXSEAL_BUILD103_BLOCK": "1", "PIXSEAL_BUILD103_PAIR": "1", "PIXSEAL_BUILD103_POSITION": "1", "PIXSEAL_BUILD103_CLOCK_TICKS": strings.TrimSpace(string(hz)), "PIXSEAL_BUILD103_COOLDOWN_MS": "0"}
			var env []string
			for _, value := range os.Environ() {
				key, _, _ := strings.Cut(value, "=")
				if _, replaced := values[key]; !replaced {
					env = append(env, value)
				}
			}
			for key, value := range values {
				env = append(env, key+"="+value)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, binary, "-test.run=^TestExperimentalV4Build103ReplayProcess$", "-test.v", "-test.timeout=0")
			cmd.Env = env
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("fresh process failed: %v\n%s", err, output)
			}
			if !strings.Contains(string(output), "build103-row\tsynthetic\t") || !strings.Contains(string(output), "\t"+variant+"\t"+debug+"\t") || !strings.Contains(string(output), "\tPASS\n") {
				t.Fatalf("missing semantic row: %s", output)
			}
		})
	}
}
