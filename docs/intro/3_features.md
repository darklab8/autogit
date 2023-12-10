# Features

- [git conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) validations (and optional other ones) on pre-commit hook
- validation of your commit history on `autogit changelog --validate` request
- automatic next [semantic versioning](https://semver.org/spec/v2.0.0.html) calculation for your product release. `autogit semver`
  - creating and pushing tag on `autogit semver --tag --push`. Changelog is automatically auto inserted into this tag.
- auto generating [changelogs](https://github.com/darklab8/darklab_autogit/releases/tag/v0.3.0-rc.2) new features and bug fixes for your next product release
- CI friendly binary file for any OS and arhictecture. Development with CI in mind. [CI examples](https://github.com/darklab8/darklab_autogit/tree/master/.github/workflows). Compiled for:
  - linux-amd64
  - linux-arm64
  - linux-386
  - linux-arm
  - windows-amd64
  - windows-386.exe
  - windows-arm64.exe
  - windows-arm.exe
  - macos-amd64
  - macos-arm64
- Contains inbuilt git. Not requiring git to be installed for its functionality
