package generator

type Family struct {
	Name    string
	Author  string
	License string
}

var families = map[string]Family{
	"bios": {
		Name:    "bios",
		Author:  "IBM",
		License: "",
	},
	"cga": {
		Name:    "cga",
		Author:  "IBM",
		License: "",
	},
	"clean": {
		Name:    "clean",
		Author:  "Dale Schumacher",
		License: "MIT",
	},
	"cozette": {
		Name:    "cozette",
		Author:  "Slavfox",
		License: "",
	},
	"creep": {
		Name:    "creep",
		Author:  "romeovs",
		License: "",
	},
	"ega": {
		Name:    "ega",
		Author:  "IBM",
		License: "",
	},
	"haxor": {
		Name:    "haxor",
		Author:  "ucw.cz",
		License: "",
	},
	"helvetica": {
		Name:    "helvetica",
		Author:  "Adobe",
		License: "",
	},
	"ibm": {
		Name:    "ibm",
		Author:  "IBM",
		License: "",
	},
	"knxt": {
		Name:    "knxt",
		Author:  "Konstantin Bytensky",
		License: "",
	},
	"mda": {
		Name:    "mda",
		Author:  "",
		License: "",
	},
	"peep": {
		Name:    "peep",
		Author:  "Zevv",
		License: "",
	},
	"psevdoazbuka": {
		Name:    "psevdoazbuka",
		Author:  "ucw.cz",
		License: "",
	},
	"raccoon": {
		Name:    "raccoon",
		Author:  "Raccoon",
		License: "",
	},
	"scientifica": {
		Name:    "scientifica",
		Author:  "romeovs",
		License: "",
	},
	"spleen": {
		Name:    "spleen",
		Author:  "Frederic Cambus",
		License: "BSD 2-Clause",
	},
	"tahoma": {
		Name:    "tahoma",
		Author:  "Larry Snyder",
		License: "",
	},
	"ucs": {
		Name:    "ucs",
		Author:  "",
		License: "public domain",
	},
	"vga": {
		Name:    "vga",
		Author:  "IBM",
		License: "",
	},
	"zx": {
		Name:    "zx",
		Author:  "Sinclair",
		License: "",
	},
}
