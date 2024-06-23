package shares

type Share struct {
	Figi string
	Uid  string
}

type Shares []Share

func (o Shares) GetFigis() []string {
	r := make([]string, 0)
	for _, v := range o {
		r = append(r, v.Figi)
	}
	return r
}

func (o Shares) GetUids() []string {
	r := make([]string, 0)
	for _, v := range o {
		r = append(r, v.Uid)
	}
	return r
}
