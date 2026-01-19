package term

import (
	"bufio"
	"fmt"
	"os"

	termios "github.com/pkg/term/termios"
	unix "golang.org/x/sys/unix"
)

// Disables input echo in the current terminal. Intended to be used in conjunction with [EnableEcho].
// Returns an error if [termios] does.
//
// Example:
//
//	fd := os.Stdin.Fd()
//	if err := term.DisableEcho(fd); err != nil {
//	    panic("no!")
//	}
//	defer term.EnableEcho(fd)
//
// [termios]: https://pkg.go.dev/github.com/pkg/term/termios
func DisableEcho(fd uintptr) error {
	t := &unix.Termios{}
	err := termios.Tcgetattr(fd, t)
	if err != nil {
		return err
	}
	t.Lflag = t.Lflag & (^uint32(unix.ECHO))
	t.Lflag = t.Lflag | unix.ISIG
	err = termios.Tcsetattr(fd, termios.TCSANOW, t)
	if err != nil {
		return err
	}
	return nil
}

// Enables input echo in the current terminal. Intended to be used in conjunction with [DisableEcho].
// Returns an error if [termios] does.
//
// Example:
//
//	fd := os.Stdin.Fd()
//	if err := term.DisableEcho(fd); err != nil {
//	    panic("no!")
//	}
//	defer term.EnableEcho(fd)
//
// [termios]: https://pkg.go.dev/github.com/pkg/term/termios
func EnableEcho(fd uintptr) error {
	t := &unix.Termios{}
	err := termios.Tcgetattr(fd, t)
	if err != nil {
		return err
	}
	t.Lflag = t.Lflag | unix.ECHO
	err = termios.Tcsetattr(fd, termios.TCSANOW, t)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Mask using *. Requires reading char-by-char & understanding erase characters.

// Reads a string from stdin until newline is encountered, temporarily disabling terminal echo until done.
// Errors will be returned if editing terminal properties fails or if user sends a premature EOF (e.g. via CTRL+D).
// More details are in the [Linux termios manpage] and [sys/unix] Go docs.
//
// [Linux termios manpage]: https://www.man7.org/linux/man-pages/man3/termios.3.html
// [sys/unix]: https://pkg.go.dev/golang.org/x/sys@v0.40.0/unix#ECHO
func PromptSensitive(prompt string) (string, error) {
	fmt.Fprintf(os.Stderr, "%s: ", prompt)
	fd := os.Stdin.Fd()
	if err := DisableEcho(fd); err != nil {
		return "", err
	}
	defer EnableEcho(fd)

	lineScanner := bufio.NewScanner(os.Stdin)
	var input string
	if lineScanner.Scan() {
		input = lineScanner.Text()
	} else {
		return "", fmt.Errorf("no input given")
	}
	fmt.Fprintln(os.Stderr)

	return input, nil
}
