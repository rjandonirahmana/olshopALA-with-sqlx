package payment

type Payment struct {
	ID      int    `db:"id"`
	Methode string `db:"methode"`
}
