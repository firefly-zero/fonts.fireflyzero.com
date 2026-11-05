package generator

type Encoding struct {
	ID    uint8
	Slug  string
	Title string
	Emoji string
}

var encodings = []Encoding{
	{
		ID:    0x0,
		Slug:  "ascii",
		Title: "ASCII",
		Emoji: "🇺🇸",
	},
	{
		ID:    0x1,
		Slug:  "iso_8859_1",
		Title: "Latin-1, Western European",
		Emoji: "🇵🇹",
	},
	{
		ID:    0x2,
		Slug:  "iso_8859_2",
		Title: "Latin-2, Central European",
		Emoji: "🇩🇪",
	},
	{
		ID:    0x3,
		Slug:  "iso_8859_3",
		Title: "Latin-3, South European",
		Emoji: "🇲🇹",
	},
	{
		ID:    0x4,
		Slug:  "iso_8859_4",
		Title: "Latin-4, North European",
		Emoji: "🇪🇪",
	},
	{
		ID:    0x5,
		Slug:  "iso_8859_9",
		Title: "Latin-5, Turkish",
		Emoji: "🇹🇷",
	},
	{
		ID:    0x6,
		Slug:  "iso_8859_10",
		Title: "Latin-6, Nordic",
		Emoji: "🇳🇴",
	},
	{
		ID:    0x7,
		Slug:  "iso_8859_13",
		Title: "Latin-7, Baltic Rim",
		Emoji: "🇵🇱",
	},
	{
		ID:    0x8,
		Slug:  "iso_8859_14",
		Title: "Latin-8, Celtic",
		Emoji: "🇮🇪",
	},
	{
		ID:    0x9,
		Slug:  "iso_8859_15",
		Title: "Latin-9 (revised Latin-1)",
		Emoji: "🇵🇹",
	},
	{
		ID:    0xA,
		Slug:  "iso_8859_16",
		Title: "Latin-10: South-East European",
		Emoji: "🇷🇴",
	},
	{
		ID:    0xB,
		Slug:  "iso_8859_5",
		Title: "Latin/Cyrillic",
		Emoji: "🇷🇺",
	},
	{
		ID:    0xC,
		Slug:  "iso_8859_7",
		Title: "Latin/Greek",
		Emoji: "🇬🇷",
	},
	{
		ID:    0xD,
		Slug:  "jis_x0201",
		Title: "Japanese katakana(halfwidth)",
		Emoji: "🇯🇵",
	},
}
