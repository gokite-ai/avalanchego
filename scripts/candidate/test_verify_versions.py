import unittest
from verify_versions import verify

COMMIT = "a" * 40
NODE = f"avalanchego/1.15.3 [database=v1.4.5, rpcchainvm=47, commit={COMMIT}]"
PLUGIN = f"Subnet-EVM/v1.15.3@{COMMIT} [rpcchainvm=47]"
HANDSHAKE = "PASS: actual Subnet-EVM subprocess handshake; rpcchainvm=47"

class CandidateEvidenceTests(unittest.TestCase):
    def test_matching_pair(self):
        self.assertEqual(verify(COMMIT, NODE, PLUGIN, HANDSHAKE)["handshake"], "passed")

    def test_rejects_wrong_commit_old_protocol_version_and_missing_handshake(self):
        for args in (
            (COMMIT[:8], NODE, PLUGIN, HANDSHAKE),
            (COMMIT, NODE.replace(COMMIT, "b" * 40), PLUGIN, HANDSHAKE),
            (COMMIT, NODE, PLUGIN.replace("=47", "=46"), HANDSHAKE),
            (COMMIT, NODE.replace("1.15.3", "1.15.1"), PLUGIN, HANDSHAKE),
            (COMMIT, NODE, PLUGIN, ""),
        ):
            with self.subTest(args=args):
                with self.assertRaises(ValueError):
                    verify(*args)

if __name__ == "__main__":
    unittest.main()
