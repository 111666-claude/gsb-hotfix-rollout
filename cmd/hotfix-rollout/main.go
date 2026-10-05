// Command hotfix-rollout 跑灰度开关样例。
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/rollout/rollout"
)

func players(count int) []string {
	out := []string{}
	for index := 0; index < count; index++ {
		out = append(out, fmt.Sprintf("p-%03d", index))
	}
	return out
}

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("hotfix-rollout", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "stable", "stable / version / sorted / work")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	one, two := rollout.New(3, 1), rollout.New(3, 4)
	switch *sample {
	case "stable":
		first, _ := one.Apply(3, "flag", 50, players(200))
		second, _ := two.Apply(3, "flag", 50, players(200))
		same := fmt.Sprint(first.Sorted()) == fmt.Sprint(second.Sorted())
		body, _ := json.Marshal(map[string]bool{"same": same})
		fmt.Fprintln(stdout, string(body))
	case "version":
		_, err := one.Apply(3, "flag", 50, players(10))
		book := rollout.New(2, 1)
		_, stale := book.Apply(1, "flag", 50, players(10))
		body, _ := json.Marshal(map[string]bool{"first_ok": err == nil, "stale_ok": stale == nil})
		fmt.Fprintln(stdout, string(body))
	case "sorted":
		result, _ := one.Apply(3, "flag", 50, players(50))
		body, _ := json.Marshal(map[string]bool{"sorted": fmt.Sprint(result.Hit) == fmt.Sprint(result.Sorted())})
		fmt.Fprintln(stdout, string(body))
	case "work":
		result, _ := one.Apply(3, "flag", 50, players(3000))
		body, _ := json.Marshal(map[string]int{"scanned": one.Scanned(), "hit": len(result.Hit)})
		fmt.Fprintln(stdout, string(body))
	default:
		fmt.Fprintln(stderr, "需要 --sample stable|version|sorted|work")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
