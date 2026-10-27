package watermark

import (
	"os"
	"strconv"
	"testing"
	"time"
)

func experimentalV4Build101Config(t *testing.T) (string, []byte, int, int, string, time.Duration) {
	t.Helper()
	acquisitionDir := os.Getenv("PIXSEAL_BUILD101_ACQUISITION_DIR")
	if acquisitionDir == "" {
		t.Skip("Build101 private retained profile requires PIXSEAL_BUILD101_ACQUISITION_DIR")
	}
	key := []byte(os.Getenv("PIXSEAL_BUILD101_KEY"))
	if len(key) < 8 {
		t.Fatal("Build101 PIXSEAL_BUILD101_KEY must contain at least 8 bytes")
	}
	cw, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD101_WIDTH"))
	if err != nil || cw <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD101_WIDTH: %v", err)
	}
	ch, err := strconv.Atoi(os.Getenv("PIXSEAL_BUILD101_HEIGHT"))
	if err != nil || ch <= 0 {
		t.Fatalf("invalid PIXSEAL_BUILD101_HEIGHT: %v", err)
	}
	expectedPayload := os.Getenv("PIXSEAL_BUILD101_MESSAGE_B")
	if expectedPayload == "" {
		t.Fatal("PIXSEAL_BUILD101_MESSAGE_B is required")
	}
	cooldownMs := 5000
	if raw := os.Getenv("PIXSEAL_BUILD101_COOLDOWN_MS"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 || v > 60000 {
			t.Fatalf("invalid PIXSEAL_BUILD101_COOLDOWN_MS=%q", raw)
		}
		cooldownMs = v
	}
	return acquisitionDir, key, cw, ch, expectedPayload, time.Duration(cooldownMs) * time.Millisecond
}
