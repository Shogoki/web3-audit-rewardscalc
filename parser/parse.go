package parser

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

type GQLRoot struct {
	Data data `json:"data"`
}
type data struct {
    Repository Repository `json:"repository"`
}
type Repository struct {
	Duplicates IssueEntry `json:"duplicates"`
	OpenIssues IssueEntry `json:"openIssues"`
}

type IssueEntry struct {
	Nodes []RawIssue `json:"nodes"`
}

type LabelEntry struct {
	Nodes []Label `json:"nodes"`
}

type RawIssue struct {
	Body   string     `json:"body"`
	Title  string     `json:"title"`
	Number int        `json:"number"`
	Labels LabelEntry `json:"labels"`
}
type Label struct {
	Name string `json:"name"`
}

type Issue struct {
	isMain     bool
	title      string
	body       string   //TODO: Do we need body here, maybe better bodyHTML?
	labels     []string //TODO: Do we need to expose labels here?
	severity   string
	number     int
	author     string
	duplicates []Issue
}

func (v Issue) GetShares() float64 {

	totalFound := float64(len(v.duplicates) + 1)
	sevMultiplier := 1.0
	if v.severity == "High" {
		sevMultiplier = 5.0
	}
	return sevMultiplier * math.Pow(0.9, (totalFound-1)) / totalFound
}

type Contest struct {
	repositoryName string
	potSize        uint16
	issues         []Issue
}

func (v RawIssue) toIssue() Issue {
	iss := Issue{}
	iss.body = v.Body
	iss.number = v.Number
	iss.title = v.Title
    iss.labels = make([]string, len(v.Labels.Nodes))
	for i, label := range v.Labels.Nodes {
		iss.labels[i] = label.Name
		if label.Name == "High" {
			iss.severity = "High"
		} else if label.Name == "Medium" {
			iss.severity = "Medium"
		}
	}
	iss.author = strings.Split(v.Body, "\n")[0]

	return iss
}
func (v IssueEntry) toIssues(isMain bool) []Issue {
	issues := make([]Issue, len(v.Nodes))
	for i, iss := range v.Nodes {
		issues[i] = iss.toIssue()
		issues[i].isMain = isMain
	}
	return issues
}
func (v RawIssue) getDuplicateId() (int, error) {

	lines := strings.Split(v.Body, "\n")
	var num string
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "Duplicate of #") {
			num = strings.Replace(lines[i], "Duplicate of #", "", 1)
			break
		} else if strings.Contains(lines[i], "Duplicate of https") {
			num = lines[i][strings.LastIndex(lines[i], "/")+1:]
			break
		}
	}
	retNum, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil || retNum == 0 {
		return 0, errors.New(fmt.Sprintf("Could not get Duplicate ID for issue %d ", v.Number))
	}
	return retNum, nil
}

func (v RawIssue) hasLabel(label string) bool {
    for _, lbl := range v.Labels.Nodes {
        if label == lbl.Name {
            return true
        }
    }
    return false
}
func contains(slice []string, s string) bool {
    for _, item := range slice {
        if s == item {
            return true
        }
    }
    return false

}
func GetContestIssues(gqlResponse []byte) ([]Issue, error) {
	var root GQLRoot
    decoder := json.NewDecoder(bytes.NewReader(gqlResponse))
	decoder.DisallowUnknownFields()
	//err := json.Unmarshal(gqlResponse, &root)
    err := decoder.Decode(&root)
	if err != nil {
		log.Print("Failed to Unmarshal", err)
		return nil, err
	}
	issues := root.Data.Repository.OpenIssues.toIssues(true)
	//Adding duplicates
	for _, dup := range root.Data.Repository.Duplicates.Nodes {
		id, err := dup.getDuplicateId()
		if err != nil {
            if( dup.hasLabel("Escalation Resolved")) {
                continue
            }
			return nil, err
		}
		for _, iss := range issues {
			if iss.number == id {
				iss.duplicates = append(iss.duplicates, dup.toIssue())
			}
		}
	}
	return issues, nil

}
