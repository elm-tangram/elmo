package diagnostic

import (
	"fmt"
	"strings"

	"github.com/erizocosmico/elmo/source"
	"github.com/erizocosmico/elmo/token"
)

// Diagnostic is the common interface of an error or a warning that happened
// in any step of the compilation process.
type Diagnostic interface {
	// Severity of the diagnostic.
	Severity() Severity
	// Msg of the diagnostic
	Msg() string
	// Line where the error happened
	Line() int64
	// Column where the error happened
	Column() int64
	// HasRegion reports whether the diagnostic contains a region of code.
	HasRegion() bool
	// Lines with the region of the diagnosed code.
	Lines() []source.Line
}

type regionDiagnostic struct {
	severity Severity
	msg      Msg
	pos      *token.Position
	lines    []source.Line
}

type msgDiagnostic struct {
	severity Severity
	msg      Msg
	pos      *token.Position
}

// Msg is a human-readable message of a diagnostic.
type Msg interface {
	fmt.Stringer
}

// NewRegionDiagnostic creates a new diagnostic for a specific region of the
// source code.
func NewRegionDiagnostic(severity Severity, msg Msg, pos *token.Position, region []source.Line) Diagnostic {
	return &regionDiagnostic{severity, msg, pos, region}
}

// NewMsgDiagnostic creates a new diagnostic that is not for a specific region
// of the source code.
func NewMsgDiagnostic(severity Severity, msg Msg, pos *token.Position) Diagnostic {
	return &msgDiagnostic{severity, msg, pos}
}

func (d *regionDiagnostic) Severity() Severity   { return d.severity }
func (d *msgDiagnostic) Severity() Severity      { return d.severity }
func (d *regionDiagnostic) Msg() string          { return d.msg.String() }
func (d *msgDiagnostic) Msg() string             { return d.msg.String() }
func (d *regionDiagnostic) Line() int64          { return int64(d.pos.Line) }
func (d *msgDiagnostic) Line() int64             { return int64(d.pos.Line) }
func (d *regionDiagnostic) Column() int64        { return int64(d.pos.Column) }
func (d *msgDiagnostic) Column() int64           { return int64(d.pos.Column) }
func (d *regionDiagnostic) HasRegion() bool      { return true }
func (d *msgDiagnostic) HasRegion() bool         { return false }
func (d *regionDiagnostic) Lines() []source.Line { return d.lines }
func (d *msgDiagnostic) Lines() []source.Line    { return nil }

// UnexpectedEOF returns a diagnostic message saying that EOF was not expected,
// but one of the given token types.
func UnexpectedEOF(expecting ...token.Type) Msg {
	return &parseError{&unexpectedEOF{typeList(expecting)}}
}

type parseError struct {
	err Msg
}

func (e *parseError) String() string {
	return "I ran into something unexpected parsing your code: " + e.err.String()
}

type unexpectedEOF struct {
	expecting typeList
}

func (e *unexpectedEOF) String() string {
	return fmt.Sprintf("Unexpected end of file, I was expecting %s instead", e.expecting)
}

// Expecting returns a diagnostic message saying that the found token was not
// what the parser was expecting.
func Expecting(found token.Type, expecting ...token.Type) Msg {
	return &parseError{
		&errExpecting{
			found,
			typeList(expecting),
		},
	}
}

type errExpecting struct {
	found     token.Type
	expecting typeList
}

func (e *errExpecting) String() string {
	return fmt.Sprintf("I found %s, but I was expecting %s instead", e.found, e.expecting)
}

// ParseError returns a custom diagnostic message.
func ParseError(msg string) Msg {
	return &parseError{&msgErr{msg}}
}

type msgErr struct {
	msg string
}

func (m *msgErr) String() string {
	return m.msg
}

type typeList []token.Type

func (tl typeList) String() string {
	if len(tl) == 0 {
		return "nothing"
	}

	if len(tl) == 1 {
		return fmt.Sprint(tl[0])
	}

	var types = make([]string, len(tl)-1)
	for i, t := range tl[:len(tl)-1] {
		types[i] = fmt.Sprint(t)
	}

	return fmt.Sprintf(
		"%s or %s",
		strings.Join(types, ", "),
		tl[len(tl)-1],
	)
}