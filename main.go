package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"

	_ "embed"
)

//go:embed click.wav
var click []byte

type tickMsg time.Time

type model struct {
	bpm         int
	currentBeat int
	totalBeats  int
	playing     bool
	buffer      *beep.Buffer
}

func (m model) Init() tea.Cmd {
	if m.playing {
		return tick(bpmToDuration(m.bpm))
	}

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if m.currentBeat >= m.totalBeats {
			m.currentBeat = 1
		} else {
			m.currentBeat++
		}

		streamer := m.buffer.Streamer(0, m.buffer.Len())
		speaker.Play(streamer)

		if m.playing {
			return m, tick(bpmToDuration(m.bpm))
		}

		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.bpm = clamp(20, 200, m.bpm+1)

			return m, nil

		case "down", "j":
			m.bpm = clamp(20, 200, m.bpm-1)

			return m, nil

		case "right", "l":
			m.totalBeats = clamp(1, 10, m.totalBeats+1)

			return m, nil

		case "left", "h":
			m.totalBeats = clamp(1, 10, m.totalBeats-1)

			return m, nil

		case " ", "p":
			m.playing = !m.playing

			if m.playing {
				return m, tick(bpmToDuration(m.bpm))
			}

			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	var playingStatus string
	if m.playing {
		playingStatus = "Playing"
	} else {
		playingStatus = "Paused"
	}

	header := strconv.Itoa(int(m.bpm)) + " bpm | " +
		strconv.Itoa(m.totalBeats) + " beats | " +
		playingStatus

	var indicator string
	for i := 1; i <= m.totalBeats; i++ {
		if i == m.currentBeat {
			indicator += "■ "
		} else {
			indicator += "▪ "
		}
	}

	return header + "\n" + indicator + "\n"
}

func tick(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(t time.Time) tea.Msg {
		return tickMsg(t)
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
	streamer, format, err := wav.Decode(bytes.NewReader(click))
	if err != nil {
		log.Fatal(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/100))

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	streamer.Close()

	m := model{
		bpm:         60,
		currentBeat: 0,
		totalBeats:  4,
		playing:     false,
		buffer:      buffer,
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no, it didn't work:", err)
		os.Exit(1)
	}
}
