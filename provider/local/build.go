package local

import (
	"fmt"
	"io"
	"time"

	"github.com/convox/praxis/types"
)

func (p *Provider) BuildCreate(app, url string, opts types.BuildCreateOptions) (*types.Build, error) {
	_, err := p.AppGet(app)
	if err != nil {
		return nil, err
	}

	id := types.Id("B", 10)

	build := &types.Build{
		Id:      id,
		App:     app,
		Status:  "created",
		Created: time.Now(),
	}

	if err := p.Store(fmt.Sprintf("apps/%s/builds/%s", app, id), build); err != nil {
		return nil, err
	}

	pid, err := p.ProcessStart(app, types.ProcessStartOptions{
		Command: fmt.Sprintf("build -id %s -url %s", id, url),
		Environment: map[string]string{
			"BUILD_APP":  app,
			"BUILD_PUSH": "foo",
		},
		Name:    fmt.Sprintf("%s-build-%s", app, id),
		Image:   "convox/praxis:test7",
		Service: "build",
		Volumes: map[string]string{
			"/var/run/docker.sock": "/var/run/docker.sock",
		},
	})
	if err != nil {
		return nil, err
	}

	build.Process = pid

	if err := p.Store(fmt.Sprintf("apps/%s/builds/%s", app, id), build); err != nil {
		return nil, err
	}

	return build, nil
}

func (p *Provider) BuildGet(app, id string) (build *types.Build, err error) {
	err = p.Load(fmt.Sprintf("apps/%s/builds/%s", app, id), &build)
	return
}

func (p *Provider) BuildList(app string) (types.Builds, error) {
	return nil, nil
}

func (p *Provider) BuildLogs(app, id string) (io.ReadCloser, error) {
	build, err := p.BuildGet(app, id)
	if err != nil {
		return nil, err
	}

	switch build.Status {
	case "running":
		return p.ProcessLogs(app, build.Process)
	default:
		return p.ObjectFetch(app, fmt.Sprintf("convox/builds/%s/log", id))
	}
}

func (p *Provider) BuildUpdate(app, id string, opts types.BuildUpdateOptions) (*types.Build, error) {
	build, err := p.BuildGet(app, id)
	if err != nil {
		return nil, err
	}

	if !opts.Ended.IsZero() {
		build.Ended = opts.Ended
	}

	if opts.Manifest != "" {
		build.Manifest = opts.Manifest
	}

	if opts.Release != "" {
		build.Release = opts.Release
	}

	if !opts.Started.IsZero() {
		build.Started = opts.Started
	}

	if opts.Status != "" {
		build.Status = opts.Status
	}

	if err := p.Store(fmt.Sprintf("apps/%s/builds/%s", app, id), build); err != nil {
		return nil, err
	}

	return build, nil
}
