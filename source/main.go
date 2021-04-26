package main

import (
	"fmt"
)

var (
	config configGlobal
)

func main() {
	if err := parseConfig(); err != nil {
		fmt.Println(err)
		return
	}

	/*
		Get statistics for all defined accounts
	*/
	var data []string
	for _, a := range config.Accounts {
		x, err := get(a)
		if err != nil {
			fmt.Println(err)
			continue
		}
		data = append(data, x)
	}

	/*
		Format the gathered statistics and print them to the stdout
	*/
	var out string
	for k, v := range data {
		if k > 0 {
			out += config.Delimiter
		}
		out += v
	}
	fmt.Println(out)
}

/*
	Get informations for specified account and format it
*/
func get(a Account) (ret string, err error) {
	c, err := NewChecker(a)
	if err != nil {
		return
	}

	err = c.Connect()
	if err != nil {
		return
	}
	defer c.Disconnect()

	s, err := c.GetStatus()
	if err != nil {
		return
	}

	ret = fmt.Sprintf("%s: %d/%d", a.Name, s.Unseen, s.Total)

	return
}
