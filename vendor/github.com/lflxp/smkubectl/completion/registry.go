package completion

const FZF_OPTION = "FZF_DEFAULT_OPTS=\"--height ${FZF_TMUX_HEIGHT:-50%} --header-lines=1 --min-height 15 --reverse $FZF_DEFAULT_OPTS --preview 'echo {}' --preview-window down:3:wrap $FZF_COMPLETION_OPTS\" fzf -m"

type Completion struct {
	Level       string
	Cmd         []string // command命令
	Args        []string // 参数
	IsShell     bool     // 是否执行kubectl获取还是直接cmd提示
	Shell       string   // 获取提示的命令
	IsCondition bool     // 是否有条件
	Condition   []string // 正则条件判断，支持多条件判断
	Daughter    map[string]Completion
}

var Completes map[string]Completion = map[string]Completion{
	"k":       kubectl,
	"kubectl": kubectl,
	"kill":    kill,
	"git":     git,
	"go":      golang,
	"showme":  showme,
}
