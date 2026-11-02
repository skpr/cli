package rsync

import (
	"path"
	"strings"
)

// PrivateFilesPath is where the private filesystem is mounted in an environment.
const PrivateFilesPath = "/mnt/private"

// privateScope describes how a remote rsync source relates to the private filesystem.
type privateScope int

const (
	// privateScopeNone means the source does not include private files.
	privateScopeNone privateScope = iota
	// privateScopeInside means the source is the private filesystem or a path within it.
	privateScopeInside
	// privateScopeParent means the source is a parent directory of the private filesystem.
	privateScopeParent
)

// splitRemote splits an rsync remote location ("env:path") into its environment and path.
// The boolean is false when the location is a local path.
func splitRemote(location string) (string, string, bool) {
	env, remotePath, found := strings.Cut(location, ":")
	if !found || env == "" || strings.Contains(env, "/") {
		return "", "", false
	}

	return env, remotePath, true
}

// classifyRemotePath reports how a remote source path relates to the private filesystem.
// Relative paths are resolved by the remote shell and cannot be checked here, so they are
// treated as not including private files.
func classifyRemotePath(remotePath string) privateScope {
	if !path.IsAbs(remotePath) {
		return privateScopeNone
	}

	// A glob such as "/mnt/*" can expand to the private filesystem on the remote side.
	// Only its literal prefix is known here, so treat any glob that could reach the
	// private filesystem as being inside it.
	if i := strings.IndexAny(remotePath, "*?["); i >= 0 {
		prefix := path.Dir(remotePath[:i+1])
		if isWithin(prefix, PrivateFilesPath) || isWithin(PrivateFilesPath, prefix) {
			return privateScopeInside
		}

		return privateScopeNone
	}

	cleaned := path.Clean(remotePath)

	switch {
	case isWithin(cleaned, PrivateFilesPath):
		return privateScopeInside
	case isWithin(PrivateFilesPath, cleaned):
		// rsync picks the transfer root from the last path element, and "." or ".." make
		// that root hard to predict. Without a known root an exclude can't be anchored
		// safely, so require the private files flag instead.
		if trimmed := strings.TrimRight(remotePath, "/"); trimmed != "" {
			if base := path.Base(trimmed); base == "." || base == ".." {
				return privateScopeInside
			}
		}

		return privateScopeParent
	default:
		return privateScopeNone
	}
}

// isWithin reports whether child is the same path as parent or is below it.
func isWithin(child, parent string) bool {
	if child == parent || parent == "/" {
		return true
	}

	return strings.HasPrefix(child, parent+"/")
}

// privateExcludePattern returns an anchored rsync exclude pattern that matches the private
// filesystem when syncing from remotePath, which must be a parent of the private filesystem.
//
// rsync anchors patterns to the root of the transfer. With a trailing slash ("/mnt/") the
// root is the directory itself, otherwise ("/mnt") it is the directory's parent.
func privateExcludePattern(remotePath string) string {
	root := path.Clean(remotePath)
	if !strings.HasSuffix(remotePath, "/") {
		root = path.Dir(root)
	}

	rel := strings.TrimPrefix(PrivateFilesPath, root)
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}

	return rel + "/"
}
