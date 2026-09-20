package jsonview

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestExplorerTextReflowsOnResize(t *testing.T) {
	t.Parallel()
	view, err := newTextView("", gjson.Parse(`"alpha beta gamma delta"`))
	require.NoError(t, err)
	viewer := &JSONViewer{stack: []JSONView{view}}
	for _, width := range []int{24, 10, 24} {
		viewer.Update(tea.WindowSizeMsg{Width: width + borderPadding, Height: 20})
		require.Equal(t, strings.Fields(view.data.Str), strings.Fields(view.View()), "text was clipped at width %d", width)
		if width == 10 {
			require.Equal(t, 3, view.viewport.TotalLineCount())
		} else {
			require.Equal(t, 1, view.viewport.TotalLineCount())
		}
	}
}

func TestExplorerTextResizePreservesScrollingAndEscapes(t *testing.T) {
	t.Parallel()
	view, err := newTextView("", gjson.Parse(`"one two three four five six seven eight nine ten \u001b[31m"`))
	require.NoError(t, err)
	view.Resize(10, 8)
	view.viewport.GotoBottom()
	before := view.viewport.YOffset
	require.Positive(t, before)
	view.Resize(10, 9)
	require.Equal(t, before, view.viewport.YOffset)
	view.Resize(100, 20)
	require.Zero(t, view.viewport.YOffset, "rewrapping must clamp an offset past the last line")
	require.Equal(t, strings.Fields(SanitizeTerminalString(view.data.Str)), strings.Fields(view.View()))
	requireNoRawTerminalControls(t, strings.TrimSpace(view.View()))
}
