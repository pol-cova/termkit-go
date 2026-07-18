package pixel

import (
	"os"
	"strings"
)

type Protocol int

const (
	ProtocolANSI Protocol = iota
	ProtocolKitty
	ProtocolITerm2
)

func DetectProtocol() Protocol {
	program := strings.ToLower(os.Getenv("TERM_PROGRAM"))
	term := strings.ToLower(os.Getenv("TERM"))
	switch {
	case program == "kitty" || strings.Contains(term, "kitty"):
		return ProtocolKitty
	case program == "iterm.app" || strings.Contains(term, "iterm"):
		return ProtocolITerm2
	default:
		return ProtocolANSI
	}
}

func (c Canvas) Auto() (string, error) {
	return c.Render(DetectProtocol())
}

func (c Canvas) Render(protocol Protocol) (string, error) {
	switch protocol {
	case ProtocolKitty:
		return c.Kitty()
	case ProtocolITerm2:
		return c.ITerm2()
	default:
		return c.ANSI(), nil
	}
}
