// Package runchart provides a terminal-based time-series chart renderer
// inspired by the sampler runchart component. It renders data points over
// time using braille Unicode characters for high-resolution terminal output.
//
// The chart supports multiple named lines, automatic Y-axis scaling,
// time-windowed data retention, and a legend with cur/min/max statistics.
package runchart

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

// Chart is a time-series run chart that accumulates data points and renders
// them as a string suitable for terminal display.
type Chart struct {
	mu         sync.Mutex
	title      string
	lines      []*Line
	width      int // chart width in terminal columns
	height     int // chart height in terminal rows (plot area only)
	window     time.Duration
	showLegend bool
	minValue   *float64 // optional floor for the Y-axis
	maxValue   *float64 // optional ceiling for the Y-axis
	titleColor string   // ANSI color code for the title underline
}

// Line represents a single named data series on the chart.
type Line struct {
	Label  string
	Color  string // ANSI color code (e.g. "\033[31m" for red)
	points []Point
	min    float64
	max    float64
}

// Point is a single timestamped value.
type Point struct {
	Value float64
	Time  time.Time
}

// Option configures a Chart.
type Option func(*Chart)

// WithTitle sets the chart title.
func WithTitle(title string) Option {
	return func(c *Chart) { c.title = title }
}

// WithSize sets the chart dimensions in terminal columns and rows.
// Width is used for the plot area; height is for the plot area excluding
// the title, legend, and axis labels.
func WithSize(width, height int) Option {
	return func(c *Chart) {
		c.width = width
		c.height = height
	}
}

// WithWindow sets the time window of data to display.
// Points older than now minus window are pruned on each render.
func WithWindow(d time.Duration) Option {
	return func(c *Chart) { c.window = d }
}

// WithLegend enables or disables the legend display.
func WithLegend(enabled bool) Option {
	return func(c *Chart) { c.showLegend = enabled }
}

// WithMinValue sets a fixed minimum value for the Y-axis.
// This prevents the axis from displaying values below the given floor.
func WithMinValue(v float64) Option {
	return func(c *Chart) { c.minValue = &v }
}

// WithMaxValue sets a fixed maximum value for the Y-axis.
// This prevents the axis from displaying values above the given ceiling.
func WithMaxValue(v float64) Option {
	return func(c *Chart) { c.maxValue = &v }
}

// WithTitleColor sets the ANSI color code used for the title underline.
// For example, "\033[38;2;238;86;34m" for orange.
func WithTitleColor(color string) Option {
	return func(c *Chart) { c.titleColor = color }
}

