package version

// These variables are set at build time via -ldflags.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func String() string {
	return Version + " (" + Commit + ") built " + Date
}
