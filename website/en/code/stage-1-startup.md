# Stage 1: Startup Full Code

Stage 1 creates the smallest Go project: module, CLI, domain function, and tests.

This is the learning version of the complete code. The comments are intentionally denser than production code so a new Go learner can read each file without guessing the role of each part.

## Files

```text
hn-agent/go.mod
hn-agent/cmd/hnctl/main.go
hn-agent/internal/hn/story.go
hn-agent/internal/hn/story_test.go
```

## `go.mod`

```go
// hn-agent/go.mod
// module declares the identity of this Go module.
// Internal imports in this project start from this module path.
module github.com/aaron/go-with-ai/hn-agent

// go declares the Go language version targeted by this module.
// The learning notes use Go 1.27 as the baseline.
go 1.27
```

## CLI Entry

```go
// hn-agent/cmd/hnctl/main.go
// package main means this package builds into an executable program.
package main

import (
	// fmt prints formatted text.
	"fmt"
	// os exposes operating-system features such as args and exit codes.
	"os"

	// This is an internal project package.
	// The import path comes from go.mod plus the directory path.
	"github.com/aaron/go-with-ai/hn-agent/internal/hn"
)

// const defines a value that does not change at runtime.
const version = "0.1.0"

// run holds the CLI logic separately from main.
// That makes the logic easier to test later because tests do not need to call os.Exit.
func run(args []string) int {
	// os.Args[1:] may be empty, so check len(args) before reading args[0].
	if len(args) > 0 && (args[0] == "version" || args[0] == "-v") {
		fmt.Println("hnctl version is", version)
		return 0
	}

	// Create a Story value with named fields.
	story := hn.Story{
		ID:    1,
		Title: "Hello",
		URL:   "https://www.example.com",
		Score: 0,
	}

	fmt.Printf("hnctl %s\n", version)
	fmt.Printf("story valid: %v\n", hn.IsValidStory(story))
	return 0
}

// main is the entry point for a Go executable.
func main() {
	// os.Args[0] is the executable name.
	// os.Args[1:] is the real CLI argument list.
	os.Exit(run(os.Args[1:]))
}
```

## Story Model

```go
// hn-agent/internal/hn/story.go
// package hn means this file belongs to the hn package.
package hn

// Story is the smallest domain model for a Hacker News story in Stage 1.
// A struct groups named fields together.
type Story struct {
	ID    int
	Title string
	URL   string
	Score int
}

// IsValidStory is a small domain function.
// It is easy to test because it does not depend on network, files, or CLI state.
func IsValidStory(story Story) bool {
	// Stage 1 uses a minimal rule:
	// a story is valid when it has a positive ID and a non-empty title.
	return story.ID > 0 && story.Title != ""
}
```

## Story Test

```go
// hn-agent/internal/hn/story_test.go
package hn

// testing is Go's standard test package.
import "testing"

// TestIsValidStory tests IsValidStory.
// go test automatically runs functions whose names start with Test.
func TestIsValidStory(t *testing.T) {
	// This is a table-driven test.
	// []struct means a slice of anonymous structs.
	tests := []struct {
		name string
		story Story
		want bool
	}{
		{
			name: "valid story",
			story: Story{
				ID:    1,
				Title: "a valid story",
			},
			want: true,
		},
		{
			name: "missing id and title",
			story: Story{
				ID:    0,
				Title: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		// t.Run creates a named subtest.
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidStory(tt.story)
			if got != tt.want {
				// Fatalf fails this subtest and stops it immediately.
				t.Fatalf("IsValidStory() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## Acceptance

```bash
cd hn-agent
go test ./internal/hn
go run ./cmd/hnctl
go run ./cmd/hnctl version
```

