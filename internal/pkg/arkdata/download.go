package arkdata

import (
	"context"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"github.com/pkg/errors"
	"os"
)

func download(ctx context.Context, repoOwner, repoName, sha string) (string, string, error) {
	temp, err := os.MkdirTemp("", "arkdata_download-*")
	if err != nil {
		return "", "", errors.WithStack(err)
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	request, err := grab.NewRequest(temp, urlOfZipball(repoOwner, repoName, sha))
	if err != nil {
		os.RemoveAll(temp)
		return "", "", errors.WithStack(err)
	}

	request = request.WithContext(ctx)
	client := grab.NewClient()
	client.UserAgent = "FlandiaYingman/arkwaifu"
	response := client.Do(request)
	return temp, response.Filename, response.Err()
}

func urlOfZipball(repoOwner, repoName, sha string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/zipball/%s", repoOwner, repoName, sha)
}
