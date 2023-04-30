package test

import (
	"io/ioutil"
	"testing"
)

// FileToBytes reads fileName and returns the file contents as a byte array
func FileToBytes(t *testing.T, fileName string) (bytes []byte) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		t.Errorf("test.FileToString: failed reading %s", fileName)
	}

	return
}
