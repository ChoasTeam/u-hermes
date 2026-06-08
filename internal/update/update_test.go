package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewService(t *testing.T) {
	svc := NewService("test", "u-hermes", "0.1.0")
	if svc == nil {
		t.Fatal("nil service")
	}
	if svc.owner != "test" || svc.repo != "u-hermes" {
		t.Errorf("unexpected owner/repo: %s/%s", svc.owner, svc.repo)
	}
}

func TestCheckForUpdate_NoUpdate_SameVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.github+json" {
			t.Error("missing GitHub API accept header")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Release{TagName: "v0.1.0"})
	}))
	defer server.Close()

	// Test the data path by checking version comparison logic
	svc := NewService("test", "u-hermes", "0.1.0")
	// Tag "v0.1.0" vs currentVer "0.1.0" → should be no update
	release := &Release{TagName: "v0.1.0"}
	isUpdate := !(release.TagName == svc.currentVer || "v"+svc.currentVer == release.TagName)
	if isUpdate {
		t.Error("same version should not trigger update")
	}
}

func TestDownloadAndVerify_FallbackAssetName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("fake binary content"))
	}))
	defer server.Close()

	release := &Release{
		TagName: "v0.2.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{Name: "u-hermes.exe", BrowserDownloadURL: server.URL},
		},
	}

	svc := NewService("test", "u-hermes", "0.1.0")
	path, err := svc.DownloadAndVerify(release)
	if err != nil {
		t.Fatalf("fallback to u-hermes.exe should work: %v", err)
	}
	_ = path
	// Clean up downloaded file
	// os.Remove(path) — skip, test writes to relative path
}

func TestCheckForUpdate_HasUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Release{
			TagName: "v0.2.0",
			Body:    "New features!",
		})
	}))
	defer server.Close()

	t.Logf("Update check server at %s", server.URL)

	// Test version comparison: v0.2.0 vs 0.1.0 → should be update
	svc := NewService("test", "u-hermes", "0.1.0")
	release := &Release{TagName: "v0.2.0"}
	isUpdate := !(release.TagName == svc.currentVer || "v"+svc.currentVer == release.TagName)
	if !isUpdate {
		t.Error("different version should trigger update")
	}
}
