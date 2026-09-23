// Package jpeglegacy contains PixSeal's deterministic JPEG decoder.
//
// The implementation is derived from the pre-Go-1.26 standard-library
// image/jpeg decoder. PixSeal vendors the decoder so JPEG rasterization does
// not change when the host/build toolchain changes its standard library.
//
// The upstream Go source is covered by the BSD-style license in LICENSE_GO.
package jpeglegacy

// DecoderID is the stable identifier used by PixSeal documentation/tests for
// the deterministic legacy JPEG raster contract.
const DecoderID = "pixseal-jpeg-pre-go1.26-v1"
