# Autogit

<p align="center">
  <img src="assets/logo.png" style="width: 200px; height: 200px;"/>
</p>

**Communicating through git professionally**

autogit is a CLI tool to validate submitted commits according to [git conventional commits standard](https://www.conventionalcommits.org/en/v1.0.0/). the tool allows to generate changelogs for releases in different formats. When you generate changelogs and see quick feedback of an end result, it promotes you to write commits better.

as a result of a tool work, you communicate your developer work better to other developers and have more professional looking repository.

# Features

- hooks to git-hooks and works to validate your git commits to git conventional commits standard for any git tool.
  - has extra possible validating rules to configure, like having minimum 3 words in a subject of a commit.
  - `autogit hook activate --global`(flag to turn it on for all repos)
- suggests next [semantic version](https://semver.org/) with `autogit semver`
  - has options to sugest next version as alpha, beta, prerelease version with build meta data.
- generates changelogs with `autogit changelog` command
  - currently supports markdown and bbcode formats
  - has option `--validate` to run validation of commits (for CI usage)
- easy create and push of a git tag with autoinserted changelog through `autogit semver --tag --push`
- CI friendly, not requires any dependencies for its usage for everything (inbuilt git-go to access git information)

# Example

text version at ubuntu 22.04:
- `apt update && apt install -y curl` (install curl if not installed)
- install latest

# Installation

## Linux

### Install latest

- install curl if not present
- install with `rm /usr/local/bin/autogit ; curl -L $(curl -Ls -o /dev/null -w %{url_effective} https://github.com/darklab8/darklab_autogit/releases/latest | sed "s/releases\/tag/releases\/download/")/autogit-linux-amd64 -o /usr/local/bin/autogit && chmod 777 /usr/local/bin/autogit`

- check installation with `autogit version` command. Expect to see `OK autogit version: v{version}`

### install specific version

- install with `rm /usr/local/bin/autogit ; curl -L https://github.com/darklab8/darklab_autogit/releases/download/v2.1.0/autogit-linux-amd64 -o /usr/local/bin/autogit && chmod 777 /usr/local/bin/autogit`
- See other installations at {INSERT LINK TO OTHER INSTALLATIONS}

## Contacts

- [@dd84ai](https://github.com/dd84ai) at `dark.dreamflyer@gmail.com`
- open [Pull Requests with question](https://github.com/darklab8/darklab_autogit/issues)
- [Darklab Discord server](https://discord.gg/aukHmTK82J)
