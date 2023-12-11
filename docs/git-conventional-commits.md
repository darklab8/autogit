# Git Conventional Commits

<p align="center">
  <img src="../assets/logo.png" style="width: 200px; height: 200px;"/>
</p>

## Intro

Software development is often seen as a place for geeks without social skills, controversary software development requires often communication skills to gather requirements and explain what we did and why. We write code not just for machines, but for other humans to read (including by ourselves a year later)

Messaging apps change, and the history of messages is rarely preserved. [Git]([https://www.oreilly.com/library/view/head-first-git/9781492092506/](https://www.oreilly.com/library/view/head-first-git/9781492092506/)) always remains to save another commit of code change and the message attached to it. A Git repository is the ultimate source of truth regarding a project. Every cloned git repository is a full decentralized backup of it.

We can traverse for [example](https://github.com/torvalds/linux/commit/2099306c4e1d5d772b150aeac68fdd1d0331b09d) Linux repository to see... git commits used as equivalent to emailing.

It is hard writing good atomic commit messages though. It takes some getting used to operating [git with best practices](https://deepsource.com/blog/git-best-practices)

But how we can ensure that every member will be adhering to such practices? Or even how to remember to do it yourself?
The answer is, we can enforce it with [linter](../) 😄

## Standard

[This linter - autogit](../) will help enforcing [Git conventional commits v1.0.0](https://www.conventionalcommits.org/en/v1.0.0/) standard

It operates through [git hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks) of your repository. Once linked, it will enforce itself on git commit at any your Git GUI tool.

you will be encouraged writing commits like that (the tool will prevent commit being submitted unless it adheres to the set rules)

```
<type>(<optional scope>)(optional breaking change '!' char): <subject>
empty separator line
<optional body>
empty separator line
<optional footer key: footer key value>
<optional footer2 key: footer2 key value>
```


### Example of usage

```
docs: correct spelling of CHANGELOG
```

```
feat(api)!: send an email to the customer when a product is shipped
```

```
fix: prevent racing of requests

Introduce a request id and a reference to latest request. Dismiss
incoming responses other than from latest request.

Remove timeouts which were used to mitigate the racing issue but are
obsolete now.

Reviewed-by: Z
Refs: #123
```



## How it benefits you?

- it is easier to review your and others Pull requests of submitted code.
  - each code line is easier understood why was added
- u can utilize `git blame` onto your git files like that
  for each code line u will see author of the line and exact commit why it was added
  - that will help you investigating curious bugs and problems later

```
$ git blame main.go
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  1) /*
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  2) Copyright © 2022 NAME HERE <EMAIL ADDRESS>
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  3) */
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  4) package main
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  5)
d4745237 (dd84ai 2023-12-10 00:09:35 +0100  6) import (
d4745237 (dd84ai 2023-12-10 00:09:35 +0100  7)  "autogit/interface_cli"
d4745237 (dd84ai 2023-12-10 00:09:35 +0100  8) )
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100  9)
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100 10) func main() {
def7cc5c (dd84ai 2023-10-21 23:30:37 +0200 11)  interface_cli.Execute()
^3e43b5c (dd84ai 2022-12-10 19:54:42 +0100 12) }
```

- u can use tools made for git conventional commits standard parsing,
  that will generate changelogs of changes for product releases!
  - This feature is available in [autogit](../)

example of changelog:
![changelog example](../assets/changelog_example.png)

- Your repository will look professionally with neat git commits, git tags, releases and changelogs! 😎

How can we enforce it for all developers of the repository though?

- By using [CI run with autogit](../.github/workflows/validate.yml) for every commit push or pull request!

## Tools

There are many [tools for conventional commits](https://www.conventionalcommits.org/en/about/)

What makes this tool different?

- It was written with [CI usage in mind](../.github/workflows/validate.yml)
  - it is easy to insert it into any other CI instrument with such an example.
  - just ensure you clone a repository with the history of your tags and commits 😉
- as a Golang compiled binary it does not require you to install node.js of a specific version or any other heavy interpreter for its running in CI or at the local dev machine.
  - even having `git` installed is not required for this tool to operate.
- as a Golang written tool it has a good capacity to add a lot of extra features 😃
- it is written to be as automated as possible and operate as a CLI tool helping to release products more professionally.
- it can operate with all defaults out of the box, and at the same time it can be customized with `autogit.yml` config of settings to specific needs.
- with operating through git hooks, it will work for you for any Git GUI tool and any IDE.

## Recap

- We learned that communicating our work to other developers (and yourself in a year) is important.
- Git serves ultimate source of truth for that.
- For Git commits to be a more readable history, we can utilize `git conventional commits` enforced by linters.
- We can use the history of `git conventional commits` to generate automatically changelogs of changes for releases.
