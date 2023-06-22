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
		trimmedLine := strings.TrimSpace(line)
		str := storage.Stop{
			StopList: trimmedLine,
		}
		sl = append(sl, str)
	}

	return sl, nil
}
