package run9

import "fmt"

// SnapIDAtMount resolves an exact mount path in this Box view to its Snap ID.
// Use "/" for the original Attached Snap, or a Volume's configured MountPath.
// It does not normalize paths, match subdirectories, or make a network request.
// This lookup does not establish readiness; Snap operations check current state
// on the server, including the owning Box's stopped state before a Fork.
func (b BoxView) SnapIDAtMount(mountPath string) (string, error) {
	if mountPath == "/" && b.BoxSnapID != "" {
		return b.BoxSnapID, nil
	}
	for _, volume := range b.Volumes {
		if volume.MountPath == mountPath && volume.SnapID != "" {
			return volume.SnapID, nil
		}
	}
	return "", fmt.Errorf("box %q has no Snap at mount path %q", b.BoxID, mountPath)
}
