package panel

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/InazumaV/V2bX/conf"
	"github.com/goccy/go-json"
)

func TestReportNodeOnlineUsersSendsXboardOnlineCounts(t *testing.T) {
	var gotPath string
	var gotPayload map[string]json.RawMessage

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request body failed: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := New(&conf.ApiConfig{
		APIHost:  server.URL,
		Key:      "token",
		NodeType: "vless",
		NodeID:   1,
	})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}

	payload := map[int][]string{
		1: {"1.1.1.1", "2.2.2.2"},
		2: {"3.3.3.3"},
	}
	if err := client.ReportNodeOnlineUsers(&payload); err != nil {
		t.Fatalf("ReportNodeOnlineUsers failed: %v", err)
	}

	if gotPath != "/api/v2/server/report" {
		t.Fatalf("path = %s, want /api/v2/server/report", gotPath)
	}

	var online map[string]int
	if err := json.Unmarshal(gotPayload["online"], &online); err != nil {
		t.Fatalf("decode online payload failed: %v", err)
	}
	if online["1"] != 2 || online["2"] != 1 {
		t.Fatalf("online payload = %#v, want counts for alive IPs", online)
	}
}
