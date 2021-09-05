package transaction

type ServiceTrans struct {
	repo RepoTransaction
}

func (s *ServiceTrans) InsertTransaction(t Transactions) error {

	err := s.repo.InserTransaction(t)
	if err != nil {
		return err
	}

	return nil

}
