package httpsignatures

import (
	"net/http"
	"strings"
	"testing"
)

const digestBodyExample = `{"hello": "world"}`
const digestHostExample = "https://example.com"

var getDigestRequestFunc = func(b string, h string) *http.Request {
	r, _ := http.NewRequest(http.MethodPost, digestHostExample, strings.NewReader(b))
	r.Header.Set("Digest", h)
	return r
}

type testAlg struct{}

func (a testAlg) Algorithm() string {
	return "TEST"
}

func (a testAlg) Create(data []byte) ([]byte, error) {
	return []byte{}, nil
}

func (a testAlg) Verify(data []byte, digest []byte) error {
	return nil
}

func TestVerifyDigest(t *testing.T) {
	var cryptoErrType = "*httpsignatures.DigestError"
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantErr     bool
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Valid MD5 digest",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "MD5=Sd/dVLAcvNLSq16eXua5uQ=="),
			},
			want: true,
		},
		{
			name: "Valid SHA-256 digest",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "SHA-256=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE="),
			},
			want: true,
		},
		{
			name: "Valid SHA-512 digest",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "SHA-512=WZDPaVn/7XgHaAy8pmojAkGWoRx2UFChF41A2svX+TaPm+AbwAgBWnrIiYllu7BNNyealdVLvRwEmTHWXvJwew=="),
			},
			want: true,
		},
		{
			name: "Invalid MD5 digest (decode error)",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "MD5=123456"),
			},
			want:        false,
			wantErr:     true,
			wantErrType: cryptoErrType,
			wantErrMsg:  "error decode digest from base64: illegal base64 data at input byte 4",
		},
		{
			name: "Invalid MD5 wrong digest",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "MD5=X48E9qOokqqrvdts8nOJRJN3OWDUoyWxBf7kbu9DBPE="),
			},
			want:        false,
			wantErr:     true,
			wantErrType: cryptoErrType,
			wantErrMsg:  "wrong digest: wrong hash",
		},
		{
			name: "Invalid digest header",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "SHA-512="),
			},
			want:        false,
			wantErr:     true,
			wantErrType: cryptoErrType,
			wantErrMsg:  "digest parser error: empty digest value",
		},
		{
			name: "Unsupported digest hash algorithm",
			args: args{
				r: getDigestRequestFunc(digestBodyExample, "SHA-0=test"),
			},
			want:        false,
			wantErr:     true,
			wantErrType: cryptoErrType,
			wantErrMsg:  "unsupported digest hash algorithm 'SHA-0'",
		},
		{
			name: "Empty body",
			args: args{
				r: getDigestRequestFunc("", "MD5=xxx"),
			},
			want:        false,
			wantErr:     true,
			wantErrType: cryptoErrType,
			wantErrMsg:  "empty body",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			err := d.Verify(tt.args.r)
			got := err == nil
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErr, tt.wantErrMsg)
		})
	}
}

func TestDigestSetDigestHashAlgorithm(t *testing.T) {
	tests := []struct {
		name string
		arg  DigestHashAlgorithm
	}{
		{
			name: "Set new algorithm OK",
			arg:  testAlg{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDigest()
			d.SetDigestHashAlgorithm(tt.arg)
			if _, ok := d.alg["TEST"]; ok == false {
				t.Error("algorithm not found")
			}
		})
	}
}
