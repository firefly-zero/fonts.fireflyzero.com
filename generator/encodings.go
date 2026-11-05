package generator

type Encoding struct {
	ID      uint8
	Slug    string
	Title   string
	Emoji   string
	Pangram string
	Chars   []rune
}

var encodings = []Encoding{
	{
		ID:      0x0,
		Slug:    "ascii",
		Title:   "ASCII",
		Emoji:   "🇺🇸",
		Pangram: "The quick brown fox jumps over the lazy dog.",
	},
	{
		ID:      0x1,
		Slug:    "iso_8859_1",
		Title:   "Latin-1, Western European",
		Emoji:   "🇵🇹",
		Pangram: "Vejo galã sexy pôr quinze kiwis à força em baú achatado.",
	},
	{
		ID:      0x2,
		Slug:    "iso_8859_2",
		Title:   "Latin-2, Central European",
		Emoji:   "🇩🇪",
		Pangram: "Victor jagt zwölf Boxkämpfer quer über den großen Sylter Deich.",
	},
	{
		ID:      0x3,
		Slug:    "iso_8859_3",
		Title:   "Latin-3, South European",
		Emoji:   "🇲🇹",
		Pangram: "Eble ĉiu kvazaŭ-deca fuŝĥoraĵo ĝojigos homtipon.",
	},
	{
		ID:      0x4,
		Slug:    "iso_8859_4",
		Title:   "Latin-4, North European",
		Emoji:   "🇪🇪",
		Pangram: "See väike mölder jõuab rongile hüpata.",
	},
	{
		ID:      0x5,
		Slug:    "iso_8859_9",
		Title:   "Latin-5, Turkish",
		Emoji:   "🇹🇷",
		Pangram: "Pijamalı hasta yağız şoföre çabucak güvendi.",
	},
	{
		ID:      0x6,
		Slug:    "iso_8859_10",
		Title:   "Latin-6, Nordic",
		Emoji:   "🇳🇴",
		Pangram: "Sær golfer med kølle vant sexquiz på wc i hjemby.",
	},
	{
		ID:      0x7,
		Slug:    "iso_8859_13",
		Title:   "Latin-7, Baltic Rim",
		Emoji:   "🇵🇱",
		Pangram: "Pchnąć w tę łódź jeża lub ośm skrzyń fig.",
	},
	{
		ID:      0x8,
		Slug:    "iso_8859_14",
		Title:   "Latin-8, Celtic",
		Emoji:   "🇮🇪",
		Pangram: "D'ith cat mór dubh na héisc lofa go pras.",
	},
	{
		ID:      0x9,
		Slug:    "iso_8859_15",
		Title:   "Latin-9 (revised Latin-1)",
		Emoji:   "🇵🇹",
		Pangram: "Vejo galã sexy pôr quinze kiwis à força em baú achatado.",
	},
	{
		ID:      0xA,
		Slug:    "iso_8859_16",
		Title:   "Latin-10: South-East European",
		Emoji:   "🇷🇴",
		Pangram: "Încă vând gem, whisky bej și tequila roz, preț fix.",
	},
	{
		ID:      0xB,
		Slug:    "iso_8859_5",
		Title:   "Latin/Cyrillic",
		Emoji:   "🇷🇺",
		Pangram: "В чащах юга жил бы цитрус? Да, но фальшивый экземпляр!",
	},
	{
		ID:      0xC,
		Slug:    "iso_8859_7",
		Title:   "Latin/Greek",
		Emoji:   "🇬🇷",
		Pangram: "Ξεσκεπάζω την ψυχοφθόρα σας βδελυγμία",
	},
	{
		ID:      0xD,
		Slug:    "jis_x0201",
		Title:   "Japanese katakana (halfwidth)",
		Emoji:   "🇯🇵",
		Pangram: "イロハニホヘトチリヌルヲワカヨタレソツネナ",
	},
}
