package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func LoadSecrets(path string, secrets map[string]*string) error {
	if path == "" {
		return fmt.Errorf("secrets path must be provided")
	}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("provided secrets path does not exist: %s", path)
		}
		return err
	}
	f.Close()

	for name, variable := range secrets {
		secretPath := filepath.Join(path, name)
		content, err := os.ReadFile(secretPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("provided secret does not exist in %s: %s", path, name)
			}
		}
		*variable = string(content)
	}
	return nil
}
