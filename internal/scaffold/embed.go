package scaffold

import _ "embed"

//go:embed templates/SKILL.md
var SkillMD []byte

//go:embed templates/CLAUDE.md
var ClaudeMD []byte
