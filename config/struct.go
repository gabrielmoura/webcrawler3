package config

type AcceptableMimeTypes struct {
	MimeTypes []string ` json:"mime_types,omitempty" yaml:"mime_types" mapstructure:"mime_types"`
}

type AcceptableSchema struct {
	Schemas []string `json:"schemas,omitempty" yaml:"schemas" mapstructure:"schemas"`
}
type DenySuffixes struct {
	Suffixes []string ` json:"suffixes,omitempty" json:"suffixes" mapstructure:"suffixes"`
}
type LanguageStopWords struct {
	LanguageCode string   ` json:"language_code,omitempty" yaml:"languageCode" mapstructure:"language_code"`
	StopWords    []string ` json:"stop_words,omitempty" yaml:"stopWords" mapstructure:"stop_words"`
}
type CommonStopWords struct {
	LanguageStopWords []*LanguageStopWords ` json:"language_stop_words,omitempty" yaml:"languageStopWords" mapstructure:"language_stop_words"`
}

type MapWordsType map[string]uint32

type ClientHints struct {
	Enabled bool              ` json:"enabled,omitempty" yaml:"enabled" mapstructure:"enabled"`
	Hints   map[string]string ` json:"hints,omitempty" yaml:"hints" mapstructure:"hints"`
}

type CoreConfig struct {
	AcceptableMimeTypes *AcceptableMimeTypes ` json:"acceptable_mime_types,omitempty" yaml:"acceptableMimeTypes" mapstructure:"acceptable_mime_types"`
	AcceptableSchema    *AcceptableSchema    ` json:"acceptable_schema,omitempty" yaml:"acceptableSchema" mapstructure:"acceptable_schema"`
	DenySuffixes        *DenySuffixes        ` json:"deny_suffixes,omitempty" yaml:"denySuffixes" mapstructure:"deny_suffixes"`
	CommonStopWords     *CommonStopWords     ` json:"common_stop_words,omitempty" yaml:"commonStopWords" mapstructure:"common_stop_words"`
	MapWordsType        *MapWordsType        ` json:"map_words_type,omitempty" yaml:"mapWordsType" mapstructure:"map_words_type"`
	ClientHints         *ClientHints         ` json:"client_hints,omitempty" yaml:"clientHints" mapstructure:"client_hints"`
}

