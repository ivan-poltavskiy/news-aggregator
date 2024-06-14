package article

type Description string

func (d Description) String() string {
	return string(d)
}
