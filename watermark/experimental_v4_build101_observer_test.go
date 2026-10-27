package watermark

import (
	"context"
	"fmt"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"time"
)

// These probes are test-only. ReadMemStats introduces a synchronization cost,
// so observed runs are attribution evidence, never promotion timing evidence.
type experimentalV4Build101Snapshot struct {
	mem                                      runtime.MemStats
	userCPU, gcCPU, scavengeCPU              float64
	schedEvents, schedOver1ms, schedOver10ms uint64
}

func experimentalV4Build101SnapshotNow() experimentalV4Build101Snapshot {
	var s experimentalV4Build101Snapshot
	runtime.ReadMemStats(&s.mem)
	values := []metrics.Sample{
		{Name: "/cpu/classes/user:cpu-seconds"},
		{Name: "/cpu/classes/gc/total:cpu-seconds"},
		{Name: "/cpu/classes/scavenge/total:cpu-seconds"},
		{Name: "/sched/latencies:seconds"},
	}
	metrics.Read(values)
	s.userCPU = values[0].Value.Float64()
	s.gcCPU = values[1].Value.Float64()
	s.scavengeCPU = values[2].Value.Float64()
	h := values[3].Value.Float64Histogram()
	for i, n := range h.Counts {
		s.schedEvents += n
		if h.Buckets[i] >= 0.001 {
			s.schedOver1ms += n
		}
		if h.Buckets[i] >= 0.01 {
			s.schedOver10ms += n
		}
	}
	return s
}

type experimentalV4Build101Stage struct {
	name          string
	wall          time.Duration
	before, after experimentalV4Build101Snapshot
}

type experimentalV4Build101Observer struct {
	stage   string
	started time.Time
	before  experimentalV4Build101Snapshot
	stages  []experimentalV4Build101Stage
}

func (o *experimentalV4Build101Observer) begin(stage string) {
	if o == nil {
		return
	}
	if o.stage != "" {
		panic("nested Build101 observation")
	}
	pprof.SetGoroutineLabels(pprof.WithLabels(context.Background(), pprof.Labels("stage", stage)))
	o.stage = stage
	o.before = experimentalV4Build101SnapshotNow()
	o.started = time.Now()
}

func (o *experimentalV4Build101Observer) end() {
	if o == nil {
		return
	}
	if o.stage == "" {
		panic("unbalanced Build101 observation")
	}
	wall := time.Since(o.started)
	after := experimentalV4Build101SnapshotNow()
	o.stages = append(o.stages, experimentalV4Build101Stage{o.stage, wall, o.before, after})
	o.stage = ""
	pprof.SetGoroutineLabels(context.Background())
}

var experimentalV4Build101StageColumns = []string{
	"run", "image", "variant", "stage", "wall_ms", "heap_before_bytes", "heap_after_bytes",
	"allocated_bytes", "mallocs", "gc_cycles", "gc_pause_ns", "estimated_user_cpu_s", "estimated_gc_cpu_s",
	"estimated_scavenge_cpu_s", "sched_events", "sched_ge1ms", "sched_ge10ms",
}

func (s experimentalV4Build101Stage) fields(run int, image, variant string) []string {
	b, a := s.before, s.after
	return []string{
		fmt.Sprint(run), image, variant, s.name, fmt.Sprint(s.wall.Milliseconds()),
		fmt.Sprint(b.mem.HeapAlloc), fmt.Sprint(a.mem.HeapAlloc),
		fmt.Sprint(a.mem.TotalAlloc - b.mem.TotalAlloc), fmt.Sprint(a.mem.Mallocs - b.mem.Mallocs),
		fmt.Sprint(a.mem.NumGC - b.mem.NumGC), fmt.Sprint(a.mem.PauseTotalNs - b.mem.PauseTotalNs),
		fmt.Sprintf("%.6f", a.userCPU-b.userCPU), fmt.Sprintf("%.6f", a.gcCPU-b.gcCPU),
		fmt.Sprintf("%.6f", a.scavengeCPU-b.scavengeCPU), fmt.Sprint(a.schedEvents - b.schedEvents),
		fmt.Sprint(a.schedOver1ms - b.schedOver1ms), fmt.Sprint(a.schedOver10ms - b.schedOver10ms),
	}
}
