package source

// Source is the set of information about source of article.
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
	"ABC",
	"BBC",
	"NBC",
	"USA Today",
	"Washington Times",
}
