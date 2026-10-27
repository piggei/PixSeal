package watermark

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"
)

type experimentalV4Build101Result struct {
	payload []byte
	info    ExperimentalV4ExtractInfo
	tele    experimentalV4PhoneBuild98RecoveryTelemetry
	ok      bool
}

func experimentalV4Build101Semantic(payload []byte, info ExperimentalV4ExtractInfo, tele experimentalV4PhoneBuild98RecoveryTelemetry, err error) experimentalV4Build101Result {
	experimentalV4Build100ClearDurations(reflect.ValueOf(&tele).Elem())
	return experimentalV4Build101Result{append([]byte(nil), payload...), info, tele, err == nil}
}

func experimentalV4Build101Equal(a, b experimentalV4Build101Result) bool {
	return bytes.Equal(a.payload, b.payload) && a.ok == b.ok && reflect.DeepEqual(a.tele, b.tele) && (!a.ok || reflect.DeepEqual(a.info, b.info))
}

func TestExperimentalV4Build101ObservedBankExact(t *testing.T) {
	carrier, _, _, w, h, _, boundary := experimentalV4Build100Fixture(t)
	for _, candidate := range []bool{false, true} {
		var want []experimentalV4PhoneHypothesis
		var we, ws, ww int
		var wg experimentalV4PhoneBuild98GeometryTelemetry
		var wf experimentalV4PhoneBuild75FreezeTelemetry
		if candidate {
			var cg experimentalV4PhoneBuild100GeometryTelemetry
			want, we, ws, ww, cg, wf = experimentalV4PhoneBuild100BlindBank(carrier, boundary, w, h)
			wg = cg.experimentalV4PhoneBuild98GeometryTelemetry
		} else {
			want, we, ws, ww, wg, wf = experimentalV4PhoneBuild98BlindBank(carrier, boundary, w, h)
		}
		observer := &experimentalV4Build101Observer{}
		got, ge, gs, gw, gg, gf := experimentalV4Build101ObservedBank(carrier, boundary, w, h, candidate, observer)
		pg := gg.experimentalV4PhoneBuild98GeometryTelemetry
		for _, v := range []any{&pg, &wg, &gf, &wf} {
			experimentalV4Build100ClearDurations(reflect.ValueOf(v).Elem())
		}
		if !reflect.DeepEqual(got, want) || ge != we || gs != ws || gw != ww || !reflect.DeepEqual(pg, wg) || !reflect.DeepEqual(gf, wf) {
			t.Fatalf("observed candidate=%t changed full bank/order/counters", candidate)
		}
		if observer.stage != "" || len(observer.stages) == 0 || observer.stages[0].name != "freeze" {
			t.Fatal("unbalanced or missing probes")
		}
	}
}

func TestExperimentalV4Build101ObservedRecoveryExact(t *testing.T) {
	f := experimentalV4Build93EmbeddedFixture(t, ProfileCapacity)
	b := PrintBoundaryEstimate{Detected: true, Confidence: 1, TopLeft: f.hyp.quad[0], TopRight: f.hyp.quad[1], BottomLeft: f.hyp.quad[2], BottomRight: f.hyp.quad[3]}
	wp, wi, wt, we := experimentalV4PhoneBuild98Recover(f.image, b, f.key, f.width, f.height)
	want := experimentalV4Build101Semantic(wp, wi, wt, we)
	for _, candidate := range []bool{false, true} {
		o := &experimentalV4Build101Observer{}
		p, i, tele, err := experimentalV4Build101ObservedRecover(f.image, b, f.key, f.width, f.height, candidate, o)
		if !experimentalV4Build101Equal(want, experimentalV4Build101Semantic(p, i, tele.experimentalV4PhoneBuild98RecoveryTelemetry, err)) || o.stage != "" {
			t.Fatalf("observed recovery candidate=%t differs", candidate)
		}
	}
}

func TestExperimentalV4Build101RowsAndSemanticGuard(t *testing.T) {
	s := experimentalV4Build101Stage{name: "gen4", before: experimentalV4Build101SnapshotNow(), after: experimentalV4Build101SnapshotNow()}
	if len(s.fields(1, "phone-b-angle.jpg", "build98")) != len(experimentalV4Build101StageColumns) {
		t.Fatal("stage row shape")
	}
	a := experimentalV4Build101Result{payload: []byte("fixture"), ok: true}
	b := a
	b.tele.Build98GeometryProfile.Gen3FetchBlockFailed++
	if experimentalV4Build101Equal(a, b) {
		t.Fatal("accepted changed semantic counter")
	}
	b = a
	b.payload = []byte("changed")
	if experimentalV4Build101Equal(a, b) {
		t.Fatal("accepted changed payload")
	}
}

