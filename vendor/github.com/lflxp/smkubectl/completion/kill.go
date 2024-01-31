package completion

var kill = Completion{
	Level:   "kill",
	IsShell: true,
	Shell:   `ps -ef|awk '{for(i=2;i<=NF;++i) printf $i "\t";printf "\n"}'`,
}
