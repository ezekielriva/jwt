package jwt_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"

	. "github.com/gbrlsnchs/jwt/v2"
)

type testCase struct {
	signer        Signer
	verifier      Signer
	marshalingErr error
	signingErr    error
	parsingErr    error
	unmarshalErr  error
	verifyingErr  error
}

type testToken struct {
	*JWT
	Name      string  `json:"name,omitempty"`
	RandInt   int     `json:"randomInt,omitempty"`
	RandFloat float64 `json:"randomFloat,omitempty"`
}

func testJWT(t *testing.T, testCases []testCase) {
	for i, tc := range testCases {
		name := fmt.Sprintf("%s %s", tc.signer.String(), tc.verifier.String())
		t.Run(name, func(t *testing.T) {
			now := time.Now()
			kid := fmt.Sprintf("kid %s %d", t.Name(), i)
			iat := now.Unix()
			exp := now.Add(30 * time.Minute).Unix()
			nbf := now.Add(1 * time.Second).Unix()
			iss := fmt.Sprintf("%s %d", t.Name(), i)
			aud := fmt.Sprintf("test %d", i)
			sub := fmt.Sprintf("sub %d", i)
			jti := strconv.Itoa(i)
			randomInt := rand.Intn(int(^uint32(0)))
			randomFloat := rand.Float64() * 100
			jot := &testToken{
				JWT: &JWT{
					Header: &Header{
						Algorithm: tc.signer.String(),
						KeyID:     kid,
					},
					Claims: &Claims{
						IssuedAt:   iat,
						Expiration: exp,
						NotBefore:  nbf,
						Issuer:     iss,
						Audience:   aud,
						Subject:    sub,
						ID:         jti,
					},
				},
				Name:      name,
				RandInt:   randomInt,
				RandFloat: randomFloat,
			}

			// 1 - Marshal.
			payload, err := Marshal(jot)
			if want, got := tc.marshalingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.SkipNow()
			}

			// 2 - Sign.
			token, err := tc.signer.Sign(payload)
			if want, got := tc.signingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.SkipNow()
			}

			// 3 - Parse.
			payload, sig, err := ParseBytes(token)
			if want, got := tc.parsingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.SkipNow()
			}

			// 4 - Unmarshal.
			var jot2 testToken
			err = Unmarshal(payload, &jot2)
			if want, got := tc.unmarshalErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.SkipNow()
			}

			// 5 - Verify.
			err = tc.verifier.Verify(payload, sig)
			if want, got := tc.verifyingErr, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if err != nil {
				t.SkipNow()
			}

			// 6 - Check new token.
			if want, got := tc.signer.String(), jot2.Header.Algorithm; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := kid, jot2.Header.KeyID; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := iat, jot2.IssuedAt; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := exp, jot2.Expiration; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := nbf, jot2.NotBefore; want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := iss, jot2.Issuer; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := aud, jot2.Audience; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := sub, jot2.Subject; want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			if want, got := jti, jot2.ID; want != got {
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

			if want, got := reflect.ValueOf(jot.Header).Elem().NumField(), reflect.ValueOf(jot2.Header).Elem().NumField(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}

			if want, got := reflect.ValueOf(jot.Claims).Elem().NumField(), reflect.ValueOf(jot2.Claims).Elem().NumField(); want != got {
				t.Errorf("want %d, got %d", want, got)
			}
		})
	}
}
