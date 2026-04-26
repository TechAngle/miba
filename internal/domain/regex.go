package domain

import "regexp"

// list of regexes for different values
var (
	DeviceIDRegex   = regexp.MustCompile(`(?i)deviceId:\s*['"]([0-9a-fA-F-]+)['"]`)
	ROMVersionRegex = regexp.MustCompile(`(?i)romVersion:\s*['"]([\d\.]+)['"]`)
	HardwareRegex   = regexp.MustCompile(`(?i)hardwareVersion:\s*['"]([A-Za-z0-9]+)['"]`)
	KeyRegex        = regexp.MustCompile(`(?i)key:\s*['"]([a-f0-9]{32})['"]`)
	IvRegex         = regexp.MustCompile(`(?i)iv:\s*['"](\d+)['"]`)
	MacRegex        = regexp.MustCompile(`(?i)deviceId\s*=\s*['"](([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2})['"]`)
)
