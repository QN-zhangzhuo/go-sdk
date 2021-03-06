package defaults

import (
	"os"
	"path/filepath"
	"runtime"
)

// SharedCredentialsFilename returns the SDK's default file path
// for the shared credentials file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.qiniu/credentials
//   - Windows: %USERPROFILE%\.qiniu\credentials
func SharedCredentialsFilename() string {
	return filepath.Join(UserHomeDir(), ".qiniu", "credentials")
}

// SharedConfigFilename returns the SDK's default file path for
// the shared config file.
//
// Builds the shared config file path based on the OS's platform.
//
//   - Linux/Unix: $HOME/.qiniu/config
//   - Windows: %USERPROFILE%\.qiniu\config
func SharedConfigFilename() string {
	return filepath.Join(UserHomeDir(), ".qiniu", "config")
}

// UserHomeDir returns the home directory for the user the process is
// running under.
func UserHomeDir() string {
	if runtime.GOOS == "windows" { // Windows
		return os.Getenv("USERPROFILE")
	}

	// *nix
	return os.Getenv("HOME")
}