func genDefaultConfig() *CoreConfig {
	return &CoreConfig{
		AcceptableMimeTypes: &AcceptableMimeTypes{
			MimeTypes: []string{
				"text/html",
				"text/plain",
				"text/xml",
				"application/xml",
				"application/xhtml+xml",
				"application/rss+xml",
				"application/atom+xml",
				"application/rdf+xml",
				"application/json",
				"application/ld+json",
				"application/vnd.geo+json",
				"application/xml-dtd",
				"application/rss+xml",
				"application/atom+xml",
				"application/rdf+xml",
				"application/json",
				"application/ld+json",
				"application/vnd.geo+json",
			},
		},
		AcceptableSchema: &AcceptableSchema{
			Schemas: []string{
				"http",
				"https",
				"",
			},
		},
		DenySuffixes: &DenySuffixes{
			Suffixes: []string{
				".css",
				".js",
				".png",
				".jpg",
				".jpeg",
				".gif",
				".svg",
				".ico",
				".mp4",
				".mp3",
				".avi",
				".flv",
				".mpeg",
				".webp",
				".webm",
				".woff",
				".woff2",
				".ttf",
				".eot",
				".otf",
				".pdf",
				".zip",
				".tar",
				".gz",
				".bz2",
				".xz",
				".7z",
				".rar",
				".apk",
				".exe",
				".dmg",
				".img",
				".pdf",
			},
		},
		CommonStopWords: &CommonStopWords{
			LanguageStopWords: []*LanguageStopWords{
				{
					LanguageCode: "en",
					StopWords:    []string{"is", "or", "a", "and", "the", "are", "of", "to"},
				},
				{
					LanguageCode: "pt",
					StopWords: []string{"a", "à", "ao", "aos", "aquela", "aquelas", "aquele", "aqueles", "aquilo", "as", "até",
						"com", "como", "da", "das", "de", "dela", "delas", "dele", "deles", "desta", "destas",
						"deste", "destes", "do", "dos", "e", "é", "ela", "elas", "ele", "eles", "em", "entre",
						"era", "essa", "essas", "esse", "esses", "esta", "estamos", "estas", "estava", "estavam",
						"estávamos", "este", "estes", "estou", "eu", "foi", "fomos", "for", "foram", "havia",
						"isso", "isto", "já", "lhe", "lhes", "mais", "mas", "me", "mesmo", "meu", "meus", "minha",
						"minhas", "muito", "na", "nas", "nem", "no", "nos", "nossa", "nossas", "nosso", "nossos",
						"num", "numa", "o", "os", "ou", "para", "pela", "pelas", "pelo", "pelos", "por", "qual",
						"quando", "que", "quem", "se", "sem", "seu", "seus", "sua", "suas", "só", "também", "te",
						"tem", "temos", "tenha", "tenham", "teu", "teus", "teve", "tinha", "tinham", "tua", "tuas",
						"um", "uma", "você", "vocês", "vos"},
				},
				{
					LanguageCode: "ru",
					StopWords: []string{"а", "без", "более", "бы", "был", "была", "были", "было", "быть", "в", "вам", "вас", "всё",
						"все", "всего", "всех", "вы", "где", "да", "даже", "для", "до", "его", "ее", "ей", "ему",
						"если", "ест", "есть", "ещё", "ж", "же", "за", "здесь", "и", "из", "или", "им", "их", "к",
						"как", "какая", "какой", "когда", "кое", "кто", "куда", "ли", "либо", "мне", "много", "может",
						"можно", "мой", "моя", "мы", "на", "над", "надо", "наш", "не", "него", "нее", "нет", "ни",
						"них", "но", "ну", "о", "об", "однако", "он", "она", "они", "оно", "от", "очень", "по",
						"под", "после", "при", "про", "с", "сам", "сама", "сами", "само", "свое", "своего", "своей",
						"свои", "себе", "себя", "сейчас", "со", "совсем", "так", "такой", "там", "тебя", "тем",
						"теперь", "то", "тогда", "того", "тоже", "только", "том", "ты", "у", "уже", "хотя", "чего",
						"чей", "чем", "что", "чтобы", "чуть", "эта", "эти", "это", "этого", "этой", "этом", "эту",
						"я"},
				},
				{
					LanguageCode: "es",
					StopWords: []string{"el", "la", "los", "las", "de", "del", "y", "a", "en", "un", "una", "unos", "unas", "con",
						"para", "por", "su", "se", "que", "es", "soy", "eres", "somos", "son", "me", "te", "nos", "le",
						"les", "lo", "mi", "tu", "si", "no", "pero", "porque", "como", "esta", "estoy", "estas", "estamos",
						"estais", "estan", "muy", "poco", "mucho", "todo", "todos", "al", "algo", "alguien", "donde", "cuando",
						"como", "aqui", "ahi", "alli", "ahora", "antes", "despues", "hoy", "ayer", "mañana", "siempre", "nunca"},
				},
				{
					LanguageCode: "hindi",
					StopWords: []string{"का", "के", "की", "में", "है", "और", "यह", "वह", "से", "को", "पर", "इस", "होता", "ही", "हैं", "ये", "वो", "कर", "गया", "लिए",
						"अपना", "अपनी", "अपने", "कुछ", "थी", "थे", "थीं", "हुआ", "जा", "रहा", "रहे", "जाता", "जाती", "जाते", "एक", "दो", "तीन", "चार",
						"पांच", "छह", "सात", "आठ", "नौ", "दस"},
				},
				{
					LanguageCode: "ch",
					StopWords: []string{"的", "了", "在", "是", "我", "有", "和", "就", "不", "人", "这", "那", "中", "来", "上", "大", "为", "个", "国",
						"以", "说", "到", "要", "子", "你", "会", "着", "能", "里", "去", "年", "得", "他", "她", "它", "们", "地", "也",
						"自", "这", "时", "那", "儿", "可", "就", "给", "下", "都", "向", "看", "起", "还", "过", "只", "把", "对", "做",
						"当", "想", "成", "事", "被", "用", "多", "从", "面", "等", "前", "些", "于", "后", "所", "又", "经", "方", "现",
						"没", "吧", "定", "得", "该", "好好", "家", "种", "那", "里", "然", "其", "间", "什", "么", "很", "得", "哪", "些",
						"向", "生", "里", "果", "再", "两", "并", "而", "些", "定"},
				},
			},
		},
		ClientHints: &ClientHints{
			Enabled: false,
			Hints: map[string]string{
				"sec-ch-ua":          `"Chromium";v="128", "Not;A=Brand";v="24", "Crl";v="128"`,
				"sec-ch-ua-mobile":   "?0",
				"sec-ch-ua-platform": `"Linux"`,
				"sec-fetch-dest":     "document",
				"sec-fetch-mode":     "navigate",
				"sec-fetch-site":     "none",
				"sec-fetch-user":     "?1",
				"sec-gpc":            "1",
			},
		},
	}
}
