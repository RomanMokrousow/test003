package cmd

import(
	"os/exec"
	"strings"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"fmt"
	"errors"
	"runtime"
)

type TRuneArr []rune;
func (src TRuneArr)IndexOf(r rune)(int){
	for k, v := range src {if r == v {return k}}
	return -1;
}

func Tokenize(src string)(Result []string, ERROR error){
	var (
		s = []rune(strings.Trim(src," \t\r\n"))
		l int = len(s)
		i int = 0
		inQuoted = ' '
		token string
	)
	for {
		if i>=l {Result = append(Result,token); break}
		if inQuoted != ' '{
			if s[i] == inQuoted {
				if (((i+1) < l) && (s[i+1] == inQuoted)) {
					token += string(inQuoted);
					i++;
				}else{
					inQuoted = ' ';
				}
			}else{
				token += string(s[i])
			}
		}else{
			if (TRuneArr([]rune("\"'")).IndexOf(s[i])>=0){
				if(((i+1) < l) && (s[i+1] == s[i])){
					token += string(s[i]);
					i++
				}else{
					inQuoted = s[i];
				}
			}else{
				if (TRuneArr([]rune(" \t\r\n")).IndexOf(s[i])>=0){
					Result = append(Result,token);
					token = "";
				}else{
					token += string(s[i]);
				}
			}
		}
		i++;
	}
	return
}

var cpConverter encoding.Encoding;

type TLocalCommand = func(params string)([]byte, error);
var LocalCommand = map[string]TLocalCommand{
	"ver": local_ver,
	"chcp": local_chcp,
	"gc": local_GC,
	"mem": local_Mem,
}
func local_Mem(params string)([]byte, error){
	var MS runtime.MemStats;
	runtime.ReadMemStats(&MS);
	return []byte(fmt.Sprintf("%#v",MS)),nil;
}
func local_GC(params string)([]byte, error){runtime.GC();return []byte("ready"),nil}
func local_ver(params string)([]byte, error){return []byte("0.0.0"),nil}
func local_chcp(params string)([]byte, error){
	params = strings.Trim(params," \t\r\n");
	if params == "" {return []byte(fmt.Sprintf("Current codepage converter is set to [%v]",cpConverter)),nil}
	if params == "clear" {cpConverter = nil;return []byte(fmt.Sprintf("Current codepage converter is set to [%v]",cpConverter)),nil}
	if params == "help" {return []byte("Usage: chcp [help | list | clear | <CodePageName>]"),nil}
	var flag bool = false;
	var Result string;
	var Err error;
	for _,v := range charmap.All {
		if params == "list" {Result += fmt.Sprintf("%v\t",v);flag = true;continue}
		if params == fmt.Sprintf("%v",v) {cpConverter = v;flag = true;break}
	}
	if !flag {Err = errors.New(fmt.Sprintf("Unknown codepage [%s]. Type 'chcp list' for get list of available values",params));Result = ""}
	return []byte(Result),Err
}

func Execute(command string)(string, error){
	var(
		cmd_o []byte
		cmd_err error
		cmd_arr []string
		cmd *exec.Cmd
		Result_Text string
		Result_Err error
	)
	command = strings.Trim(command," \t\r\n")

	//cmd_arr = strings.Split(command," ");
	cmd_arr,_ = Tokenize(command);
	fmt.Printf("Tokenized: %s\n",strings.Join(cmd_arr,"--"));

	var lc = LocalCommand[cmd_arr[0]]; if lc != nil {
		cmd_o,Result_Err = lc(strings.Join(cmd_arr[1:]," "))
	}else{
		cmd = exec.Command(cmd_arr[0], cmd_arr[1:]...);
		//cmd := exec.Command(cmd_arr[0],strings.Join(cmd_arr[1:]," "));
		cmd_o,Result_Err = cmd.Output();
	}
	if cmd_err == nil {
		if cpConverter != nil {
			var dec = cpConverter.NewDecoder();
			var err error;
			cmd_o,err = dec.Bytes(cmd_o); if err != nil {cmd_o = []byte("ERROR: CodePage converting error.")}
		}
		Result_Text = string(cmd_o);
		}else{
			Result_Err = errors.New("ERROR: " + cmd_err.Error());
		}
		return Result_Text, Result_Err;
}