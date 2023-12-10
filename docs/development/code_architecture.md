# Code architecture

## Goals

- Unit testable first. everything else later.
- Abstractions will appear with a strict minimal interface to reduce overall complexity of a code.
- High usage of `type NewType string` for more self documentation
- Minimize third party lib dependencies
- Simplify end user installation
- No autoupdates. Everything should work offline.
- CI friendly, zero system dependencies solution

## Diagram

```mermaid
flowchart TD
  UI[Interface-CLI\nUser interface via Cobra CLI third party lib]
  UI --> Actions[Actions\nreusable actions without\nattachements to UI details]
  Actions --> SemanticGit[Semantic Git\nImplements main business logic of repository\nwith added logic of conventional commits\nAnd semantic versioning]
  Actions --> Changelog[Changelog\nHow to generate one]
  Actions --> Validator[Validator\nRules for optional validations]
  Changelog --> SemanticGit
  Validator --> SemanticGit
  SemanticGit --> SemVer[SemVer\nimplements original Semantic Version\naccording to SemVer2.0.0 standard\nImplemented in current repo]
  SemanticGit --> Git[Git\ngit wrapper to simple interface\nfor current repository logic\nimplemented in current repo]
  Git --> GitGo[Git-Go\nEngine under the hood for\nGit repository operations\nImplemented by third party]
```

## Support promises:

### First tier support

- for linux and CI usage
    - we will prioritize solving any issues.
    - Unit tests run on every commit.

### Second tier support

- for windows
    - periodically checking it works in VM with Windows 10

### Third tier support

- for macos
    - only compiled

