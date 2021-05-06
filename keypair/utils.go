package keypair

import (
	bs "encoding/base64"
	"errors"
	"github.com/robpike/filter"
	"strings"
	"os"
)

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
	base64 := strings.Split(content, "/\r\n/")
	var selec = filter.Choose(base64, filterFunction).([]string)
	join := strings.Join(selec, "")
	res := strings.Trim(join, "\r\n")
	return bs.StdEncoding.DecodeString(res)
}

func filterFunction(a string) bool {
	return ! strings.HasPrefix(a, "-----")
}
