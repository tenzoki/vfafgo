
# govifilemod: Lightweight File & Version Control Utilities for Go

## Description

The `govifilemod` module provides two main components:

- **Vcr**: Simple version control operations using `go-git`, allowing you to version your workspace or individual project directories without external git commands.
- **Vfs**: Virtual file system utilities for easy file read/write/exists operations relative to a root directory.

Both types are lightweight and designed for integration into Go projects needing basic file and version management.

---

## Public Types and Functions

### Vfs (Virtual File System)

#### Type
```go
type Vfs struct {
    root string
}
```

#### Public Functions

- **NewFS(root string) *Vfs**
  Initializes a new Vfs rooted at the given directory path.

- **(fs *Vfs) Path(parts ...string) string**
  Returns the absolute path for the given path segments relative to the root.

- **(fs *Vfs) Read(parts ...string) (string, error)**
  Reads a file (relative to root) and returns its contents as a string.

- **(fs *Vfs) Write(content string, parts ...string) error**
  Writes the given content to a file (relative to root), creating directories as needed.

- **(fs *Vfs) Exists(parts ...string) bool**
  Returns true if the file or directory (relative to root) exists.

### Vcr (Version Control Recorder)

#### Type
```go
type Vcr struct {
    workdir string
    repo    *git.Repository
}
```

#### Public Functions

- **NewVcr(user, workdir string) *Vcr**
  Initializes or opens a git repository in the specified working directory. The repository is initialized only if `user != "default"`.
- **(v *Vcr) Commit(message string) string**
  Create a commit with a message and return the commit id.
- **(v *Vcr) BranchFrom(baseTag, comment string) string**
  Create a new branch from a tag or HEAD, with a comment.
- **(v *Vcr) GetHistory() []string**
  List all commits by time (as strings).
- **(v *Vcr) Checkout(refName string) error**
  Switch to a given branch or tag.
- **(v *Vcr) RewriteToMain(sourceRef, message string) error**
  Replace main branch with content from another ref and commit with a message.
- **(v *Vcr) Purge() error**
  Remove the `.git` directory to reset versioning.

---

## Usage in Another Go Project

### 1. Install Dependencies

Make sure to include this module and its dependencies—`github.com/go-git/go-git/v5` is required.

### 2. Using the Module

#### Importing

If your module name is, for example, `github.com/tenzoki/govifilemod`, import it as:

```go
import "github.com/tenzoki/govifilemod"
```
Change the path accordingly.

#### Example

```go
package main

import (
    "fmt"
    "github.com/tenzoki/govifilemod"
)

func main() {
    // VFS Usage
    fs := gov.NewFS("./workspace")
    _ = fs.Write("hello world", "test.txt")
    data, err := fs.Read("test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println("test.txt contains:", data)

    // VCR Usage
    vcr := gov.NewVcr("john", "./workspace")
    vcr.Commit("Initial save")
    // List commit history
    fmt.Println("History:", vcr.GetHistory())
}
```
## Tests

Run tests:

```sh
go test -v
```