// The timing pass invokes pinned production/candidate APIs without stage probes.
// A separate attribution pass uses test-only structural copies, stage CPU labels
// and heap/GC snapshots. It cannot qualify performance or promote Build100.
func TestExperimentalV4Build101RetainedStudy(t *testing.T) {
	dir, key, cw, ch, expected, cooldown := experimentalV4Build101Config(t)
	out := os.Getenv("PIXSEAL_BUILD101_OUTPUT_DIR")
	if out == "" || !filepath.IsAbs(out) {
		t.Fatal("absolute PIXSEAL_BUILD101_OUTPUT_DIR required")
	}
	fmt.Printf("build101-header\tpass\trun\tcycle\torder\tposition\tvariant\t%s\n", strings.Join(experimentalV4Build99RowColumns, "\t"))
	fmt.Printf("build101-stage-header\t%s\n", strings.Join(experimentalV4Build101StageColumns, "\t"))
	references := map[string]experimentalV4Build101Result{}
	run := 0
	for _, pass := range []string{"timing", "attribution"} {
		scenarios := []struct {
			cycle   int
			name    string
			reverse bool
		}{
			{1, "phone-b-mild.jpg", false}, {1, "phone-b-angle.jpg", false},
			{2, "phone-b-angle.jpg", true}, {2, "phone-b-mild.jpg", true},
		}
		if pass == "attribution" {
			scenarios = scenarios[:2]
			scenarios[1].reverse = true
		}
		for _, sc := range scenarios {
			work, boundary, prep := experimentalV4Build94Prepare(t, dir, sc.name, cw, ch)
			variants := []string{"build98", "build100"}
			order := "98-100"
			if sc.reverse {
				variants = []string{"build100", "build98"}
				order = "100-98"
			}
			for position, variant := range variants {
				run++
				runtime.GC()
				time.Sleep(cooldown)
				fmt.Printf("build101-start\t%s\t%d\t%s\t%s\n", pass, run, sc.name, variant)
				var p []byte
				var info ExperimentalV4ExtractInfo
				var tele experimentalV4PhoneBuild98RecoveryTelemetry
				var err error
				o := &experimentalV4Build101Observer{}
				var profile *os.File
				if pass == "attribution" {
					profile, err = os.OpenFile(filepath.Join(out, fmt.Sprintf("run%02d-%s-%s.cpu.prof", run, sc.name, variant)), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
					if err != nil {
						t.Fatal(err)
					}
					if err = pprof.StartCPUProfile(profile); err != nil {
						_ = profile.Close()
						t.Fatal(err)
					}
				}
				started := time.Now()
				if pass == "attribution" {
					var observed experimentalV4PhoneBuild100RecoveryTelemetry
					p, info, observed, err = experimentalV4Build101ObservedRecover(work, boundary, key, cw, ch, variant == "build100", o)
					tele = observed.experimentalV4PhoneBuild98RecoveryTelemetry
				} else if variant == "build98" {
					p, info, tele, err = experimentalV4PhoneBuild98Recover(work, boundary, key, cw, ch)
				} else {
					var candidate experimentalV4PhoneBuild100RecoveryTelemetry
					p, info, candidate, err = experimentalV4PhoneBuild100Recover(work, boundary, key, cw, ch)
					tele = candidate.experimentalV4PhoneBuild98RecoveryTelemetry
				}
				wall := time.Since(started)
				if profile != nil {
					pprof.StopCPUProfile()
					if e := profile.Close(); e != nil {
						t.Fatal(e)
					}
				}
				experimentalV4Build99AssertExpected(t, sc.name, p, tele, err, expected)
				semantic := experimentalV4Build101Semantic(p, info, tele, err)
				if reference, ok := references[sc.name]; ok {
					if !experimentalV4Build101Equal(reference, semantic) {
						t.Fatalf("run %d %s/%s full semantic mismatch", run, pass, variant)
					}
				} else {
					references[sc.name] = semantic
				}
				fields := []string{pass, strconv.Itoa(run), strconv.Itoa(sc.cycle), order, strconv.Itoa(position + 1), variant}
				fields = append(fields, experimentalV4Build99RowFields(sc.name, prep, wall, tele, err)...)
				fmt.Printf("build101-row\t%s\n", strings.Join(fields, "\t"))
				for _, stage := range o.stages {
					fmt.Printf("build101-stage\t%s\n", strings.Join(stage.fields(run, sc.name, variant), "\t"))
				}
				if o.stage != "" {
					t.Fatal("unbalanced retained probes")
				}
			}
		}
	}
}
