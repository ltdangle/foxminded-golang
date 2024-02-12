package model

import (
	"github.com/jmoiron/sqlx"
)

type WriteVoteRepo struct {
	db *sqlx.DB
}

func NewWriteVoteRepo(db *sqlx.DB) *WriteVoteRepo {
	return &WriteVoteRepo{db: db}
}
func (repo *WriteVoteRepo) Save(vote *Vote) error {
	query := `	
INSERT INTO vote 
	(from_user, to_user, vote, updated_at) 
VALUES 
	(:from_user, :to_user, :vote, :updated_at);`
	_, err := repo.db.NamedExec(query, vote)

	if err != nil {
		return err
	}

	return nil
}
func (repo *WriteVoteRepo) Update(vote *Vote) error {
	query := `
UPDATE vote SET vote = :vote WHERE from_user=:from_user AND to_user=:to_user; `

	_, err := repo.db.NamedExec(query, vote)

	if err != nil {
		return err
	}

	return nil
}
func (repo *WriteVoteRepo) FindUserVote(from_user string, to_user string) (*Vote, error) {
	var vote []Vote
	query := `
SELECT from_user, to_user, vote, updated_at FROM vote
WHERE from_user = ? AND to_user = ?;`
	err := repo.db.Select(&vote, query, from_user, to_user)
	if err != nil {
		return nil, err
	}
	if vote == nil {
		return nil, nil
	}

	return &vote[0], nil

}
