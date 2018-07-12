package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"gopkg.in/src-d/go-git.v4"
	"github.com/mgutz/ansi"
	"os/exec"
	"bufio"
	"regexp"
	"sync"
	"time"
	"strconv"
)

var wg sync.WaitGroup

type Repo struct {
	url    string
	folder string
	status bool
}

func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

func red() func(string) string {
	return ansi.ColorFunc("red+")
}

func green() func(string) string {
	return ansi.ColorFunc("green+")
}

func yellow() func(string) string {
	return ansi.ColorFunc("yellow+")
}

func printRepo(r []Repo) {

	for _, rep := range r {
		if rep.status == true {
			fmt.Printf("%s %s\n", PadRight(rep.url, " ", 100), green()(getRepoStatus(rep.status)))
		} else {
			fmt.Printf("%s %s\n", PadRight(rep.url, " ", 100), yellow()(getRepoStatus(rep.status)))
		}
	}
}

func updateRepos(r []Repo, noUpdate bool) {
	fmt.Printf("\n\n")
	for _, repo := range r {
		if (repo.status == true) { //update

			if (noUpdate) {
				fmt.Printf(red()("Skip update of %s\n"), repo.folder)
			} else {
				fmt.Printf(green()("Updating repository %s\n"), repo.folder)
				repo, _ := git.PlainOpen("/home/sajith/stats/" + repo.folder)
				repo.Pull(&git.PullOptions{RemoteName: "origin"})
			}
		} else { //pull
			fmt.Printf(yellow()("Pulling repository %s\n"), repo.folder)
			_, err := git.PlainClone("/home/sajith/stats/"+repo.folder, false, &git.CloneOptions{URL: repo.url, Progress: os.Stdout,
			})

			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
		}
	}
}

func printStats(r []Repo, user string, from string, to string) {
	fmt.Printf(green()("\nGenerating statistics for user:%s \t From %s  To %s\n"), user, from, to)
	for _, repo := range r {
		fmt.Printf(yellow()(repo.folder + "___________________________________________________________________________\n"))
		output, err := exec.Command("/home/sajith/scratch/gitstats/getstat.sh", repo.folder, "kanishka.desilva@pagero.com", from, to).Output()
		if (err != nil) {
			log.Fatal(err)
		}
		fmt.Printf("%s", output)
	}
}

func getRepoStatus(stats bool) string {
	if stats == true {
		return "UPDATE"
	} else {
		return "PULL"
	}
}

func worker(output chan string, i int) {
	defer wg.Done()
	var r int = 2 * i
	fmt.Printf("running id %d with time %d\n", i, r)
	time.Sleep(time.Duration(r) * time.Second)
	fmt.Printf("done id %d with time %d\n", i, r)
	output <- "Woke up after " + strconv.Itoa(r)
}

func xxx(workers int) {
	var output = make(chan string, 100)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(output, i)
	}

	wg.Wait()
	close(output)

	for result := range output {
		fmt.Println(result)
	}
}

func main() {

	var noUpdate bool
	var user string
	var from string
	var to string

	flag.BoolVar(&noUpdate, "noupdate", true, "a string var")
	flag.StringVar(&user, "user", "", "a string var")
	flag.StringVar(&from, "from", "", "a string var")
	flag.StringVar(&to, "to", "", "a string var")
	flag.Parse()

	fmt.Println("user:", user)
	fmt.Println("from:", from)
	fmt.Println("to:", to)
	fmt.Println("noUpdate:", noUpdate)

	file, err := os.Open("git_repos.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close() //defer so it cleanup.

	fmt.Printf(green()("Repositories....\n"))
	var repoList = make([]Repo, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		//re, _ := regexp.Compile(`[a-z]+\.(git)`)
		re, _ := regexp.Compile(`(\w+-)?\w+\.git`)
		repo := scanner.Text()
		res := re.FindAllStringSubmatch(repo, -1)

		repoExists := true
		if _, err := os.Stat("/home/sajith/stats/" + res[0][0]); os.IsNotExist(err) {
			repoExists = false
		}
		repository := Repo{repo, res[0][0], repoExists}
		repoList = append(repoList, repository)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	printRepo(repoList)
	xxx(3)

	/*
	updateRepos(repoList, noUpdate)
	fmt.Print(green()("\nAll repos are update\n"))

	if len(user) > 0 {
		printStats(repoList, user, from, to)
	}*/

}
