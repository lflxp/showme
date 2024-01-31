package cmd

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"

	. "github.com/lflxp/smkubectl/completion"
)

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// 解析命令
// 只负责执行 不负责命令补全
func ParseCmd(in string) {
	var err error
	// 记录命令最后一个字符是否是空格
	isLastWorkSpace := false
	tmp := strings.Split(in, " ")
	result := []string{}
	// 清除空格和空字符串
	// args参数包含回车
	for index, v := range tmp {
		if v != "" && v != " " {
			result = append(result, v)
		}
		if index == len(tmp)-1 {
			if v == "" || v == " " {
				slog.Debug("命令最后包含空格或空字符")
				isLastWorkSpace = true
				// fmt.Printf("0 islaster %d %d works |%s| %v\n", index, len(tmp)-1, v, isLastWorkSpace)
			}
		}
		slog.Debug("in value", "index", index, "value", v, "result", result, "isLastWorkSpace", isLastWorkSpace)
	}

	slog.Debug("result", "data", result)
	if len(result) == 0 {
		slog.Debug("命令为空")
		return
	}

	if value, ok := Completes[result[0]]; ok {
		slog.Debug("命令匹配上已存在的Completes", "cmd", result[0])
		// 判断一级命令
		if len(result) == 1 {
			if value.IsShell {
				slog.Debug(result[0], " 一级命令 执行自动Shell命令 ", value.Shell)
				execCmd(value.Shell)
			} else {
				slog.Debug(result[0], " 一级命令 打印CMD长度", len(value.Cmd), "打印ARGS长度", len(value.Args))
				fmt.Println(strings.Join(value.Cmd, "\n"))
				fmt.Println(strings.Join(value.Args, "\n"))
			}
		} else {
			// 判断一级后面的最后一个命令
			if value_daughter, ok := value.Daughter[result[len(result)-1]]; ok {
				slog.Debug("仅针对第一级命令有效", "匹配命令", result[len(result)-1])
				// 仅针对第一级命令有效
				rs := []string{}
				rs, err := execCompletion(rs, result, "", &value_daughter, -1, false, true)
				if err != nil {
					slog.Error(err.Error())
					return
				}
				fmt.Printf(strings.Join(rs, "\n"))
			} else {
				// 根据命令长度智能补全命令
				if len(result) == 2 {
					slog.Debug("智能补全二级命令")
					t_one := []string{}
					for index, k := range value.Cmd {
						if index == 0 {
							found, err := regexp.MatchString("(COMMAND|NAME|NAMESPACE|ARGS).*", k)
							if err != nil {
								slog.Error("正则匹配固定标题错误", "COMMAND|NAME|NAMESPACE|ARGS", err.Error())
								return
							}

							if !found {
								slog.Debug("未匹配到固定标题，设置默认标题")
								fmt.Println("COMMAND")
							} else {
								slog.Debug("匹配到固定标题", "TITLE", k)
								fmt.Println(k)
							}
						}
						if strings.HasPrefix(k, result[1]) {
							t_one = append(t_one, strings.Replace(k, result[1], "", 1))
						}
					}
					fmt.Printf(strings.Join(t_one, "\n"))
				} else if len(result) == 3 && !isLastWorkSpace {
					// 补全命令
					slog.Debug("智能补全三级命令")
					t_two := []string{}
					if value_3, ok := value.Daughter["get"]; ok {
						t_two, err = execCompletion(t_two, result, "", &value_3, 1, true, false)
						if err != nil {
							slog.Error(err.Error())
							return
						}
					}
					fmt.Printf(strings.Join(t_two, "\n"))
				} else {
					// 补全数据
					slog.Debug("智能补全多级命令")
					t_two := []string{}
					// 优先匹配value.Daughter里面的值
					// 补全上一个命令的结果，查询并替换已有数据
					// 没有发现命令则获取上一步的命令 获取资源后进行资源匹配
					if len(t_two) == 0 && isLastWorkSpace {
						// 补全实时数据结果
						// 补全数据 有空格
						// 如果获取最后一个参数无数据 则执行该命令
						slog.Debug("补全实时数据结果", "命令是否包含空格", isLastWorkSpace)
						var cmd string

						switch result[0] {
						case "kubectl", "k", "kk", "k8s":
							// 判断最后一个参数是否是命令行参数
							isCmds := false
							// TODO: 查询速度慢 简写无法处理
							// 要去kubectl api-resources结果中获取
							for _, c := range value.Cmd {
								if strings.HasPrefix(c, result[len(result)-1]) {
									cmd = fmt.Sprintf("kubectl get %s -A", result[len(result)-1])
									isCmds = true
									break
								}
							}

							// 没有匹配上命令行参数
							// 抽离namespace和pod
							// 查询pod数据为基础
							if !isCmds {
								ns, name, _ := absoftNS(result)
								if ns != "" && name != "" {
									cmd = fmt.Sprintf("kubectl get pod -n %s %s", ns, name)
								} else if ns != "" && name == "" {
									cmd = fmt.Sprintf("kubectl get pod -n %s", ns)
								} else {
									cmd = "kubectl get pod -A"
								}
							}
							slog.Debug("自动匹配一级命令默认数据源", "一级命令", "k kubectl kk k8s", "数据执行命令", cmd, "是否是api-resources资源", isCmds)

						case "git":
							cmd = `git branch -a|tr '*' ' '|awk '{for(i=1;i<=NF;++i) printf $i "\t";printf "\n"}'`
							fmt.Println("BRANCH")
							slog.Debug("自动匹配一级命令默认数据源", "一级命令", "git", "数据执行命令", cmd, "默认标题", "BRANCH")
						default:
							cmd = "ps -ef"
							slog.Debug("自动匹配一级命令默认数据源失败", "默认数据执行命令", cmd)
						}

						t_two, err = execCompletion(t_two, result, cmd, nil, 1, true, false)
						if err != nil {
							slog.Error(err.Error())
							return
						}
					} else if len(t_two) == 0 && !isLastWorkSpace {
						// 补全命令
						slog.Debug("补全实时数据结果", "命令是否包含空格", isLastWorkSpace)
						if target, ok := value.Daughter["get"]; ok {
							t_two, err = execCompletion(t_two, result, "", &target, 1, true, false)
							if err != nil {
								slog.Error(err.Error())
								return
							}
						} else {
							slog.Error("key get not exist\n")
							return
						}

						// 补全可能缺失的数据
						slog.Debug("补全可能缺失的数据", "倒数第二个命令", result[len(result)-2], "target", 1)
						if value_maybe, ok := value.Daughter[result[len(result)-2]]; ok {
							t_two, err = execCompletion(t_two, result, "", &value_maybe, 1, true, false)
							if err != nil {
								slog.Error(err.Error())
								return
							}
						}
						// 补全命令无效 获取上级命令的结果 并补全数据prefix数据
						slog.Debug("补全命令无效 获取上级命令的结果 并补全数据prefix数据", "倒数第三个命令", result[len(result)-2], "target", 2)
						if len(t_two) == 0 {
							if value_daughter2, ok := value.Daughter[result[len(result)-2]]; ok {
								t_two, err = execCompletion(t_two, result, "", &value_daughter2, 2, true, false)
								if err != nil {
									slog.Error(err.Error())
									return
								}
							}
						}
					} else {
						// fmt.Println("12")
						// fmt.Printf("t_two %v %b\n", t_two, isLastWorkSpace)
						slog.Debug("补全命令失败", "结果集", t_two, "命令是否包含空格", isLastWorkSpace)
					}
					fmt.Printf(strings.Join(t_two, "\n"))
				}
			}
		}
	} else {
		// fmt.Println("7")
		// execCmd(in)
		slog.Debug("未匹配到一级命令", "命令", result[0])
	}
}

