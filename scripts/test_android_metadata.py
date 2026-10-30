import importlib.util
from pathlib import Path
import sys
import tempfile
import unittest

sys.dont_write_bytecode = True

spec = importlib.util.spec_from_file_location("android_metadata", Path(__file__).with_name("android_metadata.py"))
metadata = importlib.util.module_from_spec(spec)
spec.loader.exec_module(metadata)


class AndroidMetadataTest(unittest.TestCase):
    def test_source_drift_is_rejected(self):
        source = Path(__file__).resolve().parent.parent
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "gui").mkdir()
            for file in ("VERSION", "gui/FyneApp.toml", "gui/AndroidManifest.xml"):
                (root / file).write_bytes((source / file).read_bytes())
            version, build = metadata.source_metadata(root)
            for file in ("gui/FyneApp.toml", "gui/AndroidManifest.xml"):
                original = (root / file).read_text()
                (root / file).write_text(original.replace(build, str(int(build) + 1)))
                with self.subTest(file=file), self.assertRaises(ValueError):
                    metadata.source_metadata(root)
                (root / file).write_text(original)

    def test_wrong_package_version_architecture_and_permissions_are_rejected(self):
        valid = f"package: name='{metadata.APP_ID}' versionCode='106' versionName='0.3.0'\nnative-code: 'arm64-v8a'\n"
        metadata.check_package(valid, "", "0.3.0", "106")
        for badging, permissions in ((valid.replace("106", "105"), ""), (valid.replace(metadata.APP_ID, "wrong.app"), ""), (valid.replace("arm64-v8a", "x86"), ""), (valid, "uses-permission: name='android.permission.INTERNET'")):
            with self.subTest(badging=badging, permissions=permissions), self.assertRaises(ValueError):
                metadata.check_package(badging, permissions, "0.3.0", "106")


if __name__ == "__main__":
    unittest.main()
