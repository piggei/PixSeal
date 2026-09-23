# PixSeal deterministic JPEG decoder

This directory vendors the decode side of the pre-Go-1.26 Go standard-library
`image/jpeg` implementation, plus the YCbCr drawing helper it used internally.
The checked-in source snapshot was taken from Go 1.23.2, whose decoded raster
for the canonical Build43 A/mild photograph was experimentally identical to
Go 1.25.1. Build44 locks that legacy raster contract explicitly instead of
relying on whichever standard-library JPEG decoder ships with the compiler.

Build43 physical qualification exposed a reproducibility boundary: the same
JPEG bytes rasterize differently with Go 1.25.1 and Go 1.26.0. PixSeal's phone
geometry is intentionally sensitive to image evidence, so changing the raster
can change geometry ranking before HMAC is ever attempted.

Build44 therefore makes JPEG rasterization part of the PixSeal input contract
instead of inheriting it from the compiler's standard library. PNG continues to
use `image/png`.

The vendored files retain the Go Authors copyright headers. The upstream
BSD-style license is copied as `LICENSE_GO`.
