package stream

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTargetCover_SetForm(t *testing.T) {
	target := NewTargetCover(SpiderCover, DownloadCover)
	require.Equal(t, target.Form(), SpiderCover)
	require.Equal(t, target.To(), DownloadCover)

	target.SetTo(SchedulerCover)
	require.Equal(t, target.Form(), SpiderCover)
	require.Equal(t, target.To(), SchedulerCover)

	target.SetForm(OutputDeviceCover)
	require.Equal(t, target.Form(), OutputDeviceCover)
	require.Equal(t, target.To(), SchedulerCover)

	target.SetForm(SchedulerCover)
	target.SetTo(DownloadCover)

	require.Equal(t, target.Form(), SchedulerCover)
	require.Equal(t, target.To(), DownloadCover)
}
