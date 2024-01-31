/*
Copyright © 2024 lflxp <382023823@qq.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const ZshCompletion = `compdef _k k                                                                                                                                                                                                                                                                    # 添加 _gmt 函数

__smkubectl-cli_debug()
{
    local file="$BASH_COMP_DEBUG_FILE"
    if [[ -n ${file} ]]; then
        echo "$*" >> "${file}"
    fi
}

_k() {
  local matches namespace result

  # 获取一次性结果
  result=$(command smkubectl smart ${LBUFFER})

  __smkubectl-cli_debug "namespace ${result}"
  namespace=$(command echo ${result}|head -1|awk '/^NAMESPACE/ {print "yes"}')
  __smkubectl-cli_debug "namespace ${namespace}"
  # 根据标题是否含有NAMESPACE动态切换显示结果
  if [[ -n "$namespace" ]]; then
    matches=$(command echo ${result} | FZF_DEFAULT_OPTS="--height ${FZF_TMUX_HEIGHT:-50%} --header-lines=1 --min-height 15 --reverse $FZF_DEFAULT_OPTS --preview 'echo {}' --preview-window down:3:wrap $FZF_COMPLETION_OPTS" smkubectl -m|awk '{print "-n "$1" "$2}'|tr '\n' ' ')
  else
    matches=$(command echo ${result} | FZF_DEFAULT_OPTS="--height ${FZF_TMUX_HEIGHT:-50%} --header-lines=1 --min-height 15 --reverse $FZF_DEFAULT_OPTS --preview 'echo {}' --preview-window down:3:wrap $FZF_COMPLETION_OPTS" smkubectl -m|awk '{print $1}'|tr '\n' ' ')
  fi
  #_describe 'command' ns
  if [ -n "$matches" ]; then
    LBUFFER="${LBUFFER}$matches"
    # LBUFFER="${BUFFER} $matches"
  fi
  zle reset-prompt
}

zle -N _k
#定义快捷键为： [Esc] [Esc]
` + "bindkey '`' _k"

const ZshCompletionDebug = `compdef _k k                                                                                                                                                                                                                                                                    # 添加 _gmt 函数

__smkubectl-cli_debug()
{
    local file="$BASH_COMP_DEBUG_FILE"
    if [[ -n ${file} ]]; then
        echo "$*" >> "${file}"
    fi
}

_k() {
  local matches namespace result

  # 获取一次性结果
  result=$(command smkubectl smart -d ${LBUFFER})

  __smkubectl-cli_debug "namespace ${result}"
  namespace=$(command echo ${result}|head -1|awk '/^NAMESPACE/ {print "yes"}')
  __smkubectl-cli_debug "namespace ${namespace}"
  # 根据标题是否含有NAMESPACE动态切换显示结果
  if [[ -n "$namespace" ]]; then
    matches=$(command echo ${result} | FZF_DEFAULT_OPTS="--height ${FZF_TMUX_HEIGHT:-50%} --header-lines=1 --min-height 15 --reverse $FZF_DEFAULT_OPTS --preview 'echo {}' --preview-window down:3:wrap $FZF_COMPLETION_OPTS" smkubectl -m|awk '{print "-n "$1" "$2}'|tr '\n' ' ')
  else
    matches=$(command echo ${result} | FZF_DEFAULT_OPTS="--height ${FZF_TMUX_HEIGHT:-50%} --header-lines=1 --min-height 15 --reverse $FZF_DEFAULT_OPTS --preview 'echo {}' --preview-window down:3:wrap $FZF_COMPLETION_OPTS" smkubectl -m|awk '{print $1}'|tr '\n' ' ')
  fi
  #_describe 'command' ns
  if [ -n "$matches" ]; then
    LBUFFER="${LBUFFER}$matches"
    # LBUFFER="${BUFFER} $matches"
  fi
  zle reset-prompt
}

zle -N _k
#定义快捷键为： [Esc] [Esc]
` + "bindkey '`' _k"

// zshCmd represents the zsh command
var zshCmd = &cobra.Command{
	Use:   "zsh",
	Short: "zsh completion script",
	Long:  `within kubectl k kill git`,
	Run: func(cmd *cobra.Command, args []string) {
		if debugLevel {
			fmt.Println(ZshCompletionDebug)
		} else {
			fmt.Println(ZshCompletion)
		}
	},
}

func init() {
	completionCmd.AddCommand(zshCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// zshCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// zshCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	zshCmd.Flags().BoolVarP(&debugLevel, "debug", "d", false, "Log Level")
}
