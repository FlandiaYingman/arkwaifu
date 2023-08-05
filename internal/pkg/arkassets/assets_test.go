package arkassets

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestGetGameAssets(t *testing.T) {
	type args struct {
		ctx      context.Context
		version  ark.Version
		dst      string
		patterns []string
	}
	tests := []struct {
		name        string
		args        args
		wantSuccess bool
		wantErr     bool
	}{
		{
			name: "avg/imgs/avg_img_0_0.ab",
			args: args{
				ctx:      context.Background(),
				version:  "",
				dst:      "./test_dst",
				patterns: []string{"avg/imgs/avg_img_0_0.ab"},
			},
			wantSuccess: true,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.args.dst)
			if err := GetGameAssets(tt.args.ctx, tt.args.version, tt.args.dst, tt.args.patterns); (err != nil) != tt.wantErr {
				t.Errorf("GetGameAssets() error = %v, wantErr %v", err, tt.wantErr)
			}
			files, err := pathutil.Glob([]string{"**"}, path.Join(tt.args.dst))
			if err != nil {
				t.Error(err)
			}
			if tt.wantSuccess {
				assert.NotEmpty(t, files)
			} else {
				assert.Empty(t, files)
			}
		})
	}
}
