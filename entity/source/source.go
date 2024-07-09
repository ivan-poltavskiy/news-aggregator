package source

// PathToFile describes the path to a specific source in the system
type PathToFile string

// Name describes the name of source
type Name string

// Type describes of the type of document for a particular source
type Type string

// Source is the set of information about source of news.
type Source struct {
	Name       Name
	PathToFile PathToFile
	SourceType Type
}

// Stores all types of sources provided.
const (
	RSS      Type = "RSS"
	JSON     Type = "JSON"
	UsaToday Type = "UsaToday"
)

// NewsSources stores name of resources from which news can be obtained.
var NewsSources = []string{
	"abc",
	"bbc",
	"nbc",
	"usatoday",
	"washington",
}
