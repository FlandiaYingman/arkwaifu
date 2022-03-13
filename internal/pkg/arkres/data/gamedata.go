package data

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	githubArknightsGameDataOwner = "Kengxxiao"
	githubArknightsGameDataRepo  = "ArknightsGameData"
)

// FindCommitRef checks whether specified resource version exists or not.
// If exists, it returns the corresponded commit reference, otherwise, it returns an empty string.
func FindCommitRef(resVersion string) (string, error) {
	return findCommitByResVersion(githubArknightsGameDataOwner, githubArknightsGameDataRepo, resVersion)
}

// FindLatestCommitRef finds the latest gamedata update commit of CN server.
func FindLatestCommitRef() (resVersion string, commitRef string, err error) {
	message, commitRef, err := findLatestCommit(githubArknightsGameDataOwner, githubArknightsGameDataRepo)
	if err != nil {
		return "", "", err
	}
	_, _, resVersion = parseCommitMessage(message)
	return resVersion, commitRef, nil
}

// parseCommitMessage parses the commit message of ArknightsGameData repo.
func parseCommitMessage(message string) (server string, clientVersion string, resVersion string) {
	r := regexp.MustCompile(`^\[(.*?) UPDATE] Client:(.*?) Data:(.*?)$`)
	result := r.FindStringSubmatch(message)
	if result == nil {
		return "", "", ""
	}
	return result[1], result[2], result[3]
}

// Get downloads the Arknights gamedata from https://github.com/Kengxxiao/ArknightsGameData.
//
// Firstly, it downloads the repository archive from GitHub.
// Then, it extracts the repository archive (only files under "zh_CN/gamedata/") to the specified path.
//
// For example, saying that path is "story/" and dest is "./arkwaifu/".
// What Get does is to download the full gamedata archive and extract files in the archive under "zh_CN/gamedata/story/" into "./arkwaifu/story".
func Get(resVersion string, dataPath string, dest string) error {
	commitRef, err := FindCommitRef(resVersion)
	if err != nil {
		return err
	}
	if commitRef == "" {
		return errors.Errorf("commit by res version %v not found", resVersion)
	}
	link, err := getZipballLink(githubArknightsGameDataOwner, githubArknightsGameDataRepo, commitRef)
	if err != nil {
		return err
	}

	zipball, err := downloadZipball(link, dest)
	if err != nil {
		return err
	}
	defer func(name string) {
		_ = os.Remove(name)
	}(zipball)

	basename := filepath.Base(zipball)
	purename := strings.TrimSuffix(basename, filepath.Ext(basename))
	zipballRootPath := path.Join(purename, "zh_CN", "gamedata")
	zipballDataPath := path.Join(zipballRootPath, dataPath)
	err = extractZipball(zipball, zipballRootPath, []string{zipballDataPath}, dest)
	if err != nil {
		return err
	}

	return nil
}

func GetText(gamedata string, dataPath string) (string, error) {
	fullpath := filepath.Join(gamedata, dataPath)
	file, err := ioutil.ReadFile(fullpath)
	if err != nil {
		return "", err
	}
	return string(file), err
}
