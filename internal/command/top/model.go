package top

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	lipglossv1 "github.com/charmbracelet/lipgloss"
	"github.com/skpr/api/pb"

	"charm.land/lipgloss/v2"

	skprcolor "github.com/skpr/cli/internal/color"
)

const (
	helpHeight = 2
)

var (
	bgColor = lipgloss.Color(skprcolor.HexBlack)

	chartStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBackground(bgColor).
			Background(bgColor).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Background(bgColor).
			Padding(1, 0, 0, 1)

	sectionLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Background(bgColor)

	sectionLineStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(skprcolor.HexOrange)).
				Background(bgColor)

)

type tickMsg time.Time

type metricsMsg struct {
	resourceUsage        *pb.ResourceUsageResponse
	requests             *pb.RequestsResponse
	cacheRatio           *pb.CacheRatioResponse
	invalidationRequests *pb.InvalidationRequestsResponse
	invalidationPaths    *pb.InvalidationPathsResponse
	responseTimes        *pb.ResponseTimesResponse
	responseCodes        *pb.ResponseCodesResponse
}

type metricsErrMsg struct {
	err error
}

type model struct {
	metricsMsg

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

	refreshInterval  time.Duration
	secondsRemaining int
	fetching         bool
	fetchMetrics     func() tea.Msg

	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

func newModel(initial metricsMsg, refreshInterval time.Duration, fetchMetrics func() tea.Msg) model {
	return model{
		metricsMsg:       initial,
		refreshInterval:  refreshInterval,
		secondsRemaining: int(refreshInterval.Seconds()),
		fetchMetrics:     fetchMetrics,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case tickMsg:
		m.secondsRemaining--
		if m.secondsRemaining <= 0 && !m.fetching {
			m.fetching = true
			cmds = append(cmds, func() tea.Msg {
				return m.fetchMetrics()
			})
		}
		cmds = append(cmds, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return tickMsg(t)
		}))
		return m, tea.Batch(cmds...)

	case metricsMsg:
		m.metricsMsg = msg
		m.fetching = false
		m.secondsRemaining = int(m.refreshInterval.Seconds())

		if m.ready {
			m.rebuildGraphs()
			m.viewport.SetContent(m.renderCharts())
		}

		return m, nil

	case metricsErrMsg:
		m.fetching = false
		m.secondsRemaining = int(m.refreshInterval.Seconds())
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-helpHeight)
			m.viewport.Style = lipglossv1.NewStyle().Background(lipglossv1.Color(skprcolor.HexBlack))
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - helpHeight
		}

		m.rebuildGraphs()
		m.viewport.SetContent(m.renderCharts())
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if !m.ready {
		return "Loading..."
	}

	countdown := fmt.Sprintf("Refresh in %ds", m.secondsRemaining)
	if m.fetching {
		countdown = "Refreshing..."
	}

	help := helpStyle.Render(fmt.Sprintf("↑/↓ scroll | q quit | %s", countdown))

	content := lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), help)

	return fillBg(content, m.width)
}

func (m *model) rebuildGraphs() {
	frameSize := chartStyle.GetHorizontalFrameSize()

	chartWidth2 := m.width/2 - frameSize
	if chartWidth2 < 20 {
		chartWidth2 = 20
	}

	chartWidth3 := m.width/3 - frameSize
	if chartWidth3 < 20 {
		chartWidth3 = 20
	}

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
}

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
	suffix := sectionLabelStyle.Render(" ")

	used := 3 + lipgloss.Width(name) + 1
	remaining := m.width - used
	if remaining < 0 {
		remaining = 0
	}
	trail := sectionLineStyle.Render(strings.Repeat("━", remaining))

	return "\n" + prefix + label + suffix + trail
}

func fillBg(s string, width int) string {
	const bg = "\033[48;2;0;0;0m"
	const reset = "\033[0m"

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		line = strings.ReplaceAll(line, "\x1b[m", "\x1b[m"+bg)
		line = strings.ReplaceAll(line, "\x1b[0m", "\x1b[0m"+bg)

		visible := lipgloss.Width(line)
		if visible < width {
			line += strings.Repeat(" ", width-visible)
		}
		lines[i] = bg + line + reset
	}
	return strings.Join(lines, "\n")
}

func (m model) renderCharts() string {
	kpi := m.sectionHeader("Key Performance Indicators")

	h1 := rowHeight(m.responseTimesGraph, m.responseCodesGraph)
	row1 := lipgloss.JoinHorizontal(lipgloss.Top,
		chartStyle.Height(h1).Render(m.responseTimesGraph),
		chartStyle.Height(h1).Render(m.responseCodesGraph),
	)

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
