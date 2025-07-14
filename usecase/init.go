package usecase

import "github.com/mockup-billing-engine/repo"

type Usecase struct {
	Repo *repo.Client
}

func Init(repo *repo.Client) Usecase {
	return Usecase{
		Repo: repo,
	}
}
