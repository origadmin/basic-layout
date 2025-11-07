# TODO: Rename `basic-layout` to `starter-kits`

This document outlines the necessary steps to rename the `basic-layout` directory to `starter-kits`. This change better reflects its role as a high-quality, reusable project template library.

**Note:** This is a global change. All steps should be performed carefully to ensure project integrity.

---

### Action Plan

- [ ] **1. Manually Rename the Directory**
  - Rename the directory `examples/basic-layout` to `examples/starter-kits`.

- [ ] **2. Update Go Module Paths (`go.mod`)**
  - **File**: `examples/starter-kits/simple/simple_app/go.mod`
    - **From**: `module basic-layout/simple/simple_app`
    - **To**: `module starter-kits/simple/simple_app`
  - **File**: `examples/starter-kits/multiple/multiple_sample/go.mod`
    - **From**: `module basic-layout/multiple/multiple_sample`
    - **To**: `module starter-kits/multiple/multiple_sample`

- [ ] **3. Update Go Source Code References**
  - **Scope**: All `.go` files within `examples/starter-kits/multiple/multiple_sample/`.
  - **Action**: Perform a global search and replace.
    - **Find**: `basic-layout/multiple/multiple_sample`
    - **Replace with**: `starter-kits/multiple/multiple_sample`
  - **Key Target**: Pay special attention to `go:generate` directives, for example, in `internal/features/user/data/data.go`.

- [ ] **4. Update Documentation Files (`.md`)**
  - **Scope**: All `.md` files within `examples/starter-kits/`.
  - **Action**: Perform a global search and replace.
    - **Find**: `basic-layout`
    - **Replace with**: `starter-kits`
  - **Files to check include**:
    - `README.md`
    - `docs/how-to-use-template.md`
    - All other files under the `docs/` directory.

- [ ] **5. Update Build & Configuration Files**
  - **Scope**: `examples/starter-kits/multiple/multiple_sample/`
  - **Action**: Perform a search and replace in the following files:
    - `.goreleaser.yaml`
    - `buf.yaml`
  - **Find**: `basic-layout/multiple/multiple_sample`
  - **Replace with**: `starter-kits/multiple/multiple_sample`
