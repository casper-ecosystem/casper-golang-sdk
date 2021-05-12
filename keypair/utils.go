package keypair

import (
	bs "encoding/base64"
	"errors"
	"github.com/robpike/filter"
	"os"
	"strings"
)

//ReadBase64File reads the Base64 content of a file, get rid of PEM frames
func ReadBase64File(path string) ([]byte,error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("can't read file")
	}
	resContent := bytesToString(content)
	return ReadBase64WithPEM(resContent)
}

func bytesToString(data []byte) string {
	return string(data[:])
}

func ReadBase64WithPEM(content string) ([]byte, error) {
	base64 := strings.Split(content, "\n")
	var selec = filter.Choose(base64, filterFunction).([]string)
	join := strings.Join(selec, "")
	res := strings.Trim(join, "\r\n")
	return bs.StdEncoding.DecodeString(res)
}

func filterFunction(a string) bool {
	return ! strings.HasPrefix(a, "-----")
}
