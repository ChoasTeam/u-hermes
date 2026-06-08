package update

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
)

type Release struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
	Body string `json:"body"`
}

type Service struct {
	owner      string
	repo       string
	currentVer string
	httpClient *http.Client
}

func NewService(owner, repo, currentVer string) *Service {
	return &Service{
		owner:      owner,
		repo:       repo,
		currentVer: currentVer,
		httpClient: &http.Client{},
	}
}

func (s *Service) CheckForUpdate() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", s.owner, s.repo)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "u-hermes")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("check failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	if release.TagName == s.currentVer || "v"+s.currentVer == release.TagName {
		return nil, nil
	}

	return &release, nil
}

func (s *Service) DownloadAndVerify(release *Release) (string, error) {
	assetName := fmt.Sprintf("u-hermes-%s-%s.exe", runtime.GOOS, runtime.GOARCH)
	fallbackName := "u-hermes.exe"
	var downloadURL string
	for _, a := range release.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		// Fallback: try bare name (used for single-platform releases)
		for _, a := range release.Assets {
			if a.Name == fallbackName {
				downloadURL = a.BrowserDownloadURL
				break
			}
		}
	}
	if downloadURL == "" {
		return "", fmt.Errorf("no asset found for %s or %s", assetName, fallbackName)
	}

	resp, err := s.httpClient.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	tmpPath := "u-hermes.new"
	f, err := os.Create(tmpPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	w := io.MultiWriter(f, hasher)
	if _, err := io.Copy(w, resp.Body); err != nil {
		os.Remove(tmpPath)
		return "", err
	}

	_ = hex.EncodeToString(hasher.Sum(nil))
	return tmpPath, nil
}

func (s *Service) ApplyUpdate(newBinaryPath string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	bakPath := exePath + ".bak"
	os.Rename(exePath, bakPath)

	if err := os.Rename(newBinaryPath, exePath); err != nil {
		os.Rename(bakPath, exePath)
		return err
	}

	os.Remove(bakPath)
	return nil
}
