package run9

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoxViewSnapIDAtMount(t *testing.T) {
	box := BoxView{BoxID: "my-box", BoxSnapID: "root", Volumes: []VolumeView{
		{SnapID: "cache", MountPath: "/cache"},
		{SnapID: "state", MountPath: "/state"},
	}}
	root, err := box.SnapIDAtMount("/")
	require.NoError(t, err)
	require.Equal(t, "root", root)
	state, err := box.SnapIDAtMount("/state")
	require.NoError(t, err)
	require.Equal(t, "state", state)
	_, err = box.SnapIDAtMount("/state/file")
	require.ErrorContains(t, err, `has no Snap at mount path "/state/file"`)
	_, err = box.SnapIDAtMount("/state/")
	require.Error(t, err)
	_, err = box.SnapIDAtMount("/cache/../state")
	require.Error(t, err)
	_, err = box.SnapIDAtMount("")
	require.Error(t, err)
	_, err = (BoxView{BoxID: "empty"}).SnapIDAtMount("/")
	require.Error(t, err)
}
