package main

import (
	"bytes"
	"flag"
	"log"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/flac"
	"github.com/gopxl/beep/speaker"

	_ "embed"
)

var (
	//go:embed metronome-strong-pulse.flac
	strongPulse []byte
	//go:embed metronome-weak-pulse.flac
	weakPulse []byte
)

const (
	MIN_TEMPO = 20
	MAX_TEMPO = 240

	MIN_BEATS = 1
	MAX_BEATS = 10
)

type tickMsg struct {
	time time.Time
	tag  int
}

type model struct {
	tempo   int
	beats   int
	playing bool

	currentBeat       int
	tag               int
	strongPulseBuffer *beep.Buffer
	weakPulseBuffer   *beep.Buffer
	help              help.Model
}

func (m model) Init() tea.Cmd {
	if m.playing {
		return func() tea.Msg { return tickMsg{time.Now(), m.tag} }
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.playing {
			return m, nil
		}

		if msg.tag > 0 && msg.tag != m.tag {
			return m, nil
		}

		if m.currentBeat >= m.beats {
			m.currentBeat = 1
		} else {
			m.currentBeat++
		}

		var streamer beep.StreamSeeker
		if m.currentBeat == 1 {
			streamer = m.strongPulseBuffer.Streamer(0, m.strongPulseBuffer.Len())
		} else {
			streamer = m.weakPulseBuffer.Streamer(0, m.weakPulseBuffer.Len())
		}

		speaker.Play(streamer)

		m.tag++
		return m, tick(bpmToDuration(m.tempo), m.tag)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.Up):
			m.tempo = clamp(MIN_TEMPO, MAX_TEMPO, m.tempo+1)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.Down):
			m.tempo = clamp(MIN_TEMPO, MAX_TEMPO, m.tempo-1)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.UpByFive):
			m.tempo = clamp(MIN_TEMPO, MAX_TEMPO, m.tempo+5)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.DownByFive):
			m.tempo = clamp(MIN_TEMPO, MAX_TEMPO, m.tempo-5)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.Right):
			m.beats = clamp(MIN_BEATS, MAX_BEATS, m.beats+1)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.Left):
			m.beats = clamp(MIN_BEATS, MAX_BEATS, m.beats-1)

			return m, nil

		case key.Matches(msg, DefaultKeyMap.Play):
			m.playing = !m.playing
			m.currentBeat = 0

			if m.playing {
				return m, func() tea.Msg { return tickMsg{time.Now(), m.tag} }
			}

			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	tempoStatus := strconv.Itoa(int(m.tempo)) + " bpm"
	beatStatus := strconv.Itoa(m.beats) + " beats"
	var playingStatus string
	if m.playing {
		playingStatus = "Playing"
	} else {
		playingStatus = "Paused"
	}
	separator := separatorStyle.Render(" • ")
	status := tempoStatus + separator + beatStatus + separator + playingStatus

	var indicator string
	for i := 1; i <= m.beats; i++ {
		if i == m.currentBeat {
			indicator += activeIndicatorStyle.Render("■ ")
		} else {
			indicator += "▪ "
		}
	}

	return wrapper.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			headerStyle.Render("metronome"),
			statusStyle.Render(status),
			indicatorStyle.Render(indicator),
			m.help.View(DefaultKeyMap),
		),
	)
}

func tick(t time.Duration, tag int) tea.Cmd {
	return tea.Tick(t, func(t time.Time) tea.Msg {
		return tickMsg{t, tag}
	})
}

func bpmToDuration(bpm int) time.Duration {
	return time.Duration((60 / float32(bpm)) * float32(time.Second))
}

func clamp(min, max, val int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

func main() {
	strongPulseStreamer, strongPulseFormat, err := flac.Decode(bytes.NewReader(strongPulse))
	if err != nil {
		log.Fatal(err)
	}

	strongPulseBuffer := beep.NewBuffer(strongPulseFormat)
	strongPulseBuffer.Append(strongPulseStreamer)

	if err = strongPulseStreamer.Close(); err != nil {
		log.Fatal(err)
	}

	weakPulseStreamer, weakPulseFormat, err := flac.Decode(bytes.NewReader(weakPulse))
	if err != nil {
		log.Fatal(err)
	}

	weakPulseBuffer := beep.NewBuffer(weakPulseFormat)
	weakPulseBuffer.Append(weakPulseStreamer)

	if err = weakPulseStreamer.Close(); err != nil {
		log.Fatal(err)
	}

	// NOTE: This assumes both strong and weak pulse have the same sample rate
	if err = speaker.Init(strongPulseFormat.SampleRate, strongPulseFormat.SampleRate.N(time.Second/100)); err != nil {
		log.Fatal(err)
	}

	var (
		tempo int
		beats int
		play  bool
	)
	flag.IntVar(&tempo, "t", 60, "tempo in bpm")
	flag.IntVar(&tempo, "tempo", 60, "tempo in bpm")
	flag.IntVar(&beats, "b", 4, "number of beats per measure")
	flag.IntVar(&beats, "beats", 4, "number of beats per measure")
	flag.BoolVar(&play, "p", false, "start playing")
	flag.BoolVar(&play, "play", false, "start playing")
	flag.Parse()

	m := model{
		tempo:   clamp(MIN_TEMPO, MAX_TEMPO, tempo),
		beats:   clamp(MIN_BEATS, MAX_BEATS, beats),
		playing: play,

		currentBeat:       0,
		strongPulseBuffer: strongPulseBuffer,
		weakPulseBuffer:   weakPulseBuffer,
		help:              help.New(),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		log.Fatal(err)
	}
}
