package minifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func JSON(data []byte) ([]byte, error) {
	buff := new(bytes.Buffer)
	errCompact := json.Compact(buff, data)
	if errCompact != nil {
		newErr := fmt.Errorf("failure encountered compacting json := %v", errCompact)
		return nil, newErr
	}

	b, err := ioutil.ReadAll(buff)
	if err != nil {
		readErr := fmt.Errorf("read buffer error encountered := %v", err)
		return nil, readErr
	}

	return b, nil
}
