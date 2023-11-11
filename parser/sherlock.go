package parser

import (
	"encoding/json"
	"net/http"
    "errors"
)

const SHERLOCK_API = "https://mainnet-contest.sherlock.xyz/contests?&per_page=20"

type Contest struct {
    Id uint32 `json:"id"`
    RepositoryName string `json:"judging_repo_name"`
    PrizePool uint32    `json:"prize_pool"`
	Issues      []Issue
}

func (v Contest) GetTotalShares() float64{
    var shares float64
    for _, issue := range v.Issues { 
        totalFound := len(issue.Duplicates) +1
        shares += issue.GetShares() * float64(totalFound)
    }
    return shares
}

type sherlockRoot struct {
    Items   []Contest `json:"items"`

}

func GetContestDetails(repoName string) (Contest, error){

    var contest Contest
    var root sherlockRoot
    res, err := http.Get(SHERLOCK_API)
    if (err != nil) {
        return contest, err
    }
    err = json.NewDecoder(res.Body).Decode(&root)
    if ( err != nil) {
        return contest, err
    }
    for _, c := range root.Items {
        if(c.RepositoryName == repoName) {
            return c, nil
        }
    }
    //TODO: Call GH API And fetch Issues
    return contest, errors.New("Did not find Contest for reponame")

}
