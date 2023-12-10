# Algorithms

## scenario #1 - validator / Git commit validation / Changelog validation

You try to write git commit -m "feat: add rendering in format format"
your githook is activated and tries to parse your commit name accroding to git conventional commits standard. If unable, it will give you error and prevent commit

```mermaid
flowchart TD
  GitCommit[attempt to fixate commit like
  git commit -m 'feat: add md rendering'
  with 'autogit hook activate' enabled]
  RequestValidatingChangelog[Request changelog with --validate flag] --> TryParsingCommitMessage
  GitCommit --> TryParsingCommitMessage[Try parsing commit message\nto git conventional commit\ntype \ scope \ subject \ body \ footers]
  TryParsingCommitMessage --> ReportFail[Reporting errors if unable]
  TryParsingCommitMessage --> ContinuingValidation[Continue Validation]
  ContinuingValidation --> CheckOptionalValidationRulesIfEnabled[Check options validation rules\nif they are enabled]
  CheckOptionalValidationRulesIfEnabled --> CommitTypeInAllowedList[Commit type is\nin allowed list]
  CommitTypeInAllowedList --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> CheckOtherEnabledRulesInSettings[Check other enabled\nrules in settings]
  CheckOtherEnabledRulesInSettings --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> WhenAppliedRules[when applied rules]
  WhenAppliedRules --> IfCommit[if it was commit,\nthen fixate if passed rules,\nor cancel fixation]
  WhenAppliedRules --> IfChangelog[if it was changelog validation\nthen report no errors and exit code 0\nfor pipeline checks]
```

## scenario #2 - changelog / Your wish to see changelog of additions you made, what are new features, what are fixes. For user view

You wish to have changelog auto generated.

```mermaid
flowchart TD
    RequestingChangelog[Requesting changelog]
    RequestingChangelog --> ChangelogFromLatestCommitToPreviousTagVersion
    ChangelogFromLatestCommitToPreviousTagVersion[Requesting changelog from previous\ntag to latest commit]
    RequestingChangelog --> ChangelogFromChosenTagToPreviousTag
    ChangelogFromChosenTagToPreviousTag[Requesting changelog from chosen tag\nversion to previous tag version]
    ChangelogFromLatestCommitToPreviousTagVersion --> GenerateChangelog
    GenerateChangelog[Start generating changelog]
    ParseCommits[Parse commit in necessary tag range]
    ChangelogFromChosenTagToPreviousTag --> GenerateChangelog
    GenerateChangelog --> ParseCommits
    ParseCommits --> SelectAllowedTypesForRender
    SelectAllowedTypesForRender[Filter conventional commit `types` like `feat` allowed for render]
    SelectAllowedTypesForRender --> SubgroupIntoConventionalCommitScopes
    SubgroupIntoConventionalCommitScopes[Sub group commits according to conventional commit `scope`]
    GenerateChangelog --> CalculateNextSemver
    CalculateNextSemver[Calculate Next Semantic Version]
    CalculateNextSemver --> SendChangelogToRender
    SubgroupIntoConventionalCommitScopes --> SendChangelogToRender
    SendChangelogToRender[Receive changelog for render]
    SendChangelogToRender --> RenderChangelogMarkdown
    SendChangelogToRender --> RenderChangelogRst
    SendChangelogToRender --> RenderChangelogHtml
    RenderChangelogMarkdown[Render in markdown\n--implemented--]
    RenderChangelogRst[Render in rst\n--not implemented--]
    RenderChangelogHtml[Render in html\n--not implemented--]
```

### example of rendered changelog

[Full example of rendered changelog](https://github.com/darklab8/darklab_autogit/releases/tag/v0.3.0-rc.2)

## Scenario #3 - nextSemVer / You wish to know which next semantic version / semantic tag should be applied to your release.

Program checks if u made no commits, or only refactoring and styling.

- If u made no changes, then next version is same as previous one.
- If you made `fix`, then it increases PATCH version of semantic version. Your next version is `0.0.2` rendered
- if you made `feat` request, then next MINOR version is increased. Your next version is `0.1.0`
- if you made breaking changes, users should know `feat!` or `BREAKING CHANGE:`, then next version is `1.0.0`
- if u had no previous versions, it will calculate new one as `0.0.0` + calculated version changes

more detailed algorithm, accounting also prerelease version calculations:

```mermaid
flowchart TD
  RequestNextSemanticVersioning[Request next semantic versioning]
  RequestNextSemanticVersioning --> FindCommits[Find commits\nfrom HEAD^1 to previous stable semantic version like v0.3.0]
  FindCommits --> CalculateVersionChange[Calculate main version change\nChoose only ONE path]
  CalculateVersionChange --> MajorChange[if git conventional commits\nwith breaking changes\nlike feat! are present\nand it is not 0.*.* development mode\nor flag `--publish` is present,\nthen add MAJOR version\nand reset MINOR and PATCH to 0\nchange: +1.0.0]
  CalculateVersionChange --> MinorChange[If commits with `feat` type are present\nincrease MINOR version\nand reset PATCH version\nto zero change: *.+1.0,]
  CalculateVersionChange --> PatchChange[if only commits with `fix`\n are present\nthen change only PATCH\nchange: *.*.+1]
  MajorChange --> CalculatedMainVersion
  MinorChange --> CalculatedMainVersion
  PatchChange --> CalculatedMainVersion
  CalculatedMainVersion[Calculated main version]
  CalculatedMainVersion --> CalculatePrereleaseVersion[Calculate next prerelease version]
  CalculatePrereleaseVersion --> FindLatestPrerelease[Find latest alpha beta and rc versions\nwith scanning commits up to latest stable version\nExcept not counting latest commit\nChoose ONE, SEVERAL or ALL paths next:]
  FindLatestPrerelease --> AlphaFlag[if alpha flag is present\nincrease alpha version\nand mark for rendering]
  FindLatestPrerelease --> BetaFlag[if beta flag is present\nincrease beta version\nand mark for rendering]
  FindLatestPrerelease --> RcFlag[if rc-release candidate-\nflag is present\nincrease rc version\nand mark for rendering]
  AlphaFlag --> CombineIntoTotalPrereleaseVersion
  BetaFlag --> CombineIntoTotalPrereleaseVersion
  RcFlag --> CombineIntoTotalPrereleaseVersion
  CombineIntoTotalPrereleaseVersion[Combine into latest prerelease version]
  CalculatedMainVersion ----> AddBuildMetaData[Add build meta data\nas +$BuildMetaData\nto the end of version]
  CombineIntoTotalPrereleaseVersion --> OutputFinalSemanticVersion
  AddBuildMetaData --> OutputFinalSemanticVersion[Render Final Semantic Version]
```
