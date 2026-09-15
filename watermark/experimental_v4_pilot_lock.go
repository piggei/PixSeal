package watermark

// Build30 locks the exact Build24 prototype-2 pilot identity against accidental
// mutation while the rest of Format v4 remains experimental. This is a
// development candidate lock, not yet a normative on-image Format-v4 contract:
// physical v4 print-camera/scanner evidence and the v4 framing/payload codec are
// still required before normative promotion.
const (
	experimentalV4PilotCandidateLockName = "prototype-2-search-p64"
	experimentalV4PilotCandidateLockHash = "858f74305ee9a9cbb59dd3fb6ab8afc6aaf8958e4f9f517711e2c52fc053b174"
)