// 抽离namespace和name
func absoftNS(in []string) (string, string, error) {
	var namespace, name string
	var err error

	for index, value := range in {
		if value == "-n" {
			if len(in) > index+1 {
				namespace = in[index+1]
			} else {
				err = fmt.Errorf("namespace is empty")
				return namespace, name, err
			}

			if len(in) > index+2 {
				name = in[index+2]
				break
			} else {
				err = fmt.Errorf("name is empty")
				return namespace, name, err
			}
		}
	}
	return namespace, name, nil
}

// 执行Completion
// cmd不为空时 ==> isLastWorkSpace == true
// cmd为空时 ==> isLastWorkSpace == false
// keepHeader 是否保留第一行数据
// in 当前输入的数据
// needRawCmd 是否需要原始命令
func execCompletion(result, in []string, cmd string, daughter *Completion, target int, keepHeader, first bool) ([]string, error) {
	slog.Debug("执行自动匹配", "cmd", cmd, "target", target, "keepHeader", keepHeader, "first", first)
	if daughter != nil && cmd == "" {
		if daughter.IsCondition {
			slog.Debug("执行条件判断", "IsCondition", daughter.IsCondition, "Condition", strings.Join(daughter.Condition, " "))
			var (
				found bool
				err   error
			)
			// 多条件过滤
			for _, v := range daughter.Condition {
				found, err = regexp.MatchString(v, strings.Join(in, " "))
				if err != nil {
					return result, err
				}
			}
			if found {
				var shell string
				// 硬编码
				// 专门处理pod类查询
				switch daughter.Level {
				case "show containers":
					// 获取pod的namespace和name
					ns, name, err := absoftNS(in)
					if err != nil {
						return result, err
					}
					shell = fmt.Sprintf("kubectl get pod -n %s %s -o jsonpath='{.spec.containers[*].name}'|tr ' ' '\\n'", ns, name)
					fmt.Println("CONTAINERS")
					slog.Debug("kubectl -c 精确匹配", "shell", shell)
				default:
					shell = daughter.Shell
				}

				execCmd(shell)
			}
		} else {
			if daughter.IsShell {
				// 实时补全数据结果
				if first {
					slog.Debug("执行匹配Completes", "shell", daughter.Shell, "是否一级命令", first)
					execCmd(daughter.Shell)
				} else {
					slog.Debug("执行匹配Completes, 提供命令数据后缀", "shell", daughter.Shell, "keepheader", keepHeader)
					rs_string, err := execCmdString(daughter.Shell)
					if err != nil {
						return result, err
					}
					// 补全 -n中提供的数据后缀
					for index, v := range strings.Split(rs_string, "\n") {
						if keepHeader && index == 0 {
							result = append(result, v)
						} else {
							if strings.HasPrefix(v, in[len(in)-1]) {
								result = append(result, strings.Replace(v, in[len(in)-1], "", 1))
							}
						}
					}
				}
			} else {
				// 补全cmd命令后缀
				if first {
					slog.Debug("打印一级命令列表和参数列表")
					if daughter.Cmd != nil {
						fmt.Println(strings.Join(daughter.Cmd, "\n"))
					}
					if daughter.Args != nil {
						fmt.Println(strings.Join(daughter.Args, "\n"))
					}
				} else {
					slog.Debug("补全cmd命令后缀")
					for _, k := range daughter.Cmd {
						if strings.HasPrefix(k, strings.TrimSpace(in[len(in)-target])) {
							t_rs := strings.Replace(k, strings.TrimSpace(in[len(in)-target]), "", 1)
							if t_rs != "" && t_rs != " " {
								result = append(result, t_rs)
							}
						}
					}
					for _, k := range daughter.Args {
						if strings.HasPrefix(k, strings.TrimSpace(in[len(in)-target])) {
							t_rs := strings.Replace(k, strings.TrimSpace(in[len(in)-target]), "", 1)
							if t_rs != "" && t_rs != " " {
								result = append(result, t_rs)
							}
						}
					}
				}
			}
		}
	} else if daughter == nil && cmd != "" {
		slog.Debug("执行默认命令数据收集命令", "cmd", cmd)
		rs_string, err := execCmdString(cmd)
		if err != nil {
			// 执行命令错误时执行整个命令
			if strings.Contains(err.Error(), "exit") {
				CMD := strings.Join(in, " ")
				if !strings.Contains(CMD, "edit") && in[target] != "edit" {
					slog.Debug("执行命令错误时执行整个命令", "ERROR", err.Error(), "CMD", CMD, "是否包含EDIT", false)
					execCmd(CMD)
				} else {
					slog.Debug("执行命令错误时执行整个命令", "ERROR", err.Error(), "CMD", strings.Replace(CMD, "edit", "get", 1), "是否包含EDIT", false)
					execCmd(strings.Replace(CMD, "edit", "get", 1))
				}
			} else {
				return result, err
			}
		} else {
			// 补全 -n中提供的数据后缀
			slog.Debug("执行默认命令数据收集命令[后续命令补全]", "cmd", cmd)
			for _, v := range strings.Split(rs_string, "\n") {
				if strings.TrimSpace(v) != "" {
					result = append(result, v)
				}
			}
		}
	} else {
		return result, fmt.Errorf("unsupport type cmd %s", cmd)
	}
	return result, nil
}

