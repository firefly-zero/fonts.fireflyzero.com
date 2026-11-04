package generator

type Encoding struct {
	id    uint8
	slug  string
	name  string
	emoji string
	chars []rune
}

var encodings = []Encoding{
	{
		id:    0x0,
		slug:  "ascii",
		name:  "ASCII",
		emoji: "🇺🇸",
	},
	{
		id:    0x1,
		slug:  "iso_8859_1",
		name:  "Latin-1, Western European.",
		emoji: "🇵🇹",
	},
	{
		id:    0x2,
		slug:  "iso_8859_2",
		name:  "Latin-2, Central European.",
		emoji: "🇩🇪",
	},
	{
		id:    0x3,
		slug:  "iso_8859_3",
		name:  "Latin-3, South European.",
		emoji: "🇲🇹",
	},
	{
		id:    0x4,
		slug:  "iso_8859_4",
		name:  "Latin-4, North European.",
		emoji: "🇪🇪",
	},
	{
		id:    0x5,
		slug:  "iso_8859_9",
		name:  "Latin-5, Turkish.",
		emoji: "🇹🇷",
	},
	{
		id:    0x6,
		slug:  "iso_8859_10",
		name:  "Latin-6, Nordic.",
		emoji: "🇳🇴",
	},
	{
		id:    0x7,
		slug:  "iso_8859_13",
		name:  "Latin-7, Baltic Rim.",
		emoji: "🇵🇱",
	},
	{
		id:    0x8,
		slug:  "iso_8859_14",
		name:  "Latin-8, Celtic.",
		emoji: "🇮🇪",
	},
	{
		id:    0x9,
		slug:  "iso_8859_15",
		name:  "Latin-9 (revised Latin-1).",
		emoji: "🇵🇹",
	},
	{
		id:    0xA,
		slug:  "iso_8859_16",
		name:  "Latin-10: South-East European.",
		emoji: "🇷🇴",
	},
	{
		id:    0xB,
		slug:  "iso_8859_5",
		name:  "Latin/Cyrillic.",
		emoji: "🇷🇺",
	},
	{
		id:    0xC,
		slug:  "iso_8859_7",
		name:  "Latin/Greek.",
		emoji: "🇬🇷",
	},
	{
		id:    0xD,
		slug:  "jis_x0201",
		name:  "Japanese katakana(halfwidth).",
		emoji: "🇯🇵",
	},
}
