"""Check real binary versions and handshake before accepting a candidate pair."""
import json
import re
import sys
from pathlib import Path

def verify(commit, node, plugin, handshake):
    if not re.fullmatch(r"[0-9a-f]{40}", commit):
        raise ValueError("Expected a full source commit")
    node_format = (
        rf"avalanchego/1\.15\.3 \[database=[^,\]\s]+, rpcchainvm=47, "
        rf"commit={commit}(?:, go=[^,\]\s]+)?\]"
    )
    if not re.fullmatch(node_format, node.strip()):
        raise ValueError("Unexpected AvalancheGo version, commit or protocol")
    if plugin.strip() != f"Subnet-EVM/v1.15.3@{commit} [rpcchainvm=47]":
        raise ValueError("Unexpected Subnet-EVM version, commit or protocol")
    if handshake.strip() != "PASS: actual Subnet-EVM subprocess handshake; rpcchainvm=47":
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