func execCmdString(in string) (string, error) {
	cmd := exec.Command("sh", "-c", in)
	cmd.Stdin = os.Stdin
	var out bytes.Buffer
	cmd.Stdout = &out // 标准输出
	var errout bytes.Buffer
	// cmd.Stderr = os.Stderr // 标准错误
	cmd.Stderr = &errout // 标准错误
	err := cmd.Run()
	if err != nil {
		// if &errout != nil {
		// 	log.Fatalf("cmd.Run() failed with %s err: %s errout: %s\n", in, err, errout.String())
		// } else {
		// 	log.Fatalf("cmd.Run() failed with %s err: %s\n", in, err)
		// }
		slog.Error(err.Error(), "CMD", in, "errout", errout.String())
		return "", err
	}
	return out.String(), nil
}

func execCmd(in string) error {
	cmd := exec.Command("sh", "-c", in)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout // 标准输出
	// var errout bytes.Buffer
	cmd.Stderr = os.Stderr // 标准错误
	// cmd.Stderr = &errout // 标准错误
	err := cmd.Run()
	if err != nil {
		// if &errout != nil {
		// 	log.Fatalf("cmd.Run() failed with %s err: %s errout: %s\n", in, err, errout.String())
		// } else {
		// 	log.Fatalf("cmd.Run() failed with %s err: %s\n", in, err)
		// }
		slog.Error(err.Error(), "CMD", in)
		return err
	}
	return nil
}
