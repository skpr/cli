package top

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/skpr/api/pb"

	"charm.land/lipgloss/v2"

	skprcolor "github.com/skpr/cli/internal/color"
)

const (
	// Height reserved for the help bar (1 line padding + 1 line text).
	helpHeight = 2
)

var (
	chartStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(1, 0, 0, 1)

	sectionLabelStyle = lipgloss.NewStyle().
				Bold(true)

	sectionLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(skprcolor.HexOrange))
)

type model struct {
	// API responses.
	resourceUsage        *pb.ResourceUsageResponse
	requests             *pb.RequestsResponse
	cacheRatio           *pb.CacheRatioResponse
	invalidationRequests *pb.InvalidationRequestsResponse
	invalidationPaths    *pb.InvalidationPathsResponse
	responseTimes        *pb.ResponseTimesResponse
	responseCodes        *pb.ResponseCodesResponse

	// Rendered graphs.
	cpuGraph                  string
	memoryGraph               string
	replicasGraph             string
	processesGraph            string
	listenQueueGraph          string
	requestsGraph             string
	responseTimesGraph        string
	responseCodesGraph        string
	cacheHitRatioGraph        string
	invalidationRequestsGraph string
	invalidationPathsGraph    string

	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

func newModel(
	resourceUsage *pb.ResourceUsageResponse,
	requests *pb.RequestsResponse,
	cacheRatio *pb.CacheRatioResponse,
	invalidationRequests *pb.InvalidationRequestsResponse,
	invalidationPaths *pb.InvalidationPathsResponse,
	responseTimes *pb.ResponseTimesResponse,
	responseCodes *pb.ResponseCodesResponse,
) model {
	return model{
		resourceUsage:        resourceUsage,
		requests:             requests,
		cacheRatio:           cacheRatio,
		invalidationRequests: invalidationRequests,
		invalidationPaths:    invalidationPaths,
		responseTimes:        responseTimes,
		responseCodes:        responseCodes,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-helpHeight)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - helpHeight
		}

		frameSize := chartStyle.GetHorizontalFrameSize()

		// 2-column width.
		chartWidth2 := m.width/2 - frameSize
		if chartWidth2 < 20 {
			chartWidth2 = 20
		}

		// 3-column width (for Row 2: CPU | Memory | Replicas).
		chartWidth3 := m.width/3 - frameSize
		if chartWidth3 < 20 {
			chartWidth3 = 20
		}

		// Row 1: Response Times | Response Codes (2 cols).
		if graph, err := getResponseTimesGraph(m.responseTimes, chartWidth2); err != nil {
			m.responseTimesGraph = fmt.Sprintf("Response Times error: %v", err)
		} else {
			m.responseTimesGraph = graph
		}

		if graph, err := getResponseCodesGraph(m.responseCodes, chartWidth2); err != nil {
			m.responseCodesGraph = fmt.Sprintf("Response Codes error: %v", err)
		} else {
			m.responseCodesGraph = graph
		}

		// Row 2: CPU | Memory | Replicas (3 cols).
		if graph, err := getCPUGraph(m.resourceUsage, chartWidth3); err != nil {
			m.cpuGraph = fmt.Sprintf("CPU error: %v", err)
		} else {
			m.cpuGraph = graph
		}

		if graph, err := getMemoryGraph(m.resourceUsage, chartWidth3); err != nil {
			m.memoryGraph = fmt.Sprintf("Memory error: %v", err)
		} else {
			m.memoryGraph = graph
		}

		if graph, err := getReplicasGraph(m.resourceUsage, chartWidth3); err != nil {
			m.replicasGraph = fmt.Sprintf("Replicas error: %v", err)
		} else {
			m.replicasGraph = graph
		}

		// Row 3: Processes | Listen Queue (2 cols).
		if graph, err := getProcessesGraph(m.resourceUsage, chartWidth2); err != nil {
			m.processesGraph = fmt.Sprintf("Processes error: %v", err)
		} else {
			m.processesGraph = graph
		}

		if graph, err := getListenQueueGraph(m.resourceUsage, chartWidth2); err != nil {
			m.listenQueueGraph = fmt.Sprintf("Listen Queue error: %v", err)
		} else {
			m.listenQueueGraph = graph
		}

		// Row 4: Requests | Cache Hit Ratio (2 cols).
		if graph, err := getRequestsGraph(m.requests, chartWidth2); err != nil {
			m.requestsGraph = fmt.Sprintf("Requests error: %v", err)
		} else {
			m.requestsGraph = graph
		}

		if graph, err := getCacheHitRatioGraph(m.cacheRatio, chartWidth2); err != nil {
			m.cacheHitRatioGraph = fmt.Sprintf("Cache Hit Ratio error: %v", err)
		} else {
			m.cacheHitRatioGraph = graph
		}

		// Row 5: Invalidation Requests | Invalidation Paths (2 cols).
		if graph, err := getInvalidationRequestsGraph(m.invalidationRequests, chartWidth2); err != nil {
			m.invalidationRequestsGraph = fmt.Sprintf("Invalidation Requests error: %v", err)
		} else {
			m.invalidationRequestsGraph = graph
		}

		if graph, err := getInvalidationPathsGraph(m.invalidationPaths, chartWidth2); err != nil {
			m.invalidationPathsGraph = fmt.Sprintf("Invalidation Paths error: %v", err)
		} else {
			m.invalidationPathsGraph = graph
		}

		m.viewport.SetContent(m.renderCharts())
	}

	// Forward messages to the viewport for scroll handling.
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	help := helpStyle.Render("↑/↓ scroll | q quit")

	return lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), help)
}

