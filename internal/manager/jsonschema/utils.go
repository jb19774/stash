package jsonschema

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"
)

var nilTime = (time.Time{}).UnixNano()

func CompareJSON(a interface{}, b interface{}) bool {
	aBuf, _ := encode(a)
	bBuf, _ := encode(b)
	return bytes.Compare(aBuf, bBuf) == 0
}

func marshalToFile(filePath string, j interface{}) error {
	data, err := encode(j)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filePath, data, 0755); err != nil {
		return err
	}
	return nil
}

func encode(j interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(j); err != nil {
		return nil, err
	}
	// Strip the newline at the end of the file
	return bytes.TrimRight(buffer.Bytes(), "\n"), nil
}