package workspace

import (
	"encoding/json"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"io/ioutil"
)

type FileFetcher struct {
	filename string
}

func NewFileFetcher(filename string) *FileFetcher {
	return &FileFetcher{filename: filename}
}

func (f *FileFetcher) Fetch() (Workspace, bool) {
	bytes, err := ioutil.ReadFile(f.filename)
	if err != nil {
		logger.Warn("failed to read file: %v", err)
		return nil, false
	}
	var dto WorkspaceDTO
	err = json.Unmarshal(bytes, &dto)
	if err != nil {
		logger.Warn("failed to unmarshal workspace: %v", err)
		return nil, false
	}
	return NewFrom(dto), true
}

func (f *FileFetcher) Close() {}
