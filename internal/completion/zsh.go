package completion

func GenerateZsh() string {
	return `#compdef squix

_squix() {
    local -a completions
    completions=(${(f)"$(squix __complete $words[2,-1])"})

    if [[ -n "$completions" ]]; then
        _describe 'squix' completions
    else
        _files  # Fallback to files
    fi
}

# Register completion function (works when sourced in .zshrc or installed as file)
compdef _squix squix
`
}
