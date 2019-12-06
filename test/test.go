package main

import (
	"bufio"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/input"
	"os"
)

func main(){
	buf := bufio.NewReader(os.Stdin)
	prompt := fmt.Sprintf("Password to sign with '%s':", "jjon")
	passphrase, err := input.GetPassword(prompt, buf)
	fmt.Println("passphrase",passphrase, err)
}
