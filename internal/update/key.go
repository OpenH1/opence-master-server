package update

// ReleaseKey is the public half of the release key (hex), which signs each
// release's SHA256SUMS (cmd/release-sign, in the release workflow). The
// updater installs nothing it does not sign; MASTER_UPDATE_KEY replaces it
// for a fork's releases.
//
// (OpenH1's fork: the master release key of OpenH1/opence-master-server's
// releases, not upstream's)
const ReleaseKey = "3607ccfaa3ebb0c1cbbfdb8b1578522cb281b250799f9f28b529c2bb1e6e16ca"
