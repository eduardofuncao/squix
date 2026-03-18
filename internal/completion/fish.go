package completion

func GenerateFish() string {
	return `# Fish completion for squix
complete -c squix -f -a "(squix __complete (commandline -opc)[2..-1])"
`
}
