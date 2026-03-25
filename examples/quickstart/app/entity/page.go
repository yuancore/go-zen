package entity

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func NormalizePage(page, size int) (int, int) {
	if page < DefaultPage {
		page = DefaultPage
	}
	switch {
	case size <= 0:
		size = DefaultPageSize
	case size > MaxPageSize:
		size = MaxPageSize
	}
	return page, size
}
