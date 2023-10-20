package logus

import (
	"autogit/settings/types"
	"fmt"
	"log/slog"
)

func logGroupFiles() slog.Attr {
	return slog.Group("files",
		"file3", GetCallingFile(3),
		"file4", GetCallingFile(4),
	)
}

type slogGroup struct {
	params map[string]string
}

func (s slogGroup) Render() slog.Attr {
	anies := []any{}
	for key, value := range s.params {
		anies = append(anies, key)
		anies = append(anies, value)
	}

	return slog.Group("extras", anies...)
}

type slogParam func(r *slogGroup)

func newSlogGroup(opts ...slogParam) slog.Attr {
	client := &slogGroup{params: make(map[string]string)}
	for _, opt := range opts {
		opt(client)
	}

	return (*client).Render()
}

func TestParam(value int) slogParam {
	return func(c *slogGroup) {
		c.params["test_param"] = fmt.Sprintf("%d", value)
	}
}

func ConfigPath(value types.ConfigPath) slogParam {
	return func(c *slogGroup) {
		c.params["config_path"] = string(value)
	}
}

func Expected(value any) slogParam {
	return func(c *slogGroup) {
		c.params["expected"] = fmt.Sprintf("%v", value)
	}
}
func Actual(value any) slogParam {
	return func(c *slogGroup) {
		c.params["actual"] = fmt.Sprintf("%v", value)
	}
}

// func WithAllowFailedJobs[T IJob](value bool) JobPoolOption[T] {
// 	return func(c *JobPool[T]) {
// 		c.allow_failed_jobs = value
// 	}
// }

// func NewJobPool[T IJob](opts ...JobPoolOption[T]) JobPool[T] {
// 	client := &JobPool[T]{}
// 	for _, opt := range opts {
// 		opt(client)
// 	}

// 	return *client
// }

// jobPool := NewJobPool[*JobTest](
// 	WithAllowFailedJobs[*JobTest](true),
// )
