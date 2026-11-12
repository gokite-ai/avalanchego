"""Check real binary versions and handshake before accepting a candidate pair."""
import json
import re
import sys
from pathlib import Path

def verify(commit, node, plugin, handshake):
    if not re.fullmatch(r"[0-9a-f]{40}", commit):
        raise ValueError("Expected a full source commit")
    for name, output in (("node", node), ("plugin", plugin)):
        if commit not in output:
            raise ValueError(f"{name}: wrong source commit")
        if not re.search(r"(?:avalanchego/|Subnet-EVM/v)1\.15\.3\b", output):
            raise ValueError(f"{name}: expected version 1.15.3")
        if not re.search(r"rpcchainvm[= :]+47\b", output):
            raise ValueError(f"{name}: expected RPCChainVM protocol 47")
    if "PASS: actual Subnet-EVM subprocess handshake; rpcchainvm=47" not in handshake:
        raise ValueError("Real subprocess handshake did not pass")
    return {"source_commit": commit, "platform": "linux/amd64",
            "version": "1.15.3", "rpcchainvm": 47, "handshake": "passed",
            "database_restore": "not_tested", "database_downgrade": "not_tested"}

if __name__ == "__main__":
    directory = Path(sys.argv[2])
    result = verify(sys.argv[1],
                    (directory / "node-version.txt").read_text(),
                    (directory / "plugin-version.txt").read_text(),
                    (directory / "handshake.txt").read_text())
    (directory / "verification.json").write_text(json.dumps(result, indent=2) + "\n")
