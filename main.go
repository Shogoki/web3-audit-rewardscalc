package main

import (
    "fmt"
    "shogoki/audit-rewardscalc/parser"
    "os"
    "log"
    "io"
)

func main() {
    stdIn, err := io.ReadAll(os.Stdin)
    if(err != nil) {
        panic(err)
    }
    issues, err := parser.GetContestIssues(stdIn)
    if(err != nil) {
        log.Fatal("Failed to get Issues")
    }
    fmt.Println(len(issues))
}