// New creates a new Chart with the given options.
// Defaults: 80 columns wide, 15 rows tall, 60s window, legend disabled.
func New(opts ...Option) *Chart {
	c := &Chart{
		width:      80,
		height:     15,
		window:     60 * time.Second,
		showLegend: false,
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// AddLine adds a named data series to the chart.
// color is an ANSI escape code like "\033[31m" for red. Use "" for default (white).
func (c *Chart) AddLine(label, color string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lines = append(c.lines, &Line{
		Label:  label,
		Color:  color,
		points: nil,
		min:    math.MaxFloat64,
		max:    -math.MaxFloat64,
	})
}

// Push adds a data point to the named line. If the label doesn't exist
// the point is silently dropped.
func (c *Chart) Push(label string, value float64) {
	c.PushAt(label, value, time.Now())
}

// PushAt adds a data point with an explicit timestamp.
func (c *Chart) PushAt(label string, value float64, t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, l := range c.lines {
		if l.Label == label {
			l.points = append(l.points, Point{Value: value, Time: t})
			if value < l.min {
				l.min = value
			}
			if value > l.max {
				l.max = value
			}
			return
		}
	}
}

// Render returns the chart as a string suitable for printing to a terminal.
func (c *Chart) Render() string {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.prune()

	var buf strings.Builder

	// Determine global value range across all lines.
	gmin, gmax := c.globalExtrema()
	if gmin == gmax {
		gmin -= 1
		gmax += 1
	}

	// Determine time range.
	now := time.Now()
	tmin := now.Add(-c.window)
	tmax := now

	// Y-axis label width.
	yLabelWidth := max(len(formatVal(gmin)), len(formatVal(gmax))) + 1

	// Plot area dimensions.
	plotWidth := c.width - yLabelWidth - 1 // -1 for axis line
	plotHeight := c.height
	if plotWidth < 4 {
		plotWidth = 4
	}
	if plotHeight < 2 {
		plotHeight = 2
	}

	// Title.
	if c.title != "" {
		pad := (c.width - len(c.title)) / 2
		if pad < 0 {
			pad = 0
		}
		buf.WriteString(strings.Repeat(" ", pad))
		buf.WriteString("\033[1m") // bold
		buf.WriteString(c.title)
		buf.WriteString("\033[0m") // reset
		buf.WriteString("\n")
		if c.titleColor != "" {
			buf.WriteString(c.titleColor)
		}
		buf.WriteString(strings.Repeat("─", c.width))
		if c.titleColor != "" {
			buf.WriteString("\033[0m")
		}
		buf.WriteString("\n")
	}

	// Build braille canvas: each cell is 2 wide x 4 tall in dot space.
	canvasW := plotWidth * 2
	canvasH := plotHeight * 4
	canvas := make([][]map[string]bool, canvasH)
	for i := range canvas {
		canvas[i] = make([]map[string]bool, canvasW)
		for j := range canvas[i] {
			canvas[i][j] = make(map[string]bool)
		}
	}

	// Map each line's points onto the canvas.
	for _, line := range c.lines {
		color := line.Color
		if color == "" {
			color = "default"
		}
		var prevX, prevY int
		first := true
		for _, p := range line.points {
			if p.Time.Before(tmin) || p.Time.After(tmax) {
				continue
			}
			// Map time -> x in canvas coordinates.
			xf := float64(p.Time.Sub(tmin).Nanoseconds()) / float64(tmax.Sub(tmin).Nanoseconds()) * float64(canvasW-1)
			x := int(math.Round(xf))
			// Map value -> y in canvas coordinates (0=top, canvasH-1=bottom).
			yf := (1.0 - (p.Value-gmin)/(gmax-gmin)) * float64(canvasH-1)
			y := int(math.Round(yf))
			y = clamp(y, 0, canvasH-1)
			x = clamp(x, 0, canvasW-1)

			if !first {
				// Draw line between previous point and current.
				drawLine(canvas, prevX, prevY, x, y, color)
			}
			canvas[y][x][color] = true
			prevX, prevY = x, y
			first = false
		}
	}

	// Render canvas into character rows.
	reset := "\033[0m"
	for row := 0; row < plotHeight; row++ {
		// Y-axis label.
		yVal := gmax - (gmax-gmin)*float64(row)/float64(plotHeight-1)
		label := formatVal(yVal)
		buf.WriteString(fmt.Sprintf("%*s", yLabelWidth, label))
		buf.WriteString("┤")

		for col := 0; col < plotWidth; col++ {
			// Gather 2x4 braille dots for this cell.
			colors := gatherColors(canvas, row, col)
			if len(colors) == 0 {
				buf.WriteRune(' ')
				continue
			}
			// Pick a color (first line wins for overlap).
			chosenColor := ""
			for _, line := range c.lines {
				clr := line.Color
				if clr == "" {
					clr = "default"
				}
				if colors[clr] {
					chosenColor = line.Color
					break
				}
			}

			ch := brailleChar(canvas, row, col)
			if chosenColor != "" {
				buf.WriteString(chosenColor)
			}
			buf.WriteRune(ch)
			if chosenColor != "" {
				buf.WriteString(reset)
			}
		}
		buf.WriteString("\n")
	}

	// Compute tick positions along the X-axis.
	// Space labels roughly every 15 columns, with a minimum of 2 ticks (start and end).
	const labelSpacing = 15
	numTicks := plotWidth / labelSpacing
	if numTicks < 2 {
		numTicks = 2
	}

	tickPositions := make([]int, numTicks)
	for i := 0; i < numTicks; i++ {
		tickPositions[i] = i * (plotWidth - 1) / (numTicks - 1)
	}

	// Build a set for quick lookup of tick positions.
	tickSet := make(map[int]bool, len(tickPositions))
	for _, pos := range tickPositions {
		tickSet[pos] = true
	}

	// X-axis line with tick marks.
	buf.WriteString(strings.Repeat(" ", yLabelWidth))
	buf.WriteString("└")
	for col := 0; col < plotWidth; col++ {
		if tickSet[col] {
			buf.WriteRune('┴')
		} else {
			buf.WriteRune('─')
		}
	}
	buf.WriteString("\n")

	// Time labels row — place each label centered under its tick mark.
	labelRow := make([]byte, plotWidth+yLabelWidth+1)
	for i := range labelRow {
		labelRow[i] = ' '
	}

	for _, pos := range tickPositions {
		// Interpolate the time at this tick position.
		frac := float64(pos) / float64(plotWidth-1)
		tickTime := tmin.Add(time.Duration(frac * float64(tmax.Sub(tmin))))
		label := tickTime.Format("15:04:05")

		// Center the label under the tick mark.
		absPos := yLabelWidth + 1 + pos
		start := absPos - len(label)/2
		if start < 0 {
			start = 0
		}
		if start+len(label) > len(labelRow) {
			start = len(labelRow) - len(label)
		}

		copy(labelRow[start:], label)
	}

	buf.WriteString(strings.TrimRight(string(labelRow), " "))
	buf.WriteString("\n")

	// Legend.
	if c.showLegend && len(c.lines) > 0 {
		buf.WriteString("\n")
		for _, line := range c.lines {
			cur := float64(0)
			if len(line.points) > 0 {
				cur = line.points[len(line.points)-1].Value
			}
			lmin, lmax := line.min, line.max
			if lmin == math.MaxFloat64 {
				lmin = 0
			}
			if lmax == -math.MaxFloat64 {
				lmax = 0
			}
			prefix := ""
			suffix := ""
			if line.Color != "" {
				prefix = line.Color
				suffix = reset
			}
			buf.WriteString(fmt.Sprintf("  %s■%s %-12s cur: %-8s min: %-8s max: %-8s\n",
				prefix, suffix,
				line.Label,
				formatVal(cur),
				formatVal(lmin),
				formatVal(lmax),
			))
		}
	}

	return buf.String()
}

// prune removes points older than the time window.
func (c *Chart) prune() {
	cutoff := time.Now().Add(-c.window - time.Second*10) // keep small buffer
	for _, line := range c.lines {
		idx := 0
		for idx < len(line.points) && line.points[idx].Time.Before(cutoff) {
			idx++
		}
		if idx > 0 {
			line.points = append(line.points[:0], line.points[idx:]...)
			// Recompute extrema after pruning.
			line.min = math.MaxFloat64
			line.max = -math.MaxFloat64
			for _, p := range line.points {
				if p.Value < line.min {
					line.min = p.Value
				}
				if p.Value > line.max {
					line.max = p.Value
				}
			}
		}
	}
}

// globalExtrema returns the min/max across all lines.
func (c *Chart) globalExtrema() (float64, float64) {
	gmin := math.MaxFloat64
	gmax := -math.MaxFloat64
	for _, line := range c.lines {
		for _, p := range line.points {
			if p.Value < gmin {
				gmin = p.Value
			}
			if p.Value > gmax {
				gmax = p.Value
			}
		}
	}
	if gmin == math.MaxFloat64 {
		return 0, 100
	}
	// Add 5% padding.
	spread := gmax - gmin
	if spread == 0 {
		spread = 1
	}
	gmin -= spread * 0.05
	gmax += spread * 0.05

	// Apply optional minimum value floor.
	if c.minValue != nil {
		gmin = *c.minValue
	}

	// Apply optional maximum value ceiling.
	if c.maxValue != nil {
		gmax = *c.maxValue
	}

	return gmin, gmax
}

// formatVal formats a float for axis labels.
func formatVal(v float64) string {
	return fmt.Sprintf("%.0f", v)
}

// Braille character encoding.
// Each braille character encodes a 2-wide x 4-tall dot pattern.
// Dot positions in a cell (col 0-1, row 0-3):
//
//	(0,0) (1,0)    => bit 0, bit 3
//	(0,1) (1,1)    => bit 1, bit 4
//	(0,2) (1,2)    => bit 2, bit 5
//	(0,3) (1,3)    => bit 6, bit 7
var brailleDots = [4][2]uint8{
	{0, 3},
	{1, 4},
	{2, 5},
	{6, 7},
}

// brailleChar computes the braille character for the cell at (row, col)
// in the plot grid.
func brailleChar(canvas [][]map[string]bool, row, col int) rune {
	var bits uint8
	for dy := 0; dy < 4; dy++ {
		for dx := 0; dx < 2; dx++ {
			cy := row*4 + dy
			cx := col*2 + dx
			if cy < len(canvas) && cx < len(canvas[cy]) && len(canvas[cy][cx]) > 0 {
				bits |= 1 << brailleDots[dy][dx]
			}
		}
	}
	if bits == 0 {
		return ' '
	}
	return rune(0x2800 + int(bits))
}

// gatherColors returns the set of color keys present in a braille cell.
func gatherColors(canvas [][]map[string]bool, row, col int) map[string]bool {
	colors := make(map[string]bool)
	for dy := 0; dy < 4; dy++ {
		for dx := 0; dx < 2; dx++ {
			cy := row*4 + dy
			cx := col*2 + dx
			if cy < len(canvas) && cx < len(canvas[cy]) {
				for k := range canvas[cy][cx] {
					colors[k] = true
				}
			}
		}
	}
	return colors
}

// drawLine uses Bresenham's algorithm to draw a line on the canvas.
func drawLine(canvas [][]map[string]bool, x0, y0, x1, y1 int, color string) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		if y0 >= 0 && y0 < len(canvas) && x0 >= 0 && x0 < len(canvas[y0]) {
			canvas[y0][x0][color] = true
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
