package watermark

import "testing"

func TestExperimentalV4Build46HandoffClassification(t *testing.T) {
	tests := []struct {
		name                       string
		build41, frozen, qualified int
		single                     bool
		want                       ExperimentalV4PhoneBuild46Classification
	}{
		{"build41", 2, 0, 0, false, ExperimentalV4PhoneBuild46Build41Accepted},
		{"no-candidate", 0, 0, 0, false, ExperimentalV4PhoneBuild46NoBuild43Candidate},
		{"heldout", 0, 12, 0, false, ExperimentalV4PhoneBuild46HeldoutQualification},
		{"singleton-auth", 0, 32, 1, true, ExperimentalV4PhoneBuild46QualifiedEnsembleShortfall},
		{"singleton-mismatch", 0, 32, 1, false, ExperimentalV4PhoneBuild46QualifiedCandidateMismatch},
		{"bank-shortfall", 0, 32, 2, false, ExperimentalV4PhoneBuild46Build42BankShortfall},
		{"handoff", 0, 32, 3, false, ExperimentalV4PhoneBuild46QualifiedHandoffAvailable},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExperimentalV4PhoneBuild46ClassifyCounts(tc.build41, tc.frozen, tc.qualified, tc.single); got != tc.want {
				t.Fatalf("got %q want %q", got, tc.want)
			}
		})
	}
}
