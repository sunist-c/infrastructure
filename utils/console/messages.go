package console

import (
	"fmt"
	"strings"
)

func NewMessage() Messages {
	return Messages{builder: &strings.Builder{}}
}

type Messages struct {
	builder *strings.Builder
}

func (m Messages) format(color Control, message string) Messages {
	m.builder.WriteString(string(color))
	m.builder.WriteString(message)
	m.builder.WriteString(string(Reset))
	return m
}

func (m Messages) Index(idx int) Messages {
	m.builder.WriteString(fmt.Sprintf("%d.\t", idx))
	return m
}

func (m Messages) Red(message string) Messages {
	return m.format(ColorRed, message)
}

func (m Messages) Green(message string) Messages {
	return m.format(ColorGreen, message)
}

func (m Messages) Yellow(message string) Messages {
	return m.format(ColorYellow, message)
}

func (m Messages) Blue(message string) Messages {
	return m.format(ColorBlue, message)
}

func (m Messages) Magenta(message string) Messages {
	return m.format(ColorMagenta, message)
}

func (m Messages) Cyan(message string) Messages {
	return m.format(ColorCyan, message)
}

func (m Messages) White(message string) Messages {
	return m.format(ColorWhite, message)
}

func (m Messages) Bold(message string) Messages {
	return m.format(FormatBold, message)
}

func (m Messages) Dim(message string) Messages {
	return m.format(FormatDim, message)
}

func (m Messages) Underline(message string) Messages {
	return m.format(FormatUnder, message)
}

func (m Messages) Blink(message string) Messages {
	return m.format(FormatBlink, message)
}

func (m Messages) Reverse(message string) Messages {
	return m.format(FormatRev, message)
}

func (m Messages) Message(message string, args ...any) Messages {
	if len(args) == 0 {
		m.builder.WriteString(message)
		return m
	}

	m.builder.WriteString(fmt.Sprintf(message, args...))
	return m
}
