package rsync

import (
	"path"
	"strings"
)

// PrivateFilesPath is where the private filesystem is mounted.
const PrivateFilesPath = "/mnt/private"

// privateScope describes how a remote rsync source relates to the private filesystem.
type privateScope int

const (
	privateScopeNone privateScope = iota
	// privateScopeRequiresFlag means the source is, or may expand to, private files.
	privateScopeRequiresFlag
	// privateScopeParent means the source contains the private filesystem, which can be excluded.
	privateScopeParent
)

// privateFilesElements are the path elements of PrivateFilesPath.
var privateFilesElements = strings.Split(strings.TrimPrefix(PrivateFilesPath, "/"), "/")

// splitRemote splits "env:path" into its parts. The boolean is false for local paths.
func splitRemote(location string) (string, string, bool) {
	env, remotePath, found := strings.Cut(location, ":")
	if !found || env == "" || strings.Contains(env, "/") {
		return "", "", false
	}

	return env, remotePath, true
}

// classifyRemotePath reports how a remote path relates to the private filesystem.
// Relative paths can't be resolved locally, so they aren't checked.
func classifyRemotePath(remotePath string) privateScope {
	if !path.IsAbs(remotePath) {
		return privateScopeNone
	}

	// Globs expand remotely, so an exclude can't be anchored for them.
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
		// A trailing "." or ".." makes the transfer root, and so the exclude, unpredictable.
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

// globMayReachPrivate reports whether a glob could expand to the private filesystem, a path
// within it, or a parent. Globs don't cross "/", so elements are matched depth by depth.
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

// isWithin reports whether child is parent or below it.
func isWithin(child, parent string) bool {
	if child == parent || parent == "/" {
		return true
	}

	return strings.HasPrefix(child, parent+"/")
}

// privateExcludePattern returns an exclude for the private filesystem, anchored to the
// transfer root: "/mnt/" syncs from /mnt, "/mnt" syncs from /.
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
