package operators

// errBoom este o eroare dummy folosită doar în teste
type errBoom struct{}

func (errBoom) Error() string { return "boom" }
