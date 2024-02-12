package usecase

import (
	"fmt"
	"jwt/pkg/model"
	"time"
)

const UPVOTE = 1
const DOWNVOTE = -1
const VOTE_DELAY_HRS = 24

type Vote struct {
	repo *model.WriteVoteRepo
}

func NewVoteUsecase(repo *model.WriteVoteRepo) *Vote {
	return &Vote{repo: repo}
}

func (ucase *Vote) Upvote(from_user string, to_user string) error {
	return ucase.vote(from_user, to_user, UPVOTE)
}

func (ucase *Vote) Downvote(from_user string, to_user string) error {
	return ucase.vote(from_user, to_user, DOWNVOTE)
}

func (ucase *Vote) vote(from_user string, to_user string, voteVal int) error {
	// find existing vote
	vote, err := ucase.repo.FindUserVote(from_user, to_user)
	if err != nil {
		return err
	}

	// vote doesn't exist, create new and upvote
	if vote == nil {
		return ucase.createNewVote(from_user, to_user, voteVal)
	}

	// Check for vote delay.
	duration := time.Since(*vote.UpdatedAt)
	if duration.Hours() < VOTE_DELAY_HRS {
		return fmt.Errorf("cannot vote more than once in %d hours", VOTE_DELAY_HRS)
	}

	// vote exists, update it
	return ucase.updateVote(from_user, to_user, voteVal)
}

func (ucase *Vote) createNewVote(from_user string, to_user string, vote int) error {
	// Validate request.
	validationErr := ucase.validateRequest(from_user, to_user, vote)
	if validationErr != nil {
		return validationErr
	}

	now := time.Now()
	wVote := &model.Vote{
		FromUser:  &from_user,
		ToUser:    &to_user,
		Vote:      &vote,
		UpdatedAt: &now,
	}
	err := ucase.repo.Save(wVote)
	if err != nil {
		return err
	}
	return nil
}

func (ucase *Vote) validateRequest(from_user string, to_user string, vote int) error {
	// Validate vote value.
	if vote != UPVOTE && vote != DOWNVOTE {
		return fmt.Errorf("vote value must be either %d or %d", UPVOTE, DOWNVOTE)
	}

	// User cannot vote for oneself.
	if from_user == to_user {
		return fmt.Errorf("user cannot vote for oneself")
	}

	return nil
}
func (ucase *Vote) updateVote(from_user string, to_user string, vote int) error {
	now := time.Now()
	wVote := &model.Vote{
		FromUser:  &from_user,
		ToUser:    &to_user,
		Vote:      &vote,
		UpdatedAt: &now,
	}
	err := ucase.repo.Update(wVote)
	if err != nil {
		return err
	}
	return nil
}
