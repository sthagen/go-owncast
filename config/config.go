package config

import (
	"fmt"
	"time"
)

// These are runtime-set values used for configuration.

// DatabaseFilePath is the path to the file ot be used as the global database for this run of the application.
var DatabaseFilePath = "data/owncast.db"

// LogDirectory is the path to various log files.
var LogDirectory = "./data/logs"

// EnableDebugFeatures will print additional data to help in debugging.
var EnableDebugFeatures = false

// VersionNumber is the current version string.
var VersionNumber = StaticVersionNumber

// WebServerPort is the port for Owncast's webserver that is used for this execution of the service.
var WebServerPort = 8080

// Bind WebServer to this IP address. Be secure by default.
var WebServerIP = "0.0.0.0"

// InternalHLSListenerPort is the port for HLS writes that is used for this execution of the service.
var InternalHLSListenerPort = "8927"

// GitCommit is an optional commit this build was made from.
var GitCommit = ""

// BuildPlatform is the optional platform this release was built for.
var BuildPlatform = "dev"

func GetCommit() string {
	if GitCommit == "" {
		GitCommit = time.Now().Format("20060102")
	}

	return GitCommit
}

// GetReleaseString gets the version string.
func GetReleaseString() string {
	var versionNumber = VersionNumber
	var buildPlatform = BuildPlatform
	var gitCommit = GetCommit()

	return fmt.Sprintf("Owncast v%s-%s (%s)", versionNumber, buildPlatform, gitCommit)
}
