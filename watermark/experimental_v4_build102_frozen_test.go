package watermark

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Test-only immutable capture at the original Gen3 -> Gen4 barrier.
type experimentalV4Build102Frozen struct {
	plane                 *pixelPlane
	pilot                 experimentalV4PilotCandidate
	cw, ch                int
	anchor                [4]ImagePoint
	tasks                 []experimentalV4PhoneHypothesis
	evals, seeds, workers int
	profile               experimentalV4PhoneBuild100GeometryTelemetry
	freeze                experimentalV4PhoneBuild75FreezeTelemetry
}

func experimentalV4Build102Pack(plane *pixelPlane, pilot experimentalV4PilotCandidate, cw, ch int, anchor [4]ImagePoint, tasks []experimentalV4PhoneHypothesis, evals, seeds, workers int, profile experimentalV4PhoneBuild100GeometryTelemetry, freeze experimentalV4PhoneBuild75FreezeTelemetry) experimentalV4Build102Frozen {
	return experimentalV4Build102Frozen{plane, pilot, cw, ch, anchor, tasks, evals, seeds, workers, profile, freeze}
}

// Reflection reads unexported geometry fields directly, without Interface or
// pointer addresses. Float bits, nil markers, types, lengths and original order
// are included. JSON would omit these unexported fields and accept empty hashes.
func experimentalV4Build102Hash(value any) string {
	h := sha256.New()
	experimentalV4Build102HashValue(h, reflect.ValueOf(value))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func experimentalV4Build102HashValue(h hash.Hash, v reflect.Value) {
	put := func(n uint64) { var b [8]byte; binary.LittleEndian.PutUint64(b[:], n); _, _ = h.Write(b[:]) }
	typeName := v.Type().String()
	put(uint64(len(typeName)))
	_, _ = h.Write([]byte(typeName))
	switch v.Kind() {
	case reflect.Struct:
		put(uint64(v.NumField()))
		for i := 0; i < v.NumField(); i++ {
			experimentalV4Build102HashValue(h, v.Field(i))
		}
	case reflect.Array, reflect.Slice:
		if v.Kind() == reflect.Slice {
			if v.IsNil() {
				put(0)
				return
			}
			put(1)
		}
		put(uint64(v.Len()))
		for i := 0; i < v.Len(); i++ {
			experimentalV4Build102HashValue(h, v.Index(i))
		}
	case reflect.Pointer:
		if v.IsNil() {
			put(0)
		} else {
			put(1)
			experimentalV4Build102HashValue(h, v.Elem())
		}
	case reflect.Float64:
		put(math.Float64bits(v.Float()))
	case reflect.Float32:
		put(uint64(math.Float32bits(float32(v.Float()))))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		put(uint64(v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		put(v.Uint())
	case reflect.String:
		s := v.String()
		put(uint64(len(s)))
		_, _ = h.Write([]byte(s))
	case reflect.Bool:
		if v.Bool() {
			put(1)
		} else {
			put(0)
		}
	default:
		panic("unsupported Build102 fingerprint kind: " + v.Kind().String())
	}
}

// Linux process CPU counters are live at the read boundary. CLK_TCK is supplied
// by getconf, never assumed. Timing rows explicitly describe their resolution.
func experimentalV4Build102ParseCPU(raw string) (uint64, uint64, error) {
	i := strings.LastIndex(raw, ")")
	if i < 0 {
		return 0, 0, fmt.Errorf("invalid process stat command field")
	}
	fields := strings.Fields(raw[i+1:])
	if len(fields) < 13 {
		return 0, 0, fmt.Errorf("short process stat")
	}
	u, err := strconv.ParseUint(fields[11], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	s, err := strconv.ParseUint(fields[12], 10, 64)
	return u, s, err
}

func experimentalV4Build102CPU() (uint64, uint64, error) {
	b, err := os.ReadFile("/proc/self/stat")
	if err != nil {
		return 0, 0, err
	}
	return experimentalV4Build102ParseCPU(string(b))
}
