package httpsignatures

import (
	"reflect"
	"testing"
	"time"
)

const parserErrType = "*httpsignatures.ParserError"
const validSignatureHeader = `keyId="Test",algorithm="rsa-sha256",created=1402170695,expires=1402170699,headers="(request-target) (created) (expires) host date digest content-length",signature="vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE="`
const validAuthorizationHeader = `Signature ` + validSignatureHeader

var validParsedSignatureHeader = ParsedHeader{
	keyID:     "Test",
	algorithm: "rsa-sha256",
	created:   time.Unix(1402170695, 0),
	expires:   time.Unix(1402170699, 0),
	headers:   []string{"(request-target)", "(created)", "(expires)", "host", "date", "digest", "content-length"},
	signature: "vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE=",
}

var validHeadersIfNotSpecified = []string{"(created)"}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Parser
	}{
		{
			name: "Successful",
			want: &Parser{
				params: make(map[string]bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParserParseSingleFields(t *testing.T) {
	type args struct {
		header        string
		authorization bool
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Authorization: Empty header",
			args: args{
				header:        ``,
				authorization: true,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: empty header",
		},
		{
			name: "Authorization: Only Signature keyword",
			args: args{
				header:        `Signature`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Only Signature keyword with space",
			args: args{
				header:        `Signature  `,
				authorization: true,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected parameter",
		},
		{
			name: "Authorization: Wrong in keyword",
			args: args{
				header:        `Auth`,
				authorization: true,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: invalid Authorization header, must start from Signature keyword",
		},
		{
			name: "Authorization: Wrong in keyword with space char",
			args: args{
				header:        `Auth `,
				authorization: true,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: invalid Authorization header, must start from Signature keyword",
		},
		{
			name: "Authorization: Signature and keyID",
			args: args{
				header:        `Signature keyId="v1"`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				keyID:   "v1",
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and algorithm",
			args: args{
				header:        `Signature algorithm="v2"`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				algorithm: "v2",
				headers:   validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and headers",
			args: args{
				header:        `Signature headers="(request-target) (created)" `,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				headers: []string{"(request-target)", "(created)"},
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and signature param",
			args: args{
				header:        `Signature signature="test" `,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				signature: "test",
				headers:   validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and all params without spaces",
			args: args{
				header:        `Signature keyId="v1",algorithm="v2",created=1402170695,expires=1402170699,headers="v-3 v-4",signature="v5"`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				keyID:     "v1",
				algorithm: "v2",
				created:   time.Unix(1402170695, 0),
				expires:   time.Unix(1402170699, 0),
				headers:   []string{"v-3", "v-4"},
				signature: "v5",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and all params and extra spaces",
			args: args{
				header:        `Signature   keyId  ="v1", algorithm  ="v2",created = 1402170695, expires = 1402170699 , headers  =  "  v-3 v-4  ", signature="v5"   `,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				keyID:     "v1",
				algorithm: "v2",
				created:   time.Unix(1402170695, 0),
				expires:   time.Unix(1402170699, 0),
				headers:   []string{"v-3", "v-4"},
				signature: "v5",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Signature: all params",
			args: args{
				header:        `keyId="v1",algorithm="v2",created=1402170695,expires=1402170699,headers="v-3 v-4",signature="v5"`,
				authorization: false,
			},
			want: ParsedHeader{
				keyID:     "v1",
				algorithm: "v2",
				created:   time.Unix(1402170695, 0),
				expires:   time.Unix(1402170699, 0),
				headers:   []string{"v-3", "v-4"},
				signature: "v5",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Unsupported symbol in key",
			args: args{
				header: `keyId-="v1"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: found '-' — unsupported symbol in key",
		},
		{
			name: "Unsupported symbol, expected = symbol",
			args: args{
				header: `keyId :"v1"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: found ':' — unsupported symbol, expected '=' or space symbol",
		},
		{
			name: "Unsupported symbol, expected quote symbol",
			args: args{
				header: `keyId= 'v1'`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: found ''' — unsupported symbol, expected '\"' or space symbol",
		},
		{
			name: "Unknown parameter",
			args: args{
				header: `key="v1"`,
			},
			want: ParsedHeader{
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "unexpected end of header, expected equal symbol",
			args: args{
				header: `keyId`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected '=' symbol and field value",
		},
		{
			name: "Expected field value",
			args: args{
				header: `keyId `,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected field value",
		},
		{
			name: "Expected quote",
			args: args{
				header: `keyId= `,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected '\"' symbol and field value",
		},
		{
			name: "Expected quote at the end",
			args: args{
				header: `keyId="`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected '\"' symbol",
		},
		{
			name: "Empty value",
			args: args{
				header: `keyId=""`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: empty value for key 'keyId'",
		},
		{
			name: "Div symbol expected",
			args: args{
				header: `keyId="v1" algorithm="v2"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: found 'a' — unsupported symbol, expected ',' or space symbol",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			if true == tt.args.authorization {
				p.flag = "keyword"
			} else {
				p.flag = "param"
			}
			var got, err = p.parseSignature(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseCreatedExpires(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Authorization: Signature and created",
			args: args{
				header: `created=1402170695`,
			},
			want: ParsedHeader{
				created: time.Unix(1402170695, 0),
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Authorization: Signature and expires",
			args: args{
				header: `expires=1402170699`,
			},
			want: ParsedHeader{
				expires: time.Unix(1402170699, 0),
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Wrong created INT value",
			args: args{
				header: `created=9223372036854775807`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'created' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
		{
			name: "Wrong created INT value with space at the end",
			args: args{
				header: `created=9223372036854775808 `,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'created' param value: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
		},
		{
			name: "Wrong created INT value with divider",
			args: args{
				header: `created=9223372036854775809,`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'created' param value: strconv.ParseInt: parsing \"9223372036854775809\": value out of range",
		},
		{
			name: "Wrong expires INT value",
			args: args{
				header: `expires=9223372036854775807`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
		{
			name: "Wrong expires with space at the end",
			args: args{
				header: `expires=9223372036854775808 `,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775808\": value out of range",
		},
		{
			name: "Wrong expires with divider",
			args: args{
				header: `expires=9223372036854775809,`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775809\": value out of range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			p.flag = "param"
			var got, err = p.parseSignature(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseAuthorization(t *testing.T) {
	type args struct {
		header        string
		authorization bool
	}
	var validAuthParsedHeader = validParsedSignatureHeader
	validAuthParsedHeader.keyword = "Signature"
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Authorization",
			args: args{
				header:        validAuthorizationHeader,
				authorization: true,
			},
			want:        validAuthParsedHeader,
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.ParseAuthorizationHeader(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseSignature(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Signature",
			args: args{
				header: validSignatureHeader,
			},
			want:        validParsedSignatureHeader,
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.ParseSignatureHeader(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseSignatureFailed(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Current parser stage not setKeyValue",
			args: args{
				header: `keyId="Test"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected parser stage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.parseSignature(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseDigestFailed(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedDigestHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Current parser stage not set",
			args: args{
				header: `MD5=test`,
			},
			want:        ParsedDigestHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected parser stage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.parseDigest(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseAmbiguousParams(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Ambiguous Parameters",
			args: args{
				header: `keyId="v1",ambiguous="v2",digest="v3"`,
			},
			want: ParsedHeader{
				keyID:   "v1",
				headers: validHeadersIfNotSpecified,
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.ParseSignatureHeader(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseDuplicateParams(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "Duplicate keyID",
			args: args{
				header: `keyId="v1",keyId="v2"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'keyId'",
		},
		{
			name: "Duplicate algorithm",
			args: args{
				header: `algorithm="v1",algorithm="v2"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'algorithm'",
		},
		{
			name: "Duplicate created",
			args: args{
				header: `created=1402170695,created=1402170695`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'created'",
		},
		{
			name: "Duplicate expires",
			args: args{
				header: `expires=1402170699,expires=1402170699`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'expires'",
		},
		{
			name: "Duplicate headers",
			args: args{
				header: `headers="v1",headers="v2"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'headers'",
		},
		{
			name: "Duplicate signature",
			args: args{
				header: `signature="v1",signature="v2"`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: duplicate param 'signature'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			p.flag = "param"
			got, err := p.parseSignature(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestNotSpecifiedHeadersParams(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "No headers",
			args: args{
				header: `keyId="v1"`,
			},
			want: ParsedHeader{
				keyID:   "v1",
				headers: []string{"(created)"},
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "Empty headers",
			args: args{
				header: `keyId="v1",headers=""`,
			},
			want:        ParsedHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: empty value for key 'headers'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.ParseSignatureHeader(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestParserParseDigestHeader(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        ParsedDigestHeader
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "MD5 Digest",
			args: args{
				header: `MD5=ZDk5NTk4ODgxNjM3MDc5MDQ2MTgzNDQwMzExMThiZWI=`,
			},
			want: ParsedDigestHeader{
				algo:   "MD5",
				digest: "ZDk5NTk4ODgxNjM3MDc5MDQ2MTgzNDQwMzExMThiZWI=",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "SHA-1 Digest",
			args: args{
				header: `SHA-1=ZDNiMDlhYmUzMGNmZTJlZGZmNGVlOWUwYTE0MWM5M2JmNWIzYWY4Nw==`,
			},
			want: ParsedDigestHeader{
				algo:   "SHA-1",
				digest: "ZDNiMDlhYmUzMGNmZTJlZGZmNGVlOWUwYTE0MWM5M2JmNWIzYWY4Nw==",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name: "SHA-256 Digest",
			args: args{
				header: `SHA-256=NWY4ZjA0ZjZhM2E4OTJhYWFiYmRkYjZjZjI3Mzg5NDQ5Mzc3Mzk2MGQ0YTMyNWIxMDVmZWU0NmVlZjQzMDRmMQ==`,
			},
			want: ParsedDigestHeader{
				algo:   "SHA-256",
				digest: "NWY4ZjA0ZjZhM2E4OTJhYWFiYmRkYjZjZjI3Mzg5NDQ5Mzc3Mzk2MGQ0YTMyNWIxMDVmZWU0NmVlZjQzMDRmMQ==",
			},
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
		{
			name:        "Empty Digest header",
			args:        args{},
			want:        ParsedDigestHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: empty digest header",
		},
		{
			name: "Empty Digest value",
			args: args{
				header: `md5`,
			},
			want:        ParsedDigestHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: unexpected end of header, expected digest value",
		},
		{
			name: "Unsupported digest algorithm symbol",
			args: args{
				header: `md 5=`,
			},
			want:        ParsedDigestHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: found ' ' — unsupported symbol in algorithm",
		},
		{
			name: "Empty digest value",
			args: args{
				header: `MD5=`,
			},
			want:        ParsedDigestHeader{},
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: empty digest value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			var got, err = p.ParseDigestHeader(tt.args.header)
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}

func TestVerifySignatureFields(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantErrType string
		wantErrMsg  string
	}{
		{
			name: "No keyID",
			args: args{
				header: `algorithm="v1",headers="v-2",signature="v3"`,
			},
			want:        false,
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: keyId is not set in header",
		},
		{
			name: "No signature",
			args: args{
				header: `keyId="v1",headers="v-2"`,
			},
			want:        false,
			wantErrType: parserErrType,
			wantErrMsg:  "ParserError: signature is not set in header",
		},
		{
			name: "OK",
			args: args{
				header: validSignatureHeader,
			},
			want:        true,
			wantErrType: parserErrType,
			wantErrMsg:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			_, _ = p.ParseSignatureHeader(tt.args.header)
			err := p.VerifySignatureFields()
			got := err == nil
			assert(t, got, err, tt.wantErrType, tt.name, tt.want, tt.wantErrMsg)
		})
	}
}
