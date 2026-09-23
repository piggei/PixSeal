#!/usr/bin/env python3
"""Build45 laboratory-only reference-assisted geometry diagnostic.

This helper is NOT used by PixSeal production decoding. It uses OpenCV SIFT
features from the known digital marked carrier and the acquired phone photo to
estimate one homography, then writes the mapped artwork quadrilateral in the
acquisition's original pixel coordinates. The quad can be supplied only to
`pixseal v4-diagnose-phone -oracle-quad-json ...`.
"""
import argparse
import json
import sys
from pathlib import Path

try:
    import cv2
    import numpy as np
except Exception as exc:
    print(f"error: Python OpenCV + NumPy are required for the lab oracle: {exc}", file=sys.stderr)
    sys.exit(2)


def scaled_gray(path, max_dim):
    img = cv2.imread(str(path), cv2.IMREAD_GRAYSCALE)
    if img is None:
        raise RuntimeError(f"cannot read image: {path}")
    h, w = img.shape[:2]
    scale = min(1.0, float(max_dim) / float(max(w, h)))
    if scale < 1.0:
        img = cv2.resize(img, (int(round(w * scale)), int(round(h * scale))), interpolation=cv2.INTER_AREA)
    return img, w, h, scale


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--reference", required=True, help="known digital marked carrier")
    ap.add_argument("--acquired", required=True, help="phone acquisition")
    ap.add_argument("--out", required=True, help="output JSON quad")
    ap.add_argument("--max-dim", type=int, default=2600, help="feature-analysis max dimension")
    ap.add_argument("--ratio", type=float, default=0.72, help="Lowe ratio")
    ap.add_argument("--ransac", type=float, default=5.0, help="RANSAC threshold in scaled acquisition pixels")
    args = ap.parse_args()

    ref, rw, rh, rs = scaled_gray(Path(args.reference), args.max_dim)
    acq, aw, ah, aas = scaled_gray(Path(args.acquired), args.max_dim)

    sift = cv2.SIFT_create(nfeatures=14000, contrastThreshold=0.02, edgeThreshold=12)
    kr, dr = sift.detectAndCompute(ref, None)
    ka, da = sift.detectAndCompute(acq, None)
    if dr is None or da is None or len(kr) < 12 or len(ka) < 12:
        raise RuntimeError("insufficient SIFT features")

    matcher = cv2.BFMatcher(cv2.NORM_L2)
    knn = matcher.knnMatch(dr, da, k=2)
    good = [m for m, n in knn if m.distance < args.ratio * n.distance]
    if len(good) < 12:
        raise RuntimeError(f"insufficient ratio-test matches: {len(good)}")

    src = np.float32([kr[m.queryIdx].pt for m in good]).reshape(-1, 1, 2)
    dst = np.float32([ka[m.trainIdx].pt for m in good]).reshape(-1, 1, 2)
    H, mask = cv2.findHomography(src, dst, cv2.RANSAC, args.ransac)
    if H is None or mask is None:
        raise RuntimeError("homography estimation failed")
    inliers = int(mask.ravel().sum())
    if inliers < 10:
        raise RuntimeError(f"too few RANSAC inliers: {inliers}")

    # TL, TR, BL, BR in the scaled reference image. Convert mapped acquisition
    # coordinates back to the acquisition's original raster.
    corners = np.float32([[[0.0, 0.0]], [[rw * rs - 1.0, 0.0]], [[0.0, rh * rs - 1.0]], [[rw * rs - 1.0, rh * rs - 1.0]]])
    mapped = cv2.perspectiveTransform(corners, H).reshape(4, 2)
    mapped[:, 0] /= aas
    mapped[:, 1] /= aas

    out = {
        "method": "opencv-sift-ransac-reference-lab-only",
        "reference": str(Path(args.reference)),
        "acquired": str(Path(args.acquired)),
        "reference_size": [rw, rh],
        "acquired_size": [aw, ah],
        "reference_analysis_scale": rs,
        "acquired_analysis_scale": aas,
        "keypoints_reference": len(kr),
        "keypoints_acquired": len(ka),
        "matches": len(good),
        "inliers": inliers,
        "inlier_fraction": inliers / len(good),
        "quad": [{"x": float(x), "y": float(y)} for x, y in mapped],
    }
    Path(args.out).parent.mkdir(parents=True, exist_ok=True)
    Path(args.out).write_text(json.dumps(out, indent=2) + "\n", encoding="utf-8")
    print(f"Build45 lab registration: matches={len(good)} inliers={inliers} ({out['inlier_fraction']:.3f})")
    print(f"wrote {args.out}")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"error: {exc}", file=sys.stderr)
        sys.exit(1)
