package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common"
	"github.com/Jacobbrewer1/f1-data/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestDecodeJSONBody(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type testCase struct {
		name    string
		args    args
		want    any
		wantErr error
	}
	tests := []testCase{
		{
			name: "success",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"message": "hello"}`)),
			},
			want: &common.Message{Message: utils.Ptr("hello")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(common.Message)
			err := DecodeJSONBody(tt.args.r, got)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}
