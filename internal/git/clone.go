package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kroyser123/go-mod-updater/internal/logger"
)

// клонируем репозиторий
func Clone(repoURL string, token string, logger *logger.Logger) (string, error) {
	logger.Info("Cloning repository %s", repoURL)
	logger.Debug("token: %v", token)

	// создаем папку со случайным именем: modupdater-jwu8293ui(к примеру)

	tmpDir, err := os.MkdirTemp("", "modupdater-*")
	if err != nil {
		logger.Error("Failed to create a directory: %v", err)
		return "", err
	}
	clone := repoURL
	if token != "" && strings.HasPrefix(repoURL, "https://") {
		clone = strings.Replace(repoURL, "https://", "https://token:"+token+"@", 1)
	}

	logger.Debug("git clone %s to %s", clone, tmpDir)

	// клонируем только последний коммит

	cmd := exec.Command("git", "clone", "--depth", "1", clone, tmpDir)

	// собираем весь вывод: и успех, и ошибки

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Clone failed: %v", err)
		os.RemoveAll(tmpDir)
		return "", fmt.Errorf("git clone failed: %s: %v", string(output), err)
	}
	return tmpDir, nil
}
