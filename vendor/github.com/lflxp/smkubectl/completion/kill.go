package completion

var kill = Completion{
	Level:   "kill",
	IsShell: false,
	Shell:   `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
	Args: []string{
		"ARGS     DESCRIPTION",
		"-1       HUP (hang up)",
		"-2       INT (interrupt)",
		"-3       QUIT (quit)",
		"-6       ABRT (abort)",
		"-9       KILL (non-catchable, non-ignorable kill)",
		"-14      ALRM (alarm clock)",
		"-15      TERM (software termination signal)",
	},
	Daughter: map[string]Completion{
		"-1": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-2": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-3": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-6": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-9": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-14": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
		"-15": Completion{
			Level:       "-9",
			IsShell:     true,
			Shell:       `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
			IsCondition: true,
			Condition:   []string{"^kill.*"},
		},
	},
}
