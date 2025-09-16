package main

type BotMux struct {
	IssueClaim  chan IssueAction
	Bounty      chan BountyAction
	Achievement chan Achievement
	Solution    chan Solution
}

func NewBotMux(bufferSize int) *BotMux {
	return &BotMux{
		IssueClaim:  make(chan IssueAction, bufferSize),
		Bounty:      make(chan BountyAction, bufferSize),
		Achievement: make(chan Achievement, bufferSize),
		Solution:    make(chan Solution, bufferSize),
	}
}

// Initiates the bot multiplexer to list to incoming messages for various event
// streams and process them accordingly
func (bm *BotMux) Start() {
	for {
		select {
		case issue := <-bm.IssueClaim:
			if issue.Extend == true {
				ManageExtension(issue.ParticipantUsername, issue.Url)
			} else {
				ManageIssueClaim(issue.ParticipantUsername, issue.Claim, issue.Url)
			}
		case cash := <-bm.Bounty:
			ManageBounty(cash.ParticipantUsername, cash.Amount, cash.Action, cash.Url)
		case badge := <-bm.Achievement:
			ManageAchivement(badge.ParticipantUsername, badge.Type, badge.Url)
		case sol := <-bm.Solution:
			ManageSolution(sol.ParticipantUsername, sol.Merged, sol.Url)
		}
	}
}

// Gracefully shuts down the Bot multiplexer
func (bm *BotMux) Shutdown() {
	close(bm.IssueClaim)
	close(bm.Bounty)
	close(bm.Achievement)
	close(bm.Solution)
}
