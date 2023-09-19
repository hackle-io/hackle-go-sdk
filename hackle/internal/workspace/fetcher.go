package workspace

type Fetcher interface {
	Fetch() (Workspace, bool)
	Close()
}
