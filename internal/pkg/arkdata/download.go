package arkdata

import (
	"context"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
	"os"
)

func download(ctx context.Context, repoOwner, repoName, sha string) (string, error) {
	temp, err := os.MkdirTemp("", "arkdata_download")
	if err != nil {
		return "", err
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	request, err := grab.NewRequest(temp, urlOfZipball(repoOwner, repoName, sha))
	if err != nil {
		return "", err
	}

	request = request.WithContext(ctx)
	response := grab.DefaultClient.Do(request)
	return response.Filename, response.Err()
}

func urlOfZipball(repoOwner, repoName, sha string) string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", repoOwner, repoName, sha)
}
