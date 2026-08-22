package testdata

import "net/http"

type Notes struct{}

type Memory struct{}

func (n *Notes) Search(w http.ResponseWriter, r *http.Request) {}

func (m *Memory) SearchCursor(query, cursor string, limit int) (items []string, next string) {
	return nil, ""
}
