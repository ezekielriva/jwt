package jwt_test

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	. "github.com/gbrlsnchs/jwt/v3"
)

type testCase struct {
	signer        Signer
	verifier      Verifier
	marshalingErr error
	signingErr    error
	parsingErr    error
	unmarshalErr  error
	verifyingErr  error
}

type testToken struct {
	JWT
	Name      string  `json:"name,omitempty"`
	RandInt   int     `json:"randomInt,omitempty"`
	RandFloat float64 `json:"randomFloat,omitempty"`
}

func testJWT(t *testing.T, testCases []testCase) {
	for i, tc := range testCases {
		t.Run(tc.signer.String(), func(t *testing.T) {
			now := time.Now()
			kid := fmt.Sprintf("kid %s %d", t.Name(), i)
			typ := "JWT"
			cty := "JWT"
			iat := now.Unix()
			exp := now.Add(30 * time.Minute).Unix()
			nbf := now.Add(1 * time.Second).Unix()
			iss := fmt.Sprintf("%s %d", t.Name(), i)
			aud := Audience{fmt.Sprintf("test %d", i)}
			sub := fmt.Sprintf("sub %d", i)
			jti := strconv.Itoa(i)
			randomInt := rand.Intn(math.MaxUint32)
			randomFloat := rand.Float64() * 100
			jot := &testToken{
				JWT: JWT{
					Header: Header{
						KeyID:       kid,
						ContentType: cty,
					},
					Claims: &Claims{
						Issuer:         iss,
						Subject:        sub,
						Audience:       aud,
						ExpirationTime: exp,
						NotBefore:      nbf,
						IssuedAt:       iat,
						ID:             jti,
					},
				},
				Name:      t.Name(),
				RandInt:   randomInt,
				RandFloat: randomFloat,
			}

			// Sign.
			token, err := Sign(jot, tc.signer)
			if want, got := tc.signingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				return
			}

			// Verify.
			var raw RawToken
			raw, err = Verify(token, tc.verifier)
			if want, got := tc.verifyingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				return
			}

			// Decode token.
			var jot2 testToken
			err = raw.Decode(&jot2)
			if want, got := tc.unmarshalErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}

			// Check new token.
			if want, got := tc.signer.String(), jot2.Algorithm; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := kid, jot2.KeyID; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := typ, jot2.Type; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := cty, jot2.ContentType; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := iat, jot2.IssuedAt; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := exp, jot2.ExpirationTime; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := nbf, jot2.NotBefore; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := iss, jot2.Issuer; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := aud, jot2.Audience; !reflect.DeepEqual(want, got) {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := sub, jot2.Subject; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := jti, jot2.ID; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := t.Name(), jot2.Name; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := randomInt, jot2.RandInt; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := randomFloat, jot2.RandFloat; want != got {
				t.Errorf("want %f, got %f", want, got)
			}

			if want, got := reflect.ValueOf(jot).Elem().NumField(), reflect.ValueOf(&jot2).Elem().NumField(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}
