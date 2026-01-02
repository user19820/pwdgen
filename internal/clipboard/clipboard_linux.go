package clipboard

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

// windowManager is an enum of different possible
// WM who manage clipboards (since linux doesn't have
// a centralized implementation like windows or macOS).
//
// Its used in order to find out which clipboard needs
// to be used.
type windowManager string

const (
	windowManagerX11     windowManager = "x11"
	windowManagerWayland windowManager = "wayland"
	windowManagerInvalid windowManager = ""
)

func Copy(pwd string) error {
	wMng, wMngErr := detectX11OrWayland()
	if wMngErr != nil {
		return wMngErr
	}

	switch wMng {
	case windowManagerWayland:
		wlCopyFound, wlCopyErr := detectWlCopy()
		if wlCopyErr != nil {
			return wlCopyErr
		}

		if !wlCopyFound {
			return errors.New("wl-copy needs to be installed for Wayland based linux")
		}

		return copyToClipboardWayland(pwd)
	case windowManagerX11:
		xclipFound, xclipErr := detectXclip()
		if xclipErr != nil {
			return xclipErr
		}

		if !xclipFound {
			return errors.New("xclip needs to be installed for X11 based linux")
		}

		return copyToClipboardX11(pwd)
	case windowManagerInvalid:
		panic("SHOULD BE UNREACHABLE")
	default:
		panic("SHOULD BE UNREACHABLE")
	}
}

func copyToClipboardWayland(pwd string) error {
	cmd := exec.Command("wl-copy", "--paste-once")

	wlCopyIn, pipeErr := cmd.StdinPipe()
	if pipeErr != nil {
		return pipeErr
	}

	if cmdStartErr := cmd.Start(); cmdStartErr != nil {
		return cmdStartErr
	}

	if _, writeErr := wlCopyIn.Write([]byte(pwd)); writeErr != nil {
		return writeErr
	}

	if closeErr := wlCopyIn.Close(); closeErr != nil {
		return closeErr
	}

	return cmd.Wait()
}

func copyToClipboardX11(pwd string) error {
	cmd := exec.Command("xclip")

	xclipIn, pipeErr := cmd.StdinPipe()
	if pipeErr != nil {
		return pipeErr
	}

	if cmdStartErr := cmd.Start(); cmdStartErr != nil {
		return cmdStartErr
	}

	if _, writeErr := xclipIn.Write([]byte(pwd)); writeErr != nil {
		return writeErr
	}

	if closeErr := xclipIn.Close(); closeErr != nil {
		return closeErr
	}

	return cmd.Wait()
}

func detectX11OrWayland() (windowManager, error) {
	windowEnv, isSet := os.LookupEnv("XDG_SESSION_TYPE")
	if !isSet {
		return windowManagerInvalid, errors.New("XDG_SESSION_TYPE is not set")
	}

	switch windowEnv {
	case "wayland":
		return windowManagerWayland, nil
	case "x11":
		return windowManagerX11, nil
	default:
		return windowManagerInvalid, errors.New("XDG_SESSION_TYPE is invalid (needs to be wayland or x11)")
	}
}

func detectXclip() (bool, error) {
	cmd := exec.Command("xclip", "-v")

	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	if bytes.Contains(out, []byte("command not found")) {
		return false, nil
	}

	return true, nil
}

func detectWlCopy() (bool, error) {
	cmd := exec.Command("wl-copy", "-v")

	out, err := cmd.Output()
	if err != nil {
		return false, err
	}

	if bytes.Contains(out, []byte("command not found")) {
		return false, nil
	}

	return true, nil
}
