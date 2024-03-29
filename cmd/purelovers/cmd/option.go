package cmd

type Option struct {
	Login    string `short:"l" help:"Mail address."`
	Password string `short:"p" help:"Password."`
}

type Arg struct {
	Option
	Dump struct {
		Fav struct {
			Casts DumpFavoriteCasts `cmd:""`
			Shops DumpFavoriteShops `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Restore struct {
		Fav struct {
			Casts RestoreFavoriteCasts `cmd:""`
			Shops RestoreFavoriteShops `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Version PrintVersion `cmd:"" hidden:""`
}
