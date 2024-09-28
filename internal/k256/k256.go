package k256

import (
	"fmt"

	"github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/golang-jwt/jwt/v5"
)

type signingMethodK256 struct{}

// Alg implements jwt.SigningMethod.
func (s *signingMethodK256) Alg() string {
	return "ES256K"
}

// Sign implements jwt.SigningMethod.
func (s *signingMethodK256) Sign(signingString string, key interface{}) ([]byte, error) {
	priv, ok := key.(crypto.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("wrong type %T", key)
	}

	return priv.HashAndSign([]byte(signingString))
}

// Verify implements jwt.SigningMethod.
func (s *signingMethodK256) Verify(signingString string, sig []byte, key interface{}) error {
	pub, ok := key.(crypto.PublicKey)
	if !ok {
		return fmt.Errorf("wrong type %T", key)
	}

	return pub.HashAndVerifyLenient([]byte(signingString), sig)
}

func init() {
	jwt.RegisterSigningMethod("ES256K", func() jwt.SigningMethod {
		return &signingMethodK256{}
	})
}