// rowHeight returns the height needed for a chartStyle border to match the
// tallest content in the row. It accounts for lipgloss subtracting the
// vertical border size internally.
func rowHeight(graphs ...string) int {
	h := 0
	for _, g := range graphs {
		if gh := lipgloss.Height(g); gh > h {
			h = gh
		}
	}
	return h + chartStyle.GetVerticalBorderSize()
}

func (m model) sectionHeader(name string) string {
	prefix := sectionLineStyle.Render("━━ ")
	label := sectionLabelStyle.Render(name)
	suffix := " "

	// Fill the remaining width with the line.
	used := 3 + lipgloss.Width(name) + 1
	remaining := m.width - used
	if remaining < 0 {
		remaining = 0
	}
	trail := sectionLineStyle.Render(strings.Repeat("━", remaining))

	return "\n" + prefix + label + suffix + trail
}

// renderCharts composes all widget panels into a single string for the viewport.
func (m model) renderCharts() string {
	// Key Performance Indicators.
	kpi := m.sectionHeader("Key Performance Indicators")

	h1 := rowHeight(m.responseTimesGraph, m.responseCodesGraph)
	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h1).Render(m.responseTimesGraph),
		chartStyle.Height(h1).Render(m.responseCodesGraph),
	)

	// Resources.
	resources := m.sectionHeader("Resources")

	h2 := rowHeight(m.cpuGraph, m.memoryGraph, m.replicasGraph)
	row2 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h2).Render(m.cpuGraph),
		chartStyle.Height(h2).Render(m.memoryGraph),
		chartStyle.Height(h2).Render(m.replicasGraph),
	)

	h3 := rowHeight(m.processesGraph, m.listenQueueGraph)
	row3 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h3).Render(m.processesGraph),
		chartStyle.Height(h3).Render(m.listenQueueGraph),
	)

	// CDN.
	cdn := m.sectionHeader("CDN")

	h4 := rowHeight(m.requestsGraph, m.cacheHitRatioGraph)
	row4 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h4).Render(m.requestsGraph),
		chartStyle.Height(h4).Render(m.cacheHitRatioGraph),
	)

	h5 := rowHeight(m.invalidationRequestsGraph, m.invalidationPathsGraph)
	row5 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h5).Render(m.invalidationRequestsGraph),
		chartStyle.Height(h5).Render(m.invalidationPathsGraph),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		kpi, row1,
		resources, row2, row3,
		cdn, row4, row5,
	)
}
