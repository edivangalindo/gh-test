package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "No tokens detected. Hint: cat tokens.txt | gh-test")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		token := scanner.Text()

		err := testGithubToken(token)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
	}
}

func testGithubToken(token string) error {
	fmt.Println("Testing token ->", token)

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)

	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "token "+token)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	var prettyJson bytes.Buffer
	err = json.Indent(&prettyJson, body, "", "\t")

	if err != nil {
		return err
	}

	fmt.Println(string(prettyJson.Bytes()))

	return nil
}
