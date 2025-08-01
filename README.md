# vfafgo – File, Version, Encryption, HTTP, and Zip Utilities

This package provides primitives for file storage configuration, ignore pattern matching, virtual filesystem (VFS), lightweight git-style versioning, encrypted streaming, HTTP PUT helpers, and directory compression/extraction with zip. Below, every function/struct in the codebase is explained, with usage for the included Go tests.

---

## Function and Struct Explanations

### config.go

**Config struct** – Configuration for storage location, local directory, and remote URL.

- `LoadStorageConfig()` – Reads the root directory for storage from the environment or uses a default.
- `New(localDir, remoteURL string)` – Creates a `Config` struct, ensuring the directory exists and is absolute.

---

### ignore.go

**IgnoreMatcher struct** – Holds ignore/globbing patterns from a `.qignore` file.

- `LoadIgnoreMatcher(baseDir string)` – Loads ignore patterns found in `.qignore` within a directory.
- `(IgnoreMatcher) Ignore(relPath string)` – Checks whether a given path matches any of the ignore patterns.

---

### vfs.go

**Vfs struct** – Virtual filesystem abstraction with root directory.

- `NewFS(root string)` – Initializes a new virtual filesystem at the given root.
- `(*Vfs) Path(parts ...string)` – Returns full path for sub-paths combined under root.
- `(*Vfs) Read(parts ...string)` – Returns file contents for a relative path.
- `(*Vfs) Write(content string, parts ...string)` – Writes content to a relative file, creating directories as necessary.
- `(*Vfs) Exists(parts ...string)` – Returns true if the file or directory exists under the root.

---

### vcr.go

**Vcr struct** – Wrapper around a git repository for simple versioning.

- `NewVcr(user, workdir string)` – Opens/initializes a git repo in the specified workdir. Initialization only occurs if user != "default".
- `(*Vcr) Commit(message string)` – Commits *all* staged changes and returns the commit hash.
- `(*Vcr) BranchFrom(baseTag, comment string)` – Creates and checks out a new branch from tag or HEAD, returns new commit id.
- `(*Vcr) GetHistory()` – Lists all commits in the repository.
- `(*Vcr) Checkout(refName string)` – Switches repository to a branch or tag.
- `(*Vcr) RewriteToMain(sourceRef, message string)` – Overwrites main branch with the content from another branch or commit.
- `(*Vcr) Purge()` – Deletes the internal .git directory, erasing all repo history.

---

### http.go

- `encrypt(buf bytes.Buffer, key []byte)` – Encrypts data using the external `cryptogo` library and a provided key.
- `PutStreamEncrypted(remoteURL, rel, buf, key)` – Encrypts and uploads a buffer to a remote HTTP server.
- `PutStream(remoteURL, rel, buf)` – Uploads (PUTs) the buffer directly to an HTTP server.

---

### zip.go

- `Unzip(src, dest)` – Extracts all files from a zip archive into a destination folder.
- `Zip(localRoot, relInputPath, dest)` – Creates a zip archive of a folder, ignoring files listed in `.qignore` if present; result is stored to buffer and can be saved to disk.

---

## How to Use the Tests in `vfafgo_test.go`

The file `vfafgo_test.go` includes tests for all main features:

- **Virtual File System (TestVfsBasic):** Checks file creation, reading, directory handling.
- **VCR/Version Control (TestVcrBasic):** Checks git-like versioning, commits, purging version history (.git removal).
- **Config (TestConfig):** Checks config and parameter parsing.
- **Ignore Patterns (TestIgnore):** Checks pattern loading and file ignoring.
- **HTTP Helpers (TestHTTPHelpers):** Checks error handling and encrypt function (actual network calls are error-expected; focus is on logic).
- **Zip/Unzip (TestZipUnzip):** Verifies directory zipping, file ignore support, and zip extraction.

### Running the Tests

From the package directory, run:

```sh
go test -v
```

You’ll see verbose output; errors/failures will be reported if any function does not behave as expected. All tests are self-contained. No network or external setup is required unless you add new HTTP endpoints.

---

## License

This project is licensed under the [European Union Public Licence v1.2 (EUPL)](https://joinup.ec.europa.eu/collection/eupl/eupl-text-eupl-12).
