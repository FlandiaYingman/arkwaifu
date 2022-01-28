package arkres

import (
	"testing"
)

func TestGetGameData(t *testing.T) {
	txt, err := GetStoryTxt("22-01-14-17-58-34-bd68ad", "gamedata\\story\\activities\\a001\\level_a001_01_end.txt")
	if err != nil {
		panic(err)
	}
	println(txt)
}
