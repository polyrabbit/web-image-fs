package webpage

import (
	"testing"
)

func TestHTTPClient_URLJoin(t *testing.T) {
	type args struct {
		base         string
		relativePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{{
		name: "abs path",
		args: args{
			base:         "http://pic.baidu.com",
			relativePath: "/bbb",
		},
		want: "http://pic.baidu.com/bbb",
	}, {
		name: "relative path",
		args: args{
			base:         "http://pic.baidu.com",
			relativePath: "bbb",
		},
		want: "http://pic.baidu.com/bbb",
	}, {
		name: "nested base",
		args: args{
			base:         "http://pic.baidu.com/abc/def",
			relativePath: "bbb",
		},
		want: "http://pic.baidu.com/abc/bbb",
	}, {
		name: "new path",
		args: args{
			base:         "http://pic.baidu.com/abc/def",
			relativePath: "http://img.betacat.io/sit.png",
		},
		want: "http://img.betacat.io/sit.png",
	}, {
		name: "abs path no schema",
		args: args{
			base:         "http://pic.baidu.com/abc/def",
			relativePath: "//img.betacat.io/sit.png",
		},
		want: "http://img.betacat.io/sit.png",
	}, {
		name: "new path",
		args: args{
			base:         "http://pic.baidu.com/abc/def",
			relativePath: "javascript:void(0)",
		},
		want: "javascript:void(0)",
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &HTTPClient{}
			got, err := c.URLJoin(tt.args.base, tt.args.relativePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("URLJoin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("URLJoin() got = %v, want %v", got, tt.want)
			}
		})
	}
}
