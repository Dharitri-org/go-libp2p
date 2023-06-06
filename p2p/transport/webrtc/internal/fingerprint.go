package internal

import (
	"crypto"
	"crypto/x509"
	"errors"

	"github.com/libp2p/go-libp2p/p2p/transport/webrtc/internal/encoding"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
	"github.com/pion/webrtc/v3"
)

// Fingerprint is forked from pion to avoid bytes to string alloc,
// and to avoid the entire hex interspersing when we do not need it anyway

var (
	errHashUnavailable = errors.New("fingerprint: hash algorithm is not linked into the binary")
)

// Fingerprint creates a fingerprint for a certificate using the specified hash algorithm
func Fingerprint(cert *x509.Certificate, algo crypto.Hash) ([]byte, error) {
	if !algo.Available() {
		return nil, errHashUnavailable
	}
	h := algo.New()
	// Hash.Writer is specified to be never returning an error.
	// https://golang.org/pkg/hash/#Hash
	h.Write(cert.Raw)
	return h.Sum(nil), nil
}

func DecodeRemoteFingerprint(maddr ma.Multiaddr) (*mh.DecodedMultihash, error) {
	remoteFingerprintMultibase, err := maddr.ValueForProtocol(ma.P_CERTHASH)
	if err != nil {
		return nil, err
	}
	_, data, err := multibase.Decode(remoteFingerprintMultibase)
	if err != nil {
		return nil, err
	}
	return mh.Decode(data)
}

func EncodeDTLSFingerprint(fp webrtc.DTLSFingerprint) (string, error) {
	digest, err := encoding.DecodeInterspersedHexFromASCIIString(fp.Value)
	if err != nil {
		return "", err
	}
	encoded, err := mh.Encode(digest, mh.SHA2_256)
	if err != nil {
		return "", err
	}
	return multibase.Encode(multibase.Base64url, encoded)
}
