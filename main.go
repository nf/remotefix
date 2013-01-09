// remotefix switches the origin remote for GitHub repositories installed by
// the go tool from using HTTPS to SSH.
package main

import (
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("git", "remote", "-v", "show").CombinedOutput()
	check(err)
	var u *url.URL
	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		if fields[0] != "origin" {
			log.Fatalln("non-origin remote:", fields[0])
		}
		if !strings.HasPrefix(fields[1], "https://github.com/") {
			log.Fatalln("non-https URL:", fields[1])
		}
		if u != nil {
			continue
		}
		u, err = url.Parse(fields[1])
		check(err)
	}
	if u == nil {
		log.Fatal("no origin remote found")
	}
	newRemote := fmt.Sprintf("git@github.com:%s.git", u.Path[1:])
	check(exec.Command("git", "remote", "rm", "origin").Run())
	check(exec.Command("git", "remote", "add", "origin", newRemote).Run())
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
