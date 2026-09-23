package watermark

import "testing"

func TestExperimentalV4Build45Classification(t *testing.T) {
	cases := []struct {
		name string
		in   ExperimentalV4PhoneInfo
		want ExperimentalV4PhoneBuild45Classification
	}{
		{"recovered", ExperimentalV4PhoneInfo{HMACAuthenticated: true}, ExperimentalV4PhoneBuild45Recovered},
		{"data", ExperimentalV4PhoneInfo{Accepted: true}, ExperimentalV4PhoneBuild45DataChannel},
		{"qualification", ExperimentalV4PhoneInfo{Build43Attempted: true, Build43FrozenCandidates: 16, Build43QualifiedCandidates: 1}, ExperimentalV4PhoneBuild45Qualification},
		{"geometry-no-freeze", ExperimentalV4PhoneInfo{Build43Attempted: true}, ExperimentalV4PhoneBuild45Geometry},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExperimentalV4PhoneBuild45Classify(tc.in); got != tc.want {
				t.Fatalf("classify=%q want %q", got, tc.want)
			}
		})
	}
}
