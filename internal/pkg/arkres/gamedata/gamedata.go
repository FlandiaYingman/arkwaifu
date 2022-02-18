package gamedata

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	githubArknightsGameDataOwner = "Kengxxiao"
	githubArknightsGameDataRepo  = "ArknightsGameData"
)

// Get downloads the Arknights gamedata from https://github.com/Kengxxiao/ArknightsGameData.
//
// Firstly, it downloads the repository archive from GitHub.
// Then, it extracts the repository archive (only files under "zh_CN/gamedata/") to the specified path.
//
// For example, saying that path is "story/" and dest is "./arkwaifu/".
// What Get does is to download the full gamedata archive and extract files in the archive under "zh_CN/gamedata/story/" into "./arkwaifu/story".
func Get(resVersion string, dataPath string, dest string) error {
	commitRef, err := findCommitByResVersion(githubArknightsGameDataOwner, githubArknightsGameDataRepo, resVersion)
	if err != nil {
		return err
	}
	link, err := getZipballLink(githubArknightsGameDataOwner, githubArknightsGameDataRepo, commitRef)
	if err != nil {
		return err
	}

	zipball, err := downloadZipball(link, dest)
	if err != nil {
		return err
	}
	defer os.Remove(zipball)

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
