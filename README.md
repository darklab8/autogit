<p align="center">
  <img src="assets/logo.png" />
</p>

## Description

Solution for

- git conventional commits validations
- automatic semantic versioning generation
- auto generating changelogs

## Supported OS and Architectures:

- Linux amd64
- Linux arm64
- Windows amd64

... U a free to help adding more supported versions, by finding gcc compiler from linux to OS Y, Arch Z

## Installation (Draft)

1) u install binary somewhere accessable to your OS bin path
2) adding commit-msg file to your .git folder lets say

## Usages

#### scenario #1 - validator / Git commit validation

You try to write git commit -m "feat: bla bla bla"
your githook is activated and tries to parse your commit name accroding to git conventional commits standard. If unable, it will give you error and prevent commit

#### scenario #2 - changelog / Your wish to see changelog of additions you made, what are new features, what are fixes. For user view

You wish to have changelog auto generated.
program parses your conventional commits, and renders output in markdown for copy paste to github
(Changelog is generated since the last release version / last semantic tag applied to repository)
P.S. you can also request at any time changelogs from previous releases/tags. Everything is parsed from your commits

#### Scenario #3 - nextSemVer / You wish to know which next semantic version / semantic tag should be applied to your release.

Program checks if u made no commits, or only refactoring and styling. Then it says, next version is same as previous one. Nothing changed, for example `0.0.1` as first one.
If you made `fix`, then it increases PATCH version of semantic version. Your next version is `0.0.2` rendered
if you made `feat` request, then next MINOR version is increased. Your next version is `0.1.0`
if you made breaking changes, users should know `feat!` or `BREAKING CHANGE:`, then next version is `1.0.0`

#### TLDR

So in a nutshell, it takes away complexity of using git conventional commits and semantic versioning. You are auto guided and auto corrected how correctly to perform it xD
u only need correctly writing meaning/subject/description to your commits 🙂 but since u see what is rendered to users, you quickly learn how to write it better

why semantic versioning is important, to read here https://semver.org/
well, about git conventional commits is here: https://www.npmjs.com/package/git-conventional-commits

as an example, all my releases of darktool were made with similar automatation.
Changelogs and versions https://github.com/darklab8/darklab_freelancer_darktool/releases

## Future development and resources

- https://www.quora.com/What-is-the-difference-between-alpha-beta-and-RC-software-version // Adding ability of best versions
- https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional // Configurable stricter rules to validator
- https://github.com/angular/angular/blob/22b96b9/CONTRIBUTING.md#-commit-message-guidelines // Just more about git conventional commits

## Dev Requirements

- cobra generator https://github.com/spf13/cobra-cli/blob/main/README.md
- cobra guide https://github.com/spf13/cobra/blob/main/user_guide.md
- godoc
- add binary discovery for cobra-cli, godoc detection
  - `export PATH="$PATH:/usr/local/go/bin:$HOME/go/bin"`
- Git hooks of conventional commits
  - [docs](https://gist.github.com/qoomon/5dfcdf8eec66a051ecd85625518cfd13)
  - [app](https://www.npmjs.com/package/git-conventional-commits)
