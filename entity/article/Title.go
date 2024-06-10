package article

type Title string

func (t Title) String() string {
	return string(t)
}
