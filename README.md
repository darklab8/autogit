# Autogit

<p align="center">
  <img src="assets/logo.png" style="width: 200px; height: 200px;"/>
</p>

**Communicating through git professionally**

autogit is a CLI tool to validate submitted commits according to [git conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) standard. the tool allows to generate changelogs for releases in different formats. When you generate changelogs and see quick feedback of an end result, it promotes you to write commits better.

as a result of a tool work, you communicate your developer work better to other developers and have more professional looking repository.

# Features

- hooks to git-hooks and works to validate your git commits to [git conventional commits]((https://www.conventionalcommits.org/en/v1.0.0/)) standard for any git tool.
  - has extra possible validating rules to configure, like having minimum 3 words in a subject of a commit.
  - `autogit hook activate --global`(flag to turn it on for all repos)
- suggests next [semantic version](https://semver.org/) with `autogit semver`
  - has options to sugest next version as alpha, beta, prerelease version with build meta data.
- generates changelogs with `autogit changelog` command
  - currently supports markdown and bbcode formats
  - has option `--validate` to run validation of commits (for CI usage)
- easy create and push of a git tag with autoinserted changelog through `autogit semver --tag --push`
- initialize settings for more customization with `autogit init` inside git repo
  - uncomment and override desired settings
- find out more commands and options with `autogit [any set of sub commands] --help`
- CI friendly, not requires any dependencies for its usage for everything (inbuilt git-go to access git information)

# Example

text version at ubuntu 22.04:
- `apt update && apt install -y curl` (install curl if not installed)
- [install latest](#install-latest)
- install git if not present with `apt install -y git`
- init git repo if not present `git init && git config user.email "you@example.com" && git config user.name example`
- activating git hook, `autogit hook activate --global` (optionally global for all repos)
- write some commits:
  - `echo 123 >> README.md && git add -A && git commit -m "feat: init repo with first code"`
  - `echo 123 >> README.md && git add -A && git commit -m "fix: memory leak in sql connection opener"`
  - `echo 123 >> README.md && git add -A && git commit -m "feat: new super feature"`
  - `echo 123 >> README.md && git add -A && git commit -m 'feat!: new super feature`

`BREAKING CHANGE: api for endpoint status changed to users-status'`
    - due to bash using `!` as a keyword syntax, we need to use `''` single quotes
  - `echo 123 >> README.md && git add -A && git commit -m "fix(api): example of scoped bug fix"`
  - generate changelog with `autogit changelog`

![changelog example](assets/changelog_example.png)

# Installation

## Linux

### Install latest

- install curl if not present (`apt update && apt install -y git` for debian/ubuntu)
- install with `rm /usr/local/bin/autogit ; curl -L $(curl -Ls -o /dev/null -w %{url_effective} https://github.com/darklab8/darklab_autogit/releases/latest | sed "s/releases\/tag/releases\/download/")/autogit-linux-amd64 -o /usr/local/bin/autogit && chmod 777 /usr/local/bin/autogit`

- check installation with `autogit version` command. Expect to see `OK autogit version: v{version}`

### install specific version

- install with `rm /usr/local/bin/autogit ; curl -L https://github.com/darklab8/darklab_autogit/releases/download/v2.1.0/autogit-linux-amd64 -o /usr/local/bin/autogit && chmod 777 /usr/local/bin/autogit`

### See other installations

- at {INSERT LINK TO OTHER INSTALLATIONS}

## Contacts

- [@dd84ai](https://github.com/dd84ai) at `dark.dreamflyer@gmail.com`
- open [Pull Requests with question](https://github.com/darklab8/darklab_autogit/issues)
- [Darklab Discord server](https://discord.gg/aukHmTK82J)
