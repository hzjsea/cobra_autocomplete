cobra zsh-complete 模版

使用方法
假定项目名称为 echoz
0. 前期准备
修改cmd/commondd/root.go 下面的rootCmd， 将Use改成对应的项目名称 "echoz"
全局替换文件中的upx变成echoz 
1. 构建对应的二进制文件
**go build ./cmd/commondd/main.go**
2. 将二进制文件修改成对应命令名字
**mv main echoz**
3. 使用命令生成对应的补全脚本
**./echoz completion zsh > echoza**
4. 激活脚本
**source ./echoz**
5. 使用
临时使用
**./echoz \[tab\]\[tab\]**
长期使用,将下面的内容写入到zshrc中，然后**source zshrc**
```bash 
if [ $commands[qshell] ]; then
  source <(qshell completion zsh)
fi
```



Pass
七牛的sdk对cobra组合方式封装一个builder，后期有机会可以补一补 
另外生成的脚本文件并不只是适用于cobra的命令程序，如果你有实力或者时间，完全可以用其他的命令行工具去做

另外模版一共有四部分组成， 分别是头部 中部 和 尾部 还有自定义的部分
其中头部和尾部的内容如下所示， 已经被写在了 cmd/commondd/completion.go 文件中

头部 
```go
zsh_head := "#compdef upx\n"
			os.Stdout.Write([]byte(zsh_head))
			zsh_initialization := `
__upx_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__upx_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__upx_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__upx_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__upx_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__upx_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__upx_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__upx_filedir() {
	local RET OLD_IFS w qw
	__upx_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __upx_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__upx_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__upx_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
	printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__upx_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__upx_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__upx_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__upx_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__upx_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__upx_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__upx_type/g" \
	<<'BASH_COMPLETION_EOF'
`
```

尾部
```go
			tail := `
BASH_COMPLETION_EOF
}
__upx_bash_source <(__upx_convert_bash_to_zsh)
_complete upx 2>/dev/null
`
			os.Stdout.Write([]byte(tail))
```


中部
```go
			buf := new(bytes.Buffer)
			cmd.Root().GenBashCompletion(buf)
			os.Stdout.Write(buf.Bytes())
```
其中GenBashCompletion生成的代码结构如下所示
```bash
_upx_root_command()
{
    last_command="upx"

    command_aliases=()

    commands=()
    commands+=("completion") # 你加的命令
    commands+=("config")
    commands+=("help")
    commands+=("web")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--colorful") # 全局的flag
    flags+=("--config=")
    two_word_flags+=("--config")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}


_upx_config() # 子命令
{
    last_command="upx_config"

    command_aliases=()

    commands=()
    commands+=("create")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--colorful") # 子命令的flag or 全局的flag
    flags+=("--config=")
    two_word_flags+=("--config")

    must_have_one_flag=()
    must_have_one_noun=()
    must_have_one_noun+=("node")
    must_have_one_noun+=("pod")
    must_have_one_noun+=("replicationcontroller")
    must_have_one_noun+=("service")
    noun_aliases=()
}
```


另外自定义的部分，见 cmd/commands/root.go 文件 


## base one
https://github.com/qiniu/go-sdk
