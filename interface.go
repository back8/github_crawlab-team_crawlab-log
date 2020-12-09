package log

type Driver interface {
	Init() error
	Close() error
	Write(line string) error
	WriteLines(lines []string) error
	Find(pattern string, skip int, limit int) ([]string, error)
	Count(pattern string) (int, error)
}