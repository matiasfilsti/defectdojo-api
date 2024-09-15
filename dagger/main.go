package main

import (
	"context"
	"dagger/defectdojo-api/internal/dagger"
	"fmt"
	"math"
	"math/rand"
)

type DefectdojoApi struct{}

func (m *DefectdojoApi) Publish(ctx context.Context, source *dagger.Directory) *dagger.Container {
	builder := dag.Container().
		From("golang:1.22.5").
		WithDirectory("/src", source).
		WithWorkdir("/src/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", "linux").
		WithExec([]string{"go", "build", "-o", "../bin/main"})

	prodImage := dag.Container().
		From("golang:1.22.5-alpine3.20").
		WithFile("/go/bin/main", builder.File("/src/bin/main")).
		WithWorkdir("/go/bin").
		WithExec([]string{"adduser", "--disabled-password", "--gecos", "--quiet", "--shell", "/bin/bash", "--u", "1000", "noonroot"}).
		WithExec([]string{"chown", "-R", "1000:1000", "/go"}).
		WithEntrypoint([]string{"main"})

	return prodImage

}

func (m *DefectdojoApi) Vulnerability(ctx context.Context, source *dagger.Directory) (*dagger.File, error) {

	containerTag := m.Publish(ctx, source)
	containerPublish, err := containerTag.Publish(ctx, fmt.Sprintf("filstimatias/dojoapi:%.0f", math.Floor(rand.Float64()*100)))
	if err != nil {
		return nil, err
	}

	vuln := dag.Trivy().Base().WithExec([]string{"trivy", "image", containerPublish, "-o", "/trivy-report-test.json", "-f", "json"})

	return vuln.File("/trivy-report-test.json"), nil

}

func (m *DefectdojoApi) TestAll(ctx context.Context, source *dagger.Directory) (string, error) {
	result, err := m.Lint(ctx, source)
	if err != nil {
		return "", err
	}

	return result, nil
}

// Returns a container that echoes whatever string argument is provided
func (m *DefectdojoApi) Test(ctx context.Context, source *dagger.Directory) *dagger.Container {
	result := m.BuildEnv(source).
		WithExec([]string{"go", "test", "./...", "-v"}).
		WithExec([]string{"go", "mod", "verify"}).
		WithExec([]string{"go", "mod", "download"}).
		WithExec([]string{"go", "build", "-v", "./..."})
	return result
}

func (m *DefectdojoApi) Lint(ctx context.Context, source *dagger.Directory) (string, error) {
	return m.Test(ctx, source).
		WithExec([]string{"go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.2"}).
		WithExec([]string{"pwd"}).
		WithExec([]string{"ls", "-la"}).
		WithExec([]string{"golangci-lint", "run", "-v", "./src/...", "./modules/...", "--issues-exit-code=1"}).
		Stdout(ctx)
}

func (m *DefectdojoApi) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1.22.5").
		WithDirectory("/src", source).
		WithWorkdir("/src")

}
