package bench_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/lestrrat-go/jwx/internal/jwxtest"
	"github.com/lestrrat-go/jwx/jwk"
)

func runJSONBench(b *testing.B, privkey jwk.Key) {
	b.Helper()

	privkey.Set("mykey", "1234567890")
	pubkey, err := jwk.PublicKeyOf(privkey)
	if err != nil {
		b.Fatal(err)
	}

	keytypes := []struct {
		Name string
		Key  jwk.Key
	}{
		{Name: "PublicKey", Key: pubkey},
		{Name: "PrivateKey", Key: privkey},
	}

	for _, keytype := range keytypes {
		key := keytype.Key
		b.Run(keytype.Name, func(b *testing.B) {
			buf, _ := json.Marshal(key)
			s := string(buf)
			rdr := bytes.NewReader(buf)

			testcases := []Case{
				{
					Name: "jwk.Parse",
					Test: func(b *testing.B) error {
						_, err := jwk.Parse(buf)
						return err
					},
				},
				{
					Name: "jwk.ParseString",
					Test: func(b *testing.B) error {
						_, err := jwk.ParseString(s)
						return err
					},
				},
				{
					Name: "jwk.ParseReader",
					Pretest: func(b *testing.B) error {
						b.StopTimer()
						rdr.Seek(0, 0)
						b.StartTimer()
						return nil
					},
					Test: func(b *testing.B) error {
						_, err := jwk.ParseReader(rdr)
						return err
					},
				},
				{
					Name: "json.Marshal",
					Test: func(b *testing.B) error {
						_, err := json.Marshal(key)
						return err
					},
				},
			}
			for _, tc := range testcases {
				tc.Run(b)
			}
		})
	}
}

func BenchmarkJWK(b *testing.B) {
	b.Run("Serialization", func(b *testing.B) {
		b.Run("RSA", func(b *testing.B) {
			rsakey, _ := jwxtest.GenerateRsaJwk()
			runJSONBench(b, rsakey)
		})
		b.Run("EC", func(b *testing.B) {
			eckey, _ := jwxtest.GenerateEcdsaJwk()
			runJSONBench(b, eckey)
		})
		b.Run("Symmetric", func(b *testing.B) {
			symkey, _ := jwxtest.GenerateSymmetricJwk()
			runJSONBench(b, symkey)
		})
	})
}
