package settings

import "autogit/semanticgit/conventionalcommits/conventionalcommitstype"

type TypeAllowLists struct {
	SemVerMinorIncreasers []conventionalcommitstype.Type `yaml:"semver_minor_increases"`
	SemverPatchIncreasers []conventionalcommitstype.Type `yaml:"semver_patch_increases"`
	SemverNotAffects      []conventionalcommitstype.Type `yaml:"semver_is_not_affected"`
}

func (a TypeAllowLists) GetAllTypes() []conventionalcommitstype.Type {
	list := []conventionalcommitstype.Type{}
	list = append(list, a.SemVerMinorIncreasers...)
	list = append(list, a.SemverPatchIncreasers...)
	list = append(list, a.SemverNotAffects...)
	return list
}

type ValidationScheme struct {
	Sections struct {
		Hook struct {
			CommitMsg struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"commitMsg"`
		} `yaml:"hook"`
		// TODO, add ability to disable Changelog validations?
		// Changelog struct {
		// 	Enabled bool `yaml:"enabled"`
		// } `yaml:"changelog"`
	} `yaml:"sections"`
	Rules struct {
		Issue struct {
			Present bool `yaml:"present"`
		} `yaml:"issue"`
		Header struct {
			MaxLength int `yaml:"maxLength"`
			Type      struct {
				Lowercase  bool           `yaml:"lowercase"`
				Allowlists TypeAllowLists `yaml:"allowlists"`
			} `yaml:"type"`
			Scope struct {
				Present   bool                            `yaml:"present"`
				Lowercase bool                            `yaml:"lowercase"`
				Allowlist []conventionalcommitstype.Scope `yaml:"allowlist"`
			} `yaml:"scope"`
			Subject struct {
				MinWords int `yaml:"minWords"`
			} `yaml:"subject"`
		} `yaml:"header"`
	} `yaml:"rules"`
}
