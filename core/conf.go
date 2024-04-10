package core

const ASCII_ART = `
  ██████  ▒█████   ██▓███   ██░ ██  ██▓ ▄▄▄      
▒██    ▒ ▒██▒  ██▒▓██░  ██▒▓██░ ██▒▓██▒▒████▄    
░ ▓██▄   ▒██░  ██▒▓██░ ██▓▒▒██▀▀██░▒██▒▒██  ▀█▄  
  ▒   ██▒▒██   ██░▒██▄█▓▒ ▒░▓█ ░██ ░██░░██▄▄▄▄██ 
▒██████▒▒░ ████▓▒░▒██▒ ░  ░░▓█▒░██▓░██░ ▓█   ▓██▒
▒ ▒▓▒ ▒ ░░ ▒░▒░▒░ ▒▓▒░ ░  ░ ▒ ░░▒░▒░▓   ▒▒   ▓▒█░
░ ░▒  ░ ░  ░ ▒ ▒░ ░▒ ░      ▒ ░▒░ ░ ▒ ░  ▒   ▒▒ ░
░  ░  ░  ░ ░ ░ ▒  ░░        ░  ░░ ░ ▒ ░  ░   ▒   
      ░      ░ ░            ░  ░  ░ ░        ░  ░
`

type Config struct {
	// AllErrors enables the printing of all Errors, not just the first 3
	AllErrors bool
	// Debug enables debug logging
	Debug bool
	// JIT enables the just in time compiler
	JIT bool
}

var CONF = Config{
	Debug: false,
}
