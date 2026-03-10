package lesson

import (
	"fmt"
	"strings"

	"github.com/naorpeled/aitutor/internal/ui"
	"github.com/naorpeled/aitutor/pkg/types"
)

// RenderTheory renders a slice of TheoryBlocks into styled text.
func RenderTheory(blocks []types.TheoryBlock, width int) string {
	var parts []string

	for i, b := range blocks {
		switch b.Kind {
		case types.Heading:
			if i > 0 {
				parts = append(parts, "")
			}
			parts = append(parts, ui.HeadingStyle.Width(width).Render(b.Content))
		case types.Paragraph:
			parts = append(parts, ui.ParagraphStyle.Width(width).Render(b.Content))
		case types.Code:
			parts = append(parts, ui.CodeStyle.Width(width).Render(b.Content))
		case types.Callout:
			parts = append(parts, ui.CalloutStyle.Width(width-4).Render("💡 "+b.Content))
		case types.Bullet:
			lines := strings.Split(b.Content, "\n")
			for _, line := range lines {
				parts = append(parts, ui.BulletStyle.Width(width).Render(fmt.Sprintf("• %s", line)))
			}
		}
	}

	return strings.Join(parts, "\n")
}
