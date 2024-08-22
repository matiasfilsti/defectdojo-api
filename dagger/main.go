// A generated module for DefectdojoApi functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/defectdojo-api/internal/dagger"
	"fmt"
	"math"
	"math/rand"
)

type DefectdojoApi struct{}

// // Returns a container that echoes whatever string argument is provided
// func (m *DefectdojoApi) ContainerEcho(stringArg string) *dagger.Container {
// 	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
// }

// // Returns lines that match a pattern in the files of the provided Directory
// func (m *DefectdojoApi) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
// 	return dag.Container().
// 		From("alpine:latest").
// 		WithMountedDirectory("/mnt", directoryArg).
// 		WithWorkdir("/mnt").
// 		WithExec([]string{"grep", "-R", pattern, "."}).
// 		Stdout(ctx)
// }

func (m *DefectdojoApi) Publish(ctx context.Context, source *dagger.Directory) (string, error) {
	builder := dag.Container().
		From("golang:1.22.1").
		WithDirectory("/src", source).
		WithWorkdir("/src/src").
		WithEnvVariable("CGO_ENABLED", "0").
		WithEnvVariable("GOOS", "linux").
		WithExec([]string{"go", "build", "-o", "../bin/main"})

	prodImage := dag.Container().
		From("golang:1.22.1-alpine3.19").
		WithFile("/go/bin/main", builder.File("/src/bin/main")).
		WithWorkdir("/go/bin").
		WithExec([]string{"adduser", "--disabled-password", "--gecos", "--quiet", "--shell", "/bin/bash", "--u", "1000", "noonroot"}).
		WithExec([]string{"chown", "-R", "1000:1000", "/go"}).
		WithEntrypoint([]string{"main"})

	address, err := prodImage.Publish(ctx, fmt.Sprintf("filstimatias/ml-challenge:%.0f", math.Floor(rand.Float64()*100)))
	if err != nil {
		return "", err
	}
	return address, nil
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

// Build a ready-to-use development environment
func (m *DefectdojoApi) BuildEnv(source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:1.22.5").
		WithDirectory("/src", source).
		WithWorkdir("/src")

}

func (m *DefectdojoApi) HttpService(ctx context.Context, source *dagger.Directory) *dagger.Service {
	return m.BuildEnv(source).
		WithExec([]string{"go", "run", "src/main.go"}).
		WithExposedPort(8080).
		AsService()
}

// Send a request to an HTTP service and return the response
func (m *DefectdojoApi) Get(ctx context.Context, source *dagger.Directory) (string, error) {
	return dag.Container().
		From("alpine").
		WithServiceBinding("www", m.HttpService(ctx, source)).
		WithExec([]string{"wget", "-O-", "http://www:8080/categories/MLA97994"}).
		Stdout(ctx)
}
