package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func NormalizeSafeFileName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("filename is required")
	}
	if strings.ContainsRune(name, 0) {
		return "", fmt.Errorf("filename contains invalid null byte")
	}
	if strings.ContainsAny(name, `/\\`) {
		return "", fmt.Errorf("filename must not contain path separators")
	}

	cleaned := filepath.Clean(name)
	if cleaned == "." || cleaned == ".." || cleaned == "" {
		return "", fmt.Errorf("invalid filename")
	}
	if filepath.Base(cleaned) != cleaned {
		return "", fmt.Errorf("filename must be a single path segment")
	}
	return cleaned, nil
}

func NormalizeSafeSubDirPath(path string) (string, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return "", fmt.Errorf("path is required")
	}

	cleaned, err := cleanRelativePath(trimmed)
	if err != nil {
		return "", err
	}
	if cleaned == "." {
		return "", fmt.Errorf("path must not be empty")
	}

	return filepath.ToSlash(cleaned) + "/", nil
}

func ResolvePathWithinBase(baseDir, userPath string) (string, error) {
	basePath, err := canonicalBasePath(baseDir)
	if err != nil {
		return "", err
	}

	relPath, err := cleanRelativePath(userPath)
	if err != nil {
		return "", err
	}

	candidatePath := filepath.Join(basePath, relPath)
	if !pathWithinBase(basePath, candidatePath) {
		return "", fmt.Errorf("path escapes base directory")
	}

	if relPath != "." {
		parentPath := filepath.Dir(candidatePath)
		if resolvedParent, err := filepath.EvalSymlinks(parentPath); err == nil {
			if !pathWithinBase(basePath, resolvedParent) {
				return "", fmt.Errorf("path escapes base directory via symlink")
			}
			candidatePath = filepath.Join(resolvedParent, filepath.Base(candidatePath))
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}

	if info, err := os.Lstat(candidatePath); err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			resolvedTarget, err := filepath.EvalSymlinks(candidatePath)
			if err != nil {
				return "", err
			}
			if !pathWithinBase(basePath, resolvedTarget) {
				return "", fmt.Errorf("path escapes base directory via symlink")
			}
			candidatePath = resolvedTarget
		}
	} else if !os.IsNotExist(err) {
		return "", err
	}

	return candidatePath, nil
}

func canonicalBasePath(baseDir string) (string, error) {
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}

	resolvedBase, err := filepath.EvalSymlinks(absBase)
	if err != nil {
		if os.IsNotExist(err) {
			return filepath.Clean(absBase), nil
		}
		return "", err
	}

	return filepath.Clean(resolvedBase), nil
}

func cleanRelativePath(userPath string) (string, error) {
	trimmed := strings.TrimSpace(userPath)
	if trimmed == "" || trimmed == "." {
		return ".", nil
	}
	if strings.ContainsRune(trimmed, 0) {
		return "", fmt.Errorf("path contains invalid null byte")
	}
	if strings.Contains(trimmed, "\\") {
		return "", fmt.Errorf("path contains invalid separator")
	}
	if filepath.IsAbs(trimmed) {
		return "", fmt.Errorf("absolute path is not allowed")
	}

	cleaned := filepath.Clean(trimmed)
	if cleaned == "." {
		return ".", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path traversal detected")
	}

	return cleaned, nil
}

func pathWithinBase(basePath, targetPath string) bool {
	rel, err := filepath.Rel(basePath, filepath.Clean(targetPath))
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)))
}
