package httpsignatures

import (
	"reflect"
	"testing"
)

func TestCreate(t *testing.T) {
	tests := []struct {
		name string
		want *parser
	}{
		{
			name: "Successful",
			want: &parser{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Create(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParserParse(t *testing.T) {
	type args struct {
		header        string
		authorization bool
	}
	tests := []struct {
		name       string
		args       args
		want       ParsedHeader
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Authorization: Empty header",
			args: args{
				header:        ``,
				authorization: true,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "empty header",
		},
		{
			name: "Authorization: Only Signature keyword",
			args: args{
				header:        `Signature`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Only Signature keyword with space",
			args: args{
				header:        `Signature  `,
				authorization: true,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected end of header, expected key param",
		},
		{
			name: "Authorization: Wrong in keyword",
			args: args{
				header:        `Auth`,
				authorization: true,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "invalid Authorization header, must start from Signature keyword",
		},
		{
			name: "Authorization: Wrong in keyword with space char",
			args: args{
				header:        `Auth `,
				authorization: true,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "invalid Authorization header, must start from Signature keyword",
		},
		{
			name: "Authorization: Signature and keyId",
			args: args{
				header:        `Signature keyId="v1"`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				keyId:   "v1",
			},
			wantErr:    false,
			wantErrMsg: "",
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
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Signature and created",
			args: args{
				header:        `Signature created=1402170695`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				created: 1402170695,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Signature and expires",
			args: args{
				header:        `Signature expires=1402170699`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				expires: 1402170699,
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Signature and headers",
			args: args{
				header:        `Signature headers="(request-target) (created)" `,
				authorization: true,
			},
			want: ParsedHeader{
				keyword: "Signature",
				headers: []string{"(request-target)", "(created)",},
			},
			wantErr:    false,
			wantErrMsg: "",
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
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Signature and all params",
			args: args{
				header:        `Signature keyId="v1",algorithm="v2",created=1402170695,expires=1402170699,headers="v-3 v-4",signature="v5"`,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				keyId:     "v1",
				algorithm: "v2",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"v-3", "v-4",},
				signature: "v5",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Authorization: Signature and all params and extra spaces",
			args: args{
				header:        `Signature   keyId  ="v1", algorithm  ="v2",created = 1402170695, expires = 1402170699 , headers  =  "  v-3 v-4  ", signature="v5"   `,
				authorization: true,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				keyId:     "v1",
				algorithm: "v2",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"v-3", "v-4",},
				signature: "v5",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Signature: all params",
			args: args{
				header:        `keyId="v1",algorithm="v2",created=1402170695,expires=1402170699,headers="v-3 v-4",signature="v5"`,
				authorization: false,
			},
			want: ParsedHeader{
				keyId:     "v1",
				algorithm: "v2",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"v-3", "v-4",},
				signature: "v5",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Signature: real example",
			args: args{
				header:        `keyId="Test",algorithm="rsa-sha256",created=1402170695,expires=1402170699,headers="(request-target) (created) (expires) host date content-type digest content-length",signature="vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE="`,
				authorization: false,
			},
			want: ParsedHeader{
				keyId:     "Test",
				algorithm: "rsa-sha256",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"(request-target)", "(created)", "(expires)", "host", "date", "content-type", "digest", "content-length",},
				signature: "vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE=",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name: "Unsupported symbol in key",
			args: args{
				header: `keyId-="v1"`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "found '-' — unsupported symbol in key",
		},
		{
			name: "Unsupported symbol, expected = symbol",
			args: args{
				header: `keyId :"v1"`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "found ':' — unsupported symbol, expected '=' or space symbol",
		},
		{
			name: "Unsupported symbol, expected quote symbol",
			args: args{
				header: `keyId= 'v1'`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "found ''' — unsupported symbol, expected '\"' or space symbol",
		},
		{
			name: "Unknown key",
			args: args{
				header: `key="v1"`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unknown key: 'key'",
		},
		{
			name: "unexpected end of header, expected equal symbol",
			args: args{
				header: `keyId`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected end of header, expected '=' symbol and field value",
		},
		{
			name: "Expected field value",
			args: args{
				header: `keyId `,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected end of header, expected field value",
		},
		{
			name: "Expected quote",
			args: args{
				header: `keyId= `,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected end of header, expected '\"' symbol and field value",
		},
		{
			name: "Expected quote at the end",
			args: args{
				header: `keyId="`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected end of header, expected '\"' symbol",
		},
		{
			name: "Empty value",
			args: args{
				header: `keyId=""`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "empty value for key 'keyId'",
		},
		{
			name: "Div symbol expected",
			args: args{
				header: `keyId="v1" algorithm="v2"`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "found 'a' — unsupported symbol, expected ',' or space symbol",
		},
		{
			name: "Wrong created INT value",
			args: args{
				header: `created=9223372036854775807`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "wrong 'created' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
		{
			name: "Wrong expires INT value",
			args: args{
				header: `expires=9223372036854775807`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
		{
			name: "Wrong expires with space at the end",
			args: args{
				header: `expires=9223372036854775807 `,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
		{
			name: "Wrong expires with divider",
			args: args{
				header: `expires=9223372036854775807,`,
			},
			want:       ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "wrong 'expires' param value: strconv.ParseInt: parsing \"9223372036854775807\": value out of range",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Create()
			if tt.args.authorization == true {
				p.keywordNow = true
			} else {
				p.keyNow = true
			}
			got, err := p.parse(tt.args.header)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("error message = `%s`, wantErrMsg = `%s`", err.Error(), tt.wantErrMsg)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = `%v`, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v,\nwant %v", got, tt.want)
			}
		})
	}
}

func TestParserParseAuthorization(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name       string
		args       args
		want       ParsedHeader
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Authorization",
			args: args{
				header: `Signature keyId="Test",algorithm="rsa-sha256",created=1402170695,expires=1402170699,headers="(request-target) (created) (expires) host date content-type digest content-length",signature="vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE="`,
			},
			want: ParsedHeader{
				keyword:   "Signature",
				keyId:     "Test",
				algorithm: "rsa-sha256",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"(request-target)", "(created)", "(expires)", "host", "date", "content-type", "digest", "content-length",},
				signature: "vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE=",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Create()
			got, err := p.ParseAuthorization(tt.args.header)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("error message = `%s`, wantErrMsg = `%s`", err.Error(), tt.wantErrMsg)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = `%v`, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v,\nwant %v", got, tt.want)
			}
		})
	}
}

func TestParserParseSignature(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name       string
		args       args
		want       ParsedHeader
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Signature",
			args: args{
				header: `keyId="Test", algorithm="rsa-sha256", created=1402170695, expires=1402170699, headers="(request-target) (created) (expires)", signature="vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE="`,
			},
			want: ParsedHeader{
				keyId:     "Test",
				algorithm: "rsa-sha256",
				created:   1402170695,
				expires:   1402170699,
				headers:   []string{"(request-target)", "(created)", "(expires)",},
				signature: "vSdrb+dS3EceC9bcwHSo4MlyKS59iFIrhgYkz8+oVLEEzmYZZvRs8rgOp+63LEM3v+MFHB32NfpB2bEKBIvB1q52LaEUHFv120V01IL+TAD48XaERZFukWgHoBTLMhYS2Gb51gWxpeIq8knRmPnYePbF5MOkR0Zkly4zKH7s1dE=",
			},
			wantErr:    false,
			wantErrMsg: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Create()
			got, err := p.ParseSignature(tt.args.header)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("error message = `%s`, wantErrMsg = `%s`", err.Error(), tt.wantErrMsg)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = `%v`, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v,\nwant %v", got, tt.want)
			}
		})
	}
}

func TestParserParseFailed(t *testing.T) {
	type args struct {
		header string
	}
	tests := []struct {
		name       string
		args       args
		want       ParsedHeader
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Current parser stage not set",
			args: args{
				header: `keyId="Test"`,
			},
			want: ParsedHeader{},
			wantErr:    true,
			wantErrMsg: "unexpected parser stage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Create()
			got, err := p.parse(tt.args.header)
			if err != nil && err.Error() != tt.wantErrMsg {
				t.Errorf("error message = `%s`, wantErrMsg = `%s`", err.Error(), tt.wantErrMsg)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = `%v`, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parse() got = %v,\nwant %v", got, tt.want)
			}
		})
	}
}