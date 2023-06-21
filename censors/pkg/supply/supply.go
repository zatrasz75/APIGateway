package supply

import (
	"censorship/pkg/storage"
	"io/ioutil"
	"os"
	"strings"
)

func StopList() ([]storage.Stop, error) {
	f, err := os.Open("./pkg/supply/words.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	var sl []storage.Stop
	for _, line := range lines {
		str := storage.Stop{
			StopList: line,
		}
		sl = append(sl, str)
	}

	return sl, nil
}
