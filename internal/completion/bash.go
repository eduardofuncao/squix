package completion

func GenerateBash() string {
	return `# Bash completion for squix
_squix_complete() {
    local cur prev words cword
    _init_completion || return

    # Call squix __complete with current arguments
    local completions
    completions=$(squix __complete "${words[@]:1}")

    # Filter completions based on current word
    COMPREPLY=($(compgen -W "$completions" -- "$cur"))
}

complete -F _squix_complete squix
`
}
