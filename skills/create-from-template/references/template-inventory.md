# Template Inventory Reference

Quick reference for all available code templates and their locations.

## Go Templates

### Executables (package main)

| Source | Path | Use For |
|--------|------|---------|
| **Root (Canonical)** | `templates/code/go/CODE-GO-001-GO-executable.go` | CLI tools, entry points, orchestrators |
| **Compiler** | `divisions/tech/language/compiler/templates/go-main-executable.go` | Same - implementation aligned |

### Libraries (reusable packages)

| Source | Path | Use For |
|--------|------|---------|
| **Root (Canonical)** | `templates/code/go/CODE-GO-002-GO-library.go` | Reusable packages, APIs, shared logic |
| **Compiler** | `divisions/tech/language/compiler/templates/go-library-package.go` | Same - implementation aligned |

### Demo-Tests (showcase/demonstration)

| Source | Path | Use For |
|--------|------|---------|
| **Root (Canonical)** | `templates/code/go/CODE-GO-003-GO-demo-test.go` | Demonstration-based tests, showcases |
| **Compiler** | `divisions/tech/language/compiler/templates/go-demo-test.go` | Same - implementation aligned |

## Go Configuration Templates

| File | Path | Use For |
|------|------|---------|
| **go.mod** | `divisions/tech/language/compiler/templates/go-module.mod` | Module configuration |
| **go.work** | `divisions/tech/language/compiler/templates/go-workspace.work` | Workspace configuration |
| **go.sum** | `divisions/tech/language/compiler/templates/go-checksum.sum` | Checksum documentation |

## Build System Templates

| System | Path | Use For |
|--------|------|---------|
| **Makefile** | `divisions/tech/language/compiler/templates/make-project-makefile.mk` | Project-level Makefiles |
| **CMakeLists.txt** | `divisions/tech/language/compiler/templates/cmake-project-cmakelists.cmake` | CMake build configuration |

## Quick Selection Guide

**What am I creating?**

```
Creating new Go file?
├── Entry point / CLI / main() → CODE-GO-001 (executable)
├── Reusable package / API → CODE-GO-002 (library)
└── Demo / showcase / RunX() → CODE-GO-003 (demo-test)

Creating Go configuration?
├── go.mod → go-module.mod template
├── go.work → go-workspace.work template
└── go.sum → go-checksum.sum template

Creating build configuration?
├── Makefile → make-project-makefile.mk
└── CMakeLists.txt → cmake-project-cmakelists.cmake
```

## Path Constants

**Workspace root:** `/media/seanje-lenox-wise/Project/CreativeWorkzStudio_LLC`

**Root templates:** `{workspace}/templates/code/go/`

**Compiler templates:** `{workspace}/divisions/tech/language/compiler/templates/`

## Reading Templates

```bash
# Read Go executable template (root)
cat templates/code/go/CODE-GO-001-GO-executable.go

# Read Go executable template (compiler)
cat divisions/tech/language/compiler/templates/go-main-executable.go

# Read go.mod template
cat divisions/tech/language/compiler/templates/go-module.mod
```

---

**Updated:** 2025-11-28
