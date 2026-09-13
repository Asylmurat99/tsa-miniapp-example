import hashlib
import hmac
import json
import sys
from pathlib import Path


# #region check
def check_signature(fields: dict, secret: str, sign_field: str) -> bool:
    """fields: decoded values; sign_field: "hash" for the launch context,
    "sign" for the getPhone envelope."""
    fields = dict(fields)
    received = fields.pop(sign_field, "")
    canon = "\n".join(f"{k}={fields[k]}" for k in sorted(fields))
    calc = hmac.new(secret.encode(), canon.encode(), hashlib.sha256).hexdigest()
    return hmac.compare_digest(calc, received)
# #endregion check


vector = json.loads((Path(__file__).parent / ".." / "vector.json").read_text())
failed = False
for case in vector["cases"]:
    ok = check_signature(case["fields"], vector["secret"], case["sign_field"])
    print(f"{case['name']}: {'ok' if ok else 'FAIL'}")
    failed = failed or not ok
sys.exit(1 if failed else 0)
