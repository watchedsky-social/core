package feed

import (
	"context"
	"errors"
	"fmt"

	"github.com/bluesky-social/indigo/atproto/crypto"
	"github.com/bluesky-social/indigo/atproto/identity"
	"github.com/bluesky-social/indigo/atproto/syntax"
	"github.com/golang-jwt/jwt/v5"

	// don't forget to import so the init function is called
	_ "github.com/watchedsky-social/core/internal/k256"
)

var directory identity.Directory = identity.DefaultDirectory()

func VerifyAuth(ctx context.Context, jwtStr, serviceDID, nsid string) (jwt.Claims, error) {

	parser := jwt.NewParser(jwt.WithAudience(serviceDID), jwt.WithJSONNumber())
	claims := jwt.MapClaims{}
	_, err := parser.ParseWithClaims(jwtStr, claims, func(t *jwt.Token) (interface{}, error) {
		issuer, err := t.Claims.GetIssuer()
		if err != nil {
			return nil, err
		}

		issuerDID, err := syntax.ParseDID(issuer)
		if err != nil {
			return nil, err
		}

		doc, err := directory.LookupDID(ctx, issuerDID)
		if err != nil {
			return nil, err
		}

		key, ok := doc.Keys["atproto"]
		if !ok {
			return nil, errors.New("missing atproto key")
		}

		return crypto.ParsePublicMultibase(key.PublicKeyMultibase)
	})

	if err != nil {
		return nil, err
	}

	lxm, ok := claims["lxm"]
	if !ok {
		return nil, errors.New("lxm missing")
	}

	if fmt.Sprint(lxm) != nsid {
		return nil, fmt.Errorf("invalid lxm %s", lxm)
	}

	return claims, nil
}
