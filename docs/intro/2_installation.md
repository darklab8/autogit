# Installation

## 🚀 Short opinionated installation

2) `autogit semver` to verify installation and get next expected semantic version
3) (Optional) config init

   * `autogit init` to create `autogit.yml` config locally in repo.
   * `autogit init --global` to make user global settings file fallback
   * if u will do nothing, program will just fallback to using config from memory
4) `autogit hook activate` to turn validation commit hooks on

   * `autogit hook activate --global` will turn it on globally for all repos
   * `autogit hook deactivate` (also with possible flag `--global`) can serve to deactivate this functionality

## ✈️ Detailed installation

1. [download latest stable release](https://github.com/darklab8/darklab_autogit/releases) and put to env PATH searchable range

- Linux:

  - Check your PATH bin serachable locations with `echo "$PATH"` and put into any of them or add new location, change settings to allow it being executable with `chmod`
  - Recomendation to put into `/usr/local/bin`
  - Linux ubuntu one liner: `curl -o /usr/local/bin/autogit https://github.com/darklab8/darklab_autogit/releases/download/v2.1.0/autogit-linux-amd64 && chmod 777 /usr/local/bin/autogit`
- Windows:

  - Check your PATH bin locations with `echo %PATH%` and put binary file any of them or add to new one, be sure to rename from like `autogit-windows-amd64.exe` to `autogit.exe`
  - If u use `Git Bash`, recommendation to put into `~/bin` for usage in Git Bash only, or into `C:\Program Files\Git\cmd` for working in any terminal
  - U can add to any other PATH bin searchable locations or add a new one
- MacOS

  - (To be written where to put)

2. (Optional) init autogit.yml with `autogit init` command in the root of repository. change REPOSITORY_OWNER and REPOSITORY_NAME to yours
   1. Or init global one, or don't init at all (see 3d step of short opinionated instructioon)
3. run `autogit hook activate` to create `.git-hook` folder and enabling it in your git config for automated commit validation on pre-commit hook
   1. Or activate globally with flag `--global`
   2. `deactivate with `autogit hook deactivate`, using `--global` flag if necessary as well.

P.S. Current repository runs on configured autogit as well


