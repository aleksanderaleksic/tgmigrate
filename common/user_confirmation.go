package common

import (
	"bufio"
	"fmt"
	"os"
)

func AskUserToConfirm() bool {
	fmt.Print("insert y or yes if you would like to continue: \n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	return input.Text() == "y" || input.Text() == "yes"
}
