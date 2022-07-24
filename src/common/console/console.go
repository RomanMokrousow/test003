package console

import(
	"os"
	"bufio"
	"strings"
	"project.local/domain/src/common/cmd"
)

func Console(){
	var(
		cmd_text string
		cmd_result string
		cmd_err error
	)
	scn := bufio.NewScanner(os.Stdin);
	for {
		os.Stdout.Write([]byte(">"));
		scn.Scan();
		cmd_text = scn.Text();
		if (cmd_text == ""){os.Stdout.Write([]byte("ERROR: Empty line\n"));continue}
		cmd_text = strings.Trim(cmd_text," \t\r\n")
		if (cmd_text == "exit"){break}
		cmd_result, cmd_err = cmd.Execute(cmd_text);
		if cmd_err != nil {cmd_result += "\n" + cmd_err.Error()}
		os.Stdout.Write([]byte(cmd_result + "\n"));
	}	
}