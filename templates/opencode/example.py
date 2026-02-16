import os
from e2b import Sandbox

test_dir = os.path.dirname(os.path.abspath(__file__))
with open(os.path.join(test_dir, "test.sh")) as f:
    test_script = f.read()

sbx = Sandbox.create("opencode", timeout=60)
try:
    sbx.files.write("/tmp/test.sh", test_script)
    result = sbx.commands.run("bash /tmp/test.sh")
    assert result.exit_code == 0, f"Test failed:\n{result.stderr}"
    print(result.stdout)
finally:
    sbx.kill()
