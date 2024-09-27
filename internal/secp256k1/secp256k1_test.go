package secp256k1_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/golang-jwt/jwt/v5"
	secp256k1 "github.com/watchedsky-social/core/internal/secp256k1"
)

func TestVerification(t *testing.T) {

	key, err := crypto.HexToECDSA(TestKey)
	if err != nil {
		t.Fatalf("failed parsing key: %s", err)
	}

	t.Run("ES256K", func(t *testing.T) {
		t.Parallel()
		var failureCount int

		for _, sString := range ES256K_Data {
			_, err = jwt.Parse(sString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*secp256k1.SigningMethodSecp256k1); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return &key.PublicKey, nil
			})

			if err != nil {
				failureCount++
				t.Errorf("%s: %s", sString, err)
			}
		}

		t.Logf("failed %d of %d checks", failureCount, len(ES256K_Data))
	})

	t.Run("ES256K-R", func(t *testing.T) {
		t.Parallel()
		var failureCount int

		for _, sString := range ES256KR_Data {
			_, err := jwt.Parse(sString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*secp256k1.SigningMethodSecp256k1); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return &key.PublicKey, nil
			})

			if err != nil {
				failureCount++
				t.Errorf("%s: %s", sString, err)
			}
		}

		t.Logf("failed %d of %d checks", failureCount, len(ES256K_Data))
	})

}

func TestGeneration(t *testing.T) {

	key, err := crypto.HexToECDSA(TestKey)
	if err != nil {
		t.Fatalf("failed parsing key: %s", err)
	}

	for i := range int64(128) {

		t.Run("ES256K/N="+fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			token := jwt.NewWithClaims(secp256k1.SigningMethodES256K, jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Unix(i, 0)),
			})

			sString, err := token.SignedString(key)
			if err != nil {
				t.Fatalf("failed signing token: %s", err)
			}

			_, err = jwt.Parse(sString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*secp256k1.SigningMethodSecp256k1); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return &key.PublicKey, nil
			})

			if err != nil {
				t.Fatalf("failed verifying signed token")
			}
		})

		t.Run("ES256K-R/N="+fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()
			token := jwt.NewWithClaims(secp256k1.SigningMethodES256KR, jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Unix(i, 0)),
			})

			sString, err := token.SignedString(key)
			if err != nil {
				t.Fatalf("failed signing token: %s", err)
			}

			_, err = jwt.Parse(sString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*secp256k1.SigningMethodSecp256k1); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return &key.PublicKey, nil
			})

			if err != nil {
				t.Fatalf("failed verifying signed token")
			}
		})
	}

}
