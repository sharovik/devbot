package helper

import (
	"encoding/json"
	"io/ioutil"
)

//FileToBytes method gets the data from selected file and retrieve the byte value of it
func FileToBytes(filePath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

//ObjectToFile method save object in selected path
func ObjectToFile(path string, object interface{}) error {
	if err := ioutil.WriteFile(path, convertToByte(object), 0666); err != nil {
		return err
	}

	return nil
}

func convertToByte(object interface{}) []byte {
	str, _ := json.Marshal(object)
	return str
}
