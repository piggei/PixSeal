package watermark

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// experimentalV4ActiveCorpusPaths returns only the explicitly manifested
// Build32 private corpus. Large entries stay part of the corpus identity but
// are skipped before image decoding when they exceed the configured Go-test
// budget.
func experimentalV4ActiveCorpusPaths(t *testing.T) []string {
	t.Helper()
	directory := os.Getenv("PIXSEAL_V4_CORPUS_DIR")
	if directory == "" {
		t.Skip("set PIXSEAL_V4_CORPUS_DIR to run the local experimental v4 corpus")
	}
	manifest := os.Getenv("PIXSEAL_V4_CORPUS_MANIFEST")
	if manifest == "" {
		manifest = "private-corpus-active.tsv"
	}
	maxMP := 50
	if raw := os.Getenv("PIXSEAL_V4_CORPUS_MAX_MPIX"); raw != "" {
		if n, err := strconv.Atoi(raw); err == nil && n >= 0 {
			maxMP = n
		}
	}
	f, err := os.Open(manifest)
	if err != nil {
		t.Fatalf("open active corpus manifest: %v", err)
	}
	defer f.Close()
	var paths []string
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := strings.TrimSpace(scan.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 5 {
			t.Fatalf("invalid active corpus row: %q", line)
		}
		w, err1 := strconv.Atoi(fields[2])
		h, err2 := strconv.Atoi(fields[3])
		if err1 != nil || err2 != nil {
			t.Fatalf("invalid active corpus dimensions: %q", line)
		}
		if maxMP > 0 && int64(w)*int64(h) > int64(maxMP)*1000000 {
			t.Logf("SKIP active corpus %s (%s): %.1f MP exceeds Go-test budget %d MP", fields[0], fields[1], float64(w*h)/1e6, maxMP)
			continue
		}
		p := filepath.Join(directory, fields[1])
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("active corpus file %s: %v", p, err)
		}
		paths = append(paths, p)
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	if len(paths) == 0 {
		t.Fatal("active corpus has no files within Go-test budget")
	}
	return paths
}
