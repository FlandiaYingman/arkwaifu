package agdapi

import (
	"context"
	"fmt"
	"github.com/cavaliergopher/grab/v3"
)

// fetchZipball fetches the zipball of the repo versioned by resourceVersion to dstDir, and returns the path of the zipball downloaded.
func fetchZipball(ctx context.Context, resourceVersion ResourceVersion, dstDir string) (string, error) {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	request, err := grab.NewRequest(dstDir, getZipballURL(resourceVersion))
	if err != nil {
		return "", err
	}

	request = request.WithContext(ctx)
	response := grab.DefaultClient.Do(request)
	return response.Filename, response.Err()
}

func getZipballURL(resourceVersion ResourceVersion) string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/%s.zip", RepoOwner, RepoName, resourceVersion.CommitSHA)
}
