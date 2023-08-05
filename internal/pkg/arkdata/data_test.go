package arkdata

import (
	"context"
	"github.com/flandiayingman/arkwaifu/internal/pkg/ark"
	"github.com/flandiayingman/arkwaifu/internal/pkg/util/pathutil"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestGetGameData(t *testing.T) {
	type args struct {
		ctx      context.Context
		server   ark.Server
		version  ark.Version
		patterns []string
		dst      string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "specific version story_review_table",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "23-07-27-18-50-06-aeb568",
				patterns: []string{"**/story_review_table.json"},
				dst:      "./test_1",
			},
		},
		{
			name: "latest all .json",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "",
				patterns: []string{"**/*.json"},
				dst:      "./test_2",
			},
		},
		{
			name: "latest all stories",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "",
				patterns: []string{"gamedata/story/**/*.txt"},
				dst:      "./test_3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.args.dst)
			err := GetGameData(tt.args.ctx, tt.args.server, tt.args.version, tt.args.patterns, tt.args.dst)
			if err != nil {
				t.Error(err)
			}

			files, err := pathutil.Glob(tt.args.patterns, path.Join(tt.args.dst, DefaultPrefix))
			if err != nil {
				t.Error(err)
			}

			assert.NotEmpty(t, files)
		})
	}
}

func TestGetCompositeGameData(t *testing.T) {
	type args struct {
		ctx      context.Context
		server   ark.Server
		version  ark.Version
		patterns []string
		dst      string
	}
	tests := []struct {
		name       string
		args       args
		successful bool
	}{
		{
			name: "specific version story_review_table",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "23-07-27-18-50-06-aeb568",
				patterns: []string{"**/story_review_table.json"},
				dst:      "./test_1",
			},
			successful: true,
		},
		{
			name: "latest all .json",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "",
				patterns: []string{"**/*.json"},
				dst:      "./test_2",
			},
			successful: true,
		},
		{
			name: "latest all stories",
			args: args{
				ctx:      context.Background(),
				server:   ark.CnServer,
				version:  "",
				patterns: []string{"gamedata/story/**/*.txt"},
				dst:      "./test_3",
			},
			successful: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.RemoveAll(tt.args.dst)
			err := GetCompositeGameData(tt.args.ctx, tt.args.server, tt.args.version, tt.args.patterns, tt.args.dst)

			files, err := pathutil.Glob(tt.args.patterns, path.Join(tt.args.dst, DefaultPrefix))
			if err != nil {
				t.Error(err)
			}

			if tt.successful {
				assert.NotEmpty(t, files)
			} else {
				assert.Empty(t, files)
			}
		})
	}
}
