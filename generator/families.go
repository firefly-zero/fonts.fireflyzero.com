package generator

type Family struct {
	Name    string
	Author  string
	License string
	URL     string
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
		Author:  "Ines",
		License: "MIT",
		URL:     "https://github.com/the-moonwitch/Cozette",
	},
	"creep": {
		Name:    "creep",
		Author:  "romeovs",
		License: "MIT",
		URL:     "https://github.com/romeovs/creep",
	},
	"ega": {
		Name:    "ega",
		Author:  "IBM",
		License: "",
	},
	"galmuri": {
		Name:    "galmuri",
		Author:  "Lee Minseo",
		License: "OFL-1.1",
		URL:     "https://github.com/quiple/galmuri",
	},
	"haxor": {
		Name:    "haxor",
		Author:  "ucw.cz",
		License: "GPLv2+",
		URL:     "https://github.com/metan-ucw/fonts",
	},
	"helvetica": {
		Name:    "helvetica",
		Author:  "Adobe",
		License: "",
	},
	"ibm-437": {
		Name:    "ibm-437",
		Author:  "IBM",
		License: "",
	},
	"kenney": {
		Name:    "kenney",
		Author:  "kenney.nl",
		License: "CC0",
		URL:     "https://kenney.nl/assets/kenney-fonts",
	},
	"knxt": {
		Name:    "knxt",
		Author:  "Konstantin Bytensky",
		License: "",
	},
	"leovilok": {
		Name:    "leovilok",
		Author:  "lvgx",
		License: "unlicensed",
		URL:     "https://github.com/leovilok/fonts",
	},
	"mda": {
		Name:    "mda",
		Author:  "",
		License: "",
	},
	"monomacs": {
		Name:    "monomacs",
		Author:  "Bradley Jones",
		License: "OFL-1.1",
		URL:     "https://github.com/brdjns-zz/monomacs",
	},
	"peep": {
		Name:    "peep",
		Author:  "Zevv",
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
		URL:     "https://github.com/fcambus/spleen",
	},
	"spleentt": {
		Name:    "spleentt",
		Author:  "Tommy Thorn",
		License: "BSD 2-Clause",
		URL:     "https://github.com/tommythorn/spleentt-5x8-font",
	},
	"tahoma": {
		Name:    "tahoma",
		Author:  "The Wine Project",
		License: "LGPL 2.1",
		URL:     "https://gitlab.winehq.org/wine/wine",
	},
	"tom-thumb": {
		Name:    "tom-thumb",
		Author:  "Raccoon",
		License: "CC0",
		URL:     "https://robey.lag.net/2010/01/23/tiny-monospace-font.html",
	},
	"ucs": {
		Name:    "ucs",
		Author:  "",
		License: "public domain",
		URL:     "https://www.cl.cam.ac.uk/~mgk25/ucs-fonts.html",
	},
	"unifontex": {
		Name:    "unifontex",
		Author:  "stgiga",
		License: "GPL-2.0",
		URL:     "https://github.com/stgiga/UnifontEX",
	},
	"vga": {
		Name:    "vga",
		Author:  "IBM",
		License: "",
	},
	"zx": {
		Name:    "zx-spectrum",
		Author:  "Amstrad plc",
		License: "proprietary",
		URL:     "https://github.com/ha1tch/bdf-fonts/blob/main/sinclair/SINCLAIR.txt",
	},
}
