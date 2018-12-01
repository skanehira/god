package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

type history struct {
	file string
}

func newHistory() *history {
	usr, err := user.Current()
	if err != nil {
		log.Println(fmt.Sprintf("cannot get home dir: %s", err))
		os.Exit(0)
	}

	return &history{
		filepath.Join(usr.HomeDir, ".docker_cmd_history"),
	}
}

func (h history) loadHistory() []string {
	var cmdArgs []string

	file, err := os.OpenFile(h.file, os.O_RDONLY, 0666)
	if err != nil {
		return cmdArgs
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 4096)
	for {
		line, _, err := reader.ReadLine()
		cmdArgs = append(cmdArgs, string(line))
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(fmt.Sprintf("cannot read file: %s", err))
			os.Exit(0)
		}
	}

	return cmdArgs[:len(cmdArgs)-1]
}

func (h history) saveHistory(cmds []string) error {
	file, err := os.OpenFile(h.file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, strings.Join(cmds, " "))
	return nil
}

func executeCmd(cmds []string) error {
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func isEOF(err error) bool {
	if err == promptui.ErrEOF {
		return true
	}

	return false
}

func isInterrupt(err error) bool {
	if err == promptui.ErrInterrupt {
		return true
	}

	return false
}

func main() {
	// get file path
	h := newHistory()

	// get args
	flag.Parse()
	args := flag.Args()

	// if have args
	var cmdArgs []string

	if len(args) != 0 {
		cmdArgs = append([]string{"docker"}, args...)
		// save history
		h.saveHistory(cmdArgs)

	} else {
		// load history
		history := h.loadHistory()

		// select interface
		list := promptui.Select{
			Label: "command history",
			Templates: &promptui.SelectTemplates{
				Label:  ` {{ . | green }}`,
				Active: fmt.Sprintf(`%s {{ . | underline | red }}`, promptui.IconSelect),
			},
			Searcher: func(input string, index int) bool {
				cmd := history[index]
				name := strings.Replace(strings.ToLower(cmd), " ", "", -1)
				input = strings.Replace(strings.ToLower(input), " ", "", -1)
				return strings.Contains(name, input)
			},
			Items: history,
			Size:  50,
		}

		_, cmd, err := list.Run()

		if err != nil {
			if isEOF(err) || isInterrupt(err) {
				os.Exit(0)
			}

			log.Println(err)
			os.Exit(0)
		}

		// get selected
		cmdArgs = strings.Split(cmd, " ")
	}

	// execute command
	if err := executeCmd(cmdArgs); err != nil {
		log.Println(fmt.Sprintf("cannot run cmd: %s", err))
		os.Exit(0)
	}
}
