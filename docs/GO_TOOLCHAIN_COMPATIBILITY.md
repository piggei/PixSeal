# Go toolchain and JPEG reproducibility

## Build43 qualified toolchain

PixSeal v0.3.0-build43 is qualified with **Go 1.25.1**.

The repository Makefile selects that exact toolchain explicitly through
`GOTOOLCHAIN=go1.25.1`. `make build` always rebuilds `dist/pixseal`, so an older
binary compiled with a different Go release cannot be silently reused by a
physical qualification target.

Use the normal commands:

```sh
make
make v4-build43-phone-physical-test
make all-test ALL_TEST_REPORT=report_build_43.txt
```

Do not prefix those commands with a different `GOTOOLCHAIN`. The Makefile owns
the qualified selection. A machine with a newer Go installed may download/use
Go 1.25.1 automatically through Go's toolchain mechanism.

On Windows, `build-windows.bat` sets `GOTOOLCHAIN=go1.25.1` before testing and
building.

## Why Go 1.26 is not qualified for Build43

Go 1.26 replaced the standard-library JPEG encoder and decoder with new
implementations. The Go 1.26 release notes explicitly warn that programs that
expect bit-for-bit JPEG encoder/decoder output may need to be updated:

https://go.dev/doc/go1.26#image/jpeg

PixSeal's physical smartphone geometry is intentionally sensitive to small
image-domain differences. During Build43 qualification, the canonical private
`phone-a-mild.jpg` was decoded with the same source code on the same Linux host
using Go 1.25.1 and Go 1.26.0. The original JPEG bytes were identical, but the
rasterized Y/Cb/Cr planes differed.

For the canonical A/mild photograph (source SHA-256
`e2993c77aef3ec917d564b674f1250072bb786dcd69cb0e7f6719d51cffe0403`):

```text
Go 1.25.1
Y  e7635474904ae14d0ed2a2fcae6ef2305993120cfe050e13ce5df3b0d9fc541e
Cb dd2b5afb6ecca2a013bbb2414dacf2f379cc9b886af424d39cf38319f63544f6
Cr ce188a4c30f6e827f50b15cfd277f7702aa633112ee2f6458289459b2a8d9e23

Go 1.26.0
Y  21bf5b4092e030629e9f6c78f89a7292abeeeb4de0f5af7f63cb7aa088f7ab6d
Cb 865839f3869d6a2943a122e21d23000732ecea5b79f563e8ace3ed0c3ec6dc9e
Cr 687fdd7ce7b7a768593a3641721d783a019dfc72d85681fb8d2de7b66b329eec
```

That upstream raster difference changes the boundary/pilot geometry scores
enough that the qualified A/mild path can reject under Go 1.26.0. Rebuilding the
same Build43 source with Go 1.25.1 restores the qualified path and HMAC-authenticated
payload `v4-b38-phone-a`.

This is **not** evidence of a Format-v4, HMAC, ECC or Build43 side-pair algorithm
regression. It is a toolchain/input-decoder reproducibility issue discovered by
the physical corpus.

## Build43 policy

- Go 1.25.1 is the only qualified Build43 build/test toolchain.
- Go 1.26+ is intentionally not claimed as physically equivalent for JPEG input.
- PNG behavior does not remove the need for this qualification rule because the
  physical smartphone corpus is JPEG.
- The Format-v3/v4 on-image formats are unchanged.
- The Build43 geometry algorithm is not retuned to compensate for Go 1.26's
  different JPEG rasterization.

## Follow-up after Build43

A later build should remove this toolchain dependency by defining and testing a
stable image-ingest contract. Candidate approaches include a project-controlled
JPEG decoding path or another normalization step whose decoded raster is covered
by deterministic fixtures/hashes across supported platforms. That change must be
qualified separately before Go 1.26+ becomes a supported physical-build toolchain.

## Build44 deterministic-ingest qualification

Build44 implements the follow-up described above. JPEG decoding is now owned by
PixSeal through `internal/jpeglegacy`, with public deterministic Y/Cb/Cr fixture
hashes and a CLI-ingest regression. PNG remains on `image/png`.

Build44 no longer relies on the toolchain standard library for JPEG rasterization. Two explicit qualification targets were used to validate Go 1.26.0:

```sh
make v4-build44-go126-jpeg-compat-test
make v4-build44-go126-phone-physical-test
```

The first checks the frozen raster contract under Go 1.26.0. The second builds a
Go 1.26.0 executable and requires the complete private Build43 physical matrix
to remain unchanged. Both gates passed on the Surface/WSL2 qualification host.
Build44 therefore promotes **Go 1.26.0** as the normal qualified toolchain.
`make build`, `make all-test` and `build-windows.bat` now select Go 1.26.0.
Go 1.25.1 remains relevant only as the historical Build43 qualification reference.
