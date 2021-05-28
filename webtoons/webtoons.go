package webtoons

type Info struct {
	Title      string // must be folder name as well
	Subscriber string
	Rating     string
	Summary    string
	Creators   []string
	Episode    string
	End        int
}

type Option struct {
	TitleNum string
	Genre    string
	Title    string
	Start    int
	End      int
	Workers  int8
	Verbose  bool
}
