package container

import (
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/sharovik/devbot/internal/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	//We switch pointer to the root directory for control the path from which we need to generate test-data file-paths
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../../")
	_ = os.Chdir(dir)
}

func TestMain_Init(t *testing.T) {
	c := C.Init()
	assert.IsType(t, Main{}, c)

	assert.Equal(t, true, c.Config.IsInitialised())
	assert.Equal(t, true, log.IsInitialized())
}
