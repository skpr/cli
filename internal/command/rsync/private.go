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
	// privateScopeRequiresFlag means the source is, or may expand to, the private filesystem
	// or a path within it, so downloading it requires the private files flag.
	privateScopeRequiresFlag
	// privateScopeParent means the source is a parent directory of the private filesystem,
	// which can be synced with the private filesystem excluded.
	privateScopeParent
)

// privateFilesElements are the path elements of PrivateFilesPath.
var privateFilesElements = strings.Split(strings.TrimPrefix(PrivateFilesPath, "/"), "/")

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

	// A glob such as "/mnt/*" is expanded on the remote side and may include the private
	// filesystem. An exclude can't be anchored reliably for an expanded list, so require
	// the flag whenever the glob could reach it.
	if strings.ContainsAny(remotePath, "*?[") {
		if globMayReachPrivate(remotePath) {
			return privateScopeRequiresFlag
		}

		return privateScopeNone
	}

	cleaned := path.Clean(remotePath)

	switch {
	case isWithin(cleaned, PrivateFilesPath):
		return privateScopeRequiresFlag
	case isWithin(PrivateFilesPath, cleaned):
		// rsync picks the transfer root from the last path element, and "." or ".." make
		// that root hard to predict. Without a known root an exclude can't be anchored
		// safely, so require the private files flag instead.
		if trimmed := strings.TrimRight(remotePath, "/"); trimmed != "" {
			if base := path.Base(trimmed); base == "." || base == ".." {
				return privateScopeRequiresFlag
			}
		}

		return privateScopeParent
	default:
		return privateScopeNone
	}
}

// globMayReachPrivate reports whether an absolute glob could expand to the private filesystem,
// a path within it, or one of its parents. Shell globs don't match across "/", so each element
// is matched against the element of PrivateFilesPath at the same depth. Globs containing "." or
// ".." elements, or invalid patterns, are assumed to reach it.
func globMayReachPrivate(glob string) bool {
	var elements []string
	for _, e := range strings.Split(glob, "/") {
		switch e {
		case "":
			continue
		case ".", "..":
			return true
		}
		elements = append(elements, e)
	}

	for i := 0; i < len(elements) && i < len(privateFilesElements); i++ {
		matched, err := path.Match(elements[i], privateFilesElements[i])
		if err != nil {
			return true
		}
		if !matched {
			return false
		}
	}

	return true
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
