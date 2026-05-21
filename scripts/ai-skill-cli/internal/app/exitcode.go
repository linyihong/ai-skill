package app

const (
	ExitSuccess             = 0
	ExitGeneralFailure      = 1
	ExitInvalidUsage        = 2
	ExitMissingDependency   = 10
	ExitUnsupportedPlatform = 11
	ExitUnsafeRepoState     = 20
	ExitPermissionDenied    = 21
	ExitValidationFailed    = 30
	ExitGeneratedStale      = 31
	ExitPartialCloseBlocked = 40
)
