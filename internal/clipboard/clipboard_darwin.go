package clipboard

import "os/exec"

//nolint:noctx // there is no real reason to use context
func Copy(pwd string) error {
	cmd := exec.Command("pbcopy")

	in, pipeErr := cmd.StdinPipe()
	if pipeErr != nil {
		return pipeErr
	}

	if cmdStartErr := cmd.Start(); cmdStartErr != nil {
		return cmdStartErr
	}

	if _, writeErr := in.Write([]byte(pwd)); writeErr != nil {
		return writeErr
	}

	if closeErr := in.Close(); closeErr != nil {
		return closeErr
	}

	return cmd.Wait()
}
