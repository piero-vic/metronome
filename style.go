package main

import "github.com/charmbracelet/lipgloss"

var wrapper = lipgloss.NewStyle().
	Padding(1, 1, 2, 1)

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("0")).
	Background(lipgloss.Color("3")).
	Padding(0, 1).
	Margin(0, 0, 1, 0)

var statusStyle = lipgloss.NewStyle().
	Margin(0, 0, 1, 0)

var indicatorStyle = lipgloss.NewStyle().
	Margin(0, 0, 1, 0)

var activeIndicatorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("3"))
