package core

/*
	Quest language parser

	WARN: The parser is not Unicode-aware! Use ASCII characters only!
*/

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/zaklaus/rurik/src/system"
)

// Quest language keywords
const (
	KwTitle      = "title"
	KwBackground = "+background"
	KwBriefing   = "briefing"
	KwResources  = "qrc"
	KwMessage    = "message"
	KwVideo      = "video"
	KwSound      = "sound"
	KwStage      = "stage"
	KwStages     = "qst"
	KwTask       = "task"
	KwEvent      = "event"
	KwSet        = "set"
	KwAbove      = "above"
	KwBelow      = "below"
	KwEquals     = "equals"
	KwNotEquals  = "!equals"
	KwAnd        = "and"
	KwOr         = "or"
	KwXor        = "xor"
	KwComment    = "$-"
	KwScope      = ":"
	KwLeftBrace  = "("
	KwRightBrace = ")"
)

const (
	TkIdentifier = iota
	TkInteger
	TkSeparator
	TkEndOfFile
)

type QuestToken struct {
	Kind  int
	Text  string
	Value int

	WordPos int
}

type QuestTaskDef struct {
	Name           string
	Commands       []QuestCmd
	ProgramCounter int
	IsDone         bool

	IsEvent   bool
	EventArgs []float64
}

type QuestCmd struct {
	Name string
	Args []string
}

type QuestResource struct {
	Kind    int
	Content string
}

const (
	QrMessage = iota
	QrSound
	QrVideo
	QrStage
)

var QuestResourceKinds = map[string]int{
	KwMessage: QrMessage,
	KwSound:   QrSound,
	KwVideo:   QrVideo,
	KwStage:   QrStage,
}

type QuestParser struct {
	Data            []byte
	TextPos         int
	LastWordPos     int
	AllowWhitespace bool
}

func (p *QuestParser) At(idx int) rune {
	return rune(p.Data[idx])
}

func (p *QuestParser) SkipWhitespace() {
	for p.TextPos < len(p.Data) && (IsWhitespace(p.At(p.TextPos))) {
		p.TextPos++
	}
}

func (p *QuestParser) SkipSeparators() {
	for t := p.PeekToken(); t.Kind == TkSeparator; t = p.PeekToken() {
		p.ParseToken()
	}
}

func (p *QuestParser) PeekChar() rune {
	if p.TextPos >= len(p.Data)-1 {
		return 0
	}

	return p.At(p.TextPos)
}

func (p *QuestParser) NextChar() rune {
	r := p.PeekChar()
	p.TextPos++

	return r
}

func (p *QuestParser) ParseToken() QuestToken {
	p.SkipWhitespace()

	if p.TextPos >= len(p.Data)-1 {
		return p.TokenEndOfFile()
	}

	var buf string
	p.LastWordPos = p.TextPos

	al := p.AllowWhitespace
	brc := 0

	if p.TextPos < len(p.Data)-2 &&
		string(p.Data[p.TextPos:p.TextPos+2]) == KwComment {
		for r := p.PeekChar(); r != 0 && r != '\n'; r = p.PeekChar() {
			p.NextChar()
		}

		p.NextChar()
	}

	if string(p.PeekChar()) == KwLeftBrace {
		p.AllowWhitespace = true
	}

	for r := p.PeekChar(); r != 0 && (!IsWhitespace(r) || p.AllowWhitespace) && r != '\n' && string(r) != KwScope; r = p.PeekChar() {
		buf += string(p.NextChar())

		if string(r) == KwLeftBrace {
			brc++
		} else if string(r) == KwRightBrace {
			brc--

			if brc == 0 {
				break
			}
		}
	}

	p.AllowWhitespace = al

	if len(buf) == 0 && string(p.PeekChar()) == KwScope {
		p.NextChar()
		return p.TokenIdentifier(KwScope)
	} else if len(buf) == 0 && p.PeekChar() == '\n' {
		sep := 1
		p.NextChar()
		p.SkipWhitespace()

		for p.PeekChar() == '\n' {
			sep++
			p.NextChar()
			p.SkipWhitespace()
		}

		return p.TokenSeparator(sep)
	}

	if val, err := strconv.Atoi(buf); err == nil {
		return p.TokenInteger(val)
	}

	return p.TokenIdentifier(buf)
}

func (p *QuestParser) NextIdentifier() string {
	p.SkipSeparators()
	ident := p.ParseToken()

	if ident.Kind != TkIdentifier {
		log.Fatalf("Token at '%d' invalid! Expected Identifier.\n", ident.WordPos)
		return ""
	}

	return ident.Text
}

func (p *QuestParser) NextWord() string {
	p.SkipSeparators()
	t := p.ParseToken()

	if t.Kind != TkIdentifier && t.Kind != TkInteger {
		log.Fatalf("Word at '%d' invalid! Expected Word.\n", t.WordPos)
		return ""
	}

	return t.Text
}

func (p *QuestParser) NextString() string {
	p.SkipSeparators()
	p.AllowWhitespace = true
	var buf string

	for tk := p.PeekToken(); tk.Kind != TkEndOfFile && tk.Kind != TkSeparator; tk = p.PeekToken() {
		buf += p.NextWord()
	}

	p.AllowWhitespace = false
	return buf
}

func (p *QuestParser) NextTextBlock() string {
	p.SkipSeparators()
	p.AllowWhitespace = true
	var buf string

	for sep := p.PeekToken(); sep.Kind != TkEndOfFile &&
		sep.Kind != TkSeparator ||
		(sep.Kind == TkSeparator && sep.Value < 2); sep = p.PeekToken() {

		if sep.Kind == TkSeparator {
			buf += "\n"
			p.ParseToken()
		} else {
			buf += p.NextWord()
		}
	}

	p.AllowWhitespace = false
	return buf
}

func (p *QuestParser) NextNumber() int {
	p.SkipSeparators()
	tk := p.ParseToken()

	if tk.Kind != TkInteger {
		log.Fatalf("Number at '%d' invalid! Expected Number.\n", tk.WordPos)
		return -1
	}

	return tk.Value
}

func (p *QuestParser) Expect(ident string) bool {
	p.SkipSeparators()
	ok := true

	tk := p.ParseToken()

	if tk.Kind != TkIdentifier || strings.ToLower(tk.Text) != ident {
		log.Fatalf("Unexpected token '%v'! Expected: '%s'.\n", tk, ident)
		ok = false
	}

	return ok
}

func (p *QuestParser) PeekToken() QuestToken {
	op := *p

	tk := p.ParseToken()

	*p = op
	return tk
}

func (p *QuestParser) CheckResourceKind(kind string) bool {
	_, ok := QuestResourceKinds[strings.ToLower(kind)]
	return ok
}

func (p *QuestParser) ParseResources() (res map[int]QuestResource) {
	res = map[int]QuestResource{}

	p.SkipSeparators()

	for resKind := p.PeekToken(); resKind.Kind != TkEndOfFile && p.CheckResourceKind(resKind.Text); resKind = p.PeekToken() {
		p.ParseToken()
		p.Expect(KwScope)
		resourceID := p.NextNumber()
		kind, _ := QuestResourceKinds[strings.ToLower(resKind.Text)]
		content := p.NextTextBlock()

		res[resourceID] = QuestResource{
			Kind:    kind,
			Content: content,
		}

		p.SkipSeparators()
	}

	return
}

func (p *QuestParser) ParseTasks() (res []QuestTaskDef) {
	res = []QuestTaskDef{}

	// Handle headless main task (entry point)
	res = append(res, QuestTaskDef{
		Name:     "<entry-point>",
		Commands: p.ParseTask(),
	})

	p.SkipSeparators()

	for t := p.PeekToken(); t.Kind == TkIdentifier; t = p.PeekToken() {
		kw := strings.ToLower(p.NextIdentifier())

		if kw != KwTask && kw != KwEvent {
			log.Fatalf("Invalid task found at '%d'!\n", t.WordPos)
			return
		}

		taskName := p.NextIdentifier()
		p.Expect(KwScope)

		res = append(res, QuestTaskDef{
			Name:     taskName,
			Commands: p.ParseTask(),
			IsEvent:  kw == KwEvent,
		})

		taskType := "Task"

		if kw == KwEvent {
			taskType = "Event"
		}

		log.Printf("%s '%s' has been added!", taskType, taskName)

		p.SkipSeparators()
	}

	return
}

func (p *QuestParser) ParseTask() (res []QuestCmd) {
	res = []QuestCmd{}
	p.SkipSeparators()

	for t := p.PeekToken(); t.Kind == TkIdentifier; t = p.PeekToken() {
		// end of the line
		if t.Text == KwTask || t.Text == KwEvent {
			break
		}

		cmd := strings.ToLower(p.NextIdentifier())

		args := []string{}

		for pt := p.PeekToken(); pt.Kind != TkEndOfFile && pt.Kind != TkSeparator; pt = p.PeekToken() {
			args = append(args, p.NextWord())
		}

		res = append(res, QuestCmd{
			Name: cmd,
			Args: args,
		})

		p.SkipSeparators()
	}

	return
}

// QuestDef describes the Quest definition file and the opcodes
type QuestDef struct {
	Title            string
	Briefing         string
	RunsInBackground bool
	Resources        map[int]QuestResource
	TaskDef          []QuestTaskDef
}

var (
	questCache = map[string]*QuestDef{}
)

func ParseQuest(questName string) *QuestDef {
	questAsset := system.FindAsset(fmt.Sprintf("quests/%s.qst", strings.ToLower(questName)))

	if questAsset == nil {
		log.Fatalf("Quest '%s' could not be found!\n", questName)
		return nil
	}

	/* cachedQuest, ok := questCache[questName]

	if ok {
		log.Printf("Reusing existing Quest template '%s'", questName)
		return cachedQuest
	} */

	parser := QuestParser{
		Data: questAsset.Data,
	}

	def := &QuestDef{}

	for t := parser.PeekToken(); t.Kind != TkEndOfFile; t = parser.PeekToken() {
		parser.SkipSeparators()
		ident := parser.NextIdentifier()

		if ident[0] == '+' {
			parser.HandleFlag(def, ident)
			continue
		}

		parser.Expect(KwScope)

		switch strings.ToLower(ident) {
		case KwTitle:
			def.Title = parser.NextString()
		case KwBriefing:
			def.Briefing = parser.NextTextBlock()
		case KwResources:
			def.Resources = parser.ParseResources()
		case KwStages:
			def.TaskDef = parser.ParseTasks()
		default:
			log.Fatalf("Undefined token at '%d'! It says: '%s'.\n", t.WordPos, ident)
			return def
		}
	}

	questCache[questName] = def

	return def
}

func (p *QuestParser) HandleFlag(def *QuestDef, flag string) {
	switch flag {
	case KwBackground:
		def.RunsInBackground = true
	}
}

func IsAlpha(c rune) bool {
	return unicode.IsLetter(c)
}

func IsNumber(c rune) bool {
	return unicode.IsNumber(c)
}

func IsAlphaNumeric(c rune) bool {
	return IsAlpha(c) || IsNumber(c)
}

func IsWhitespace(c rune) bool {
	return unicode.IsSpace(c) && c != '\n'
}

func (p *QuestParser) TokenEndOfFile() QuestToken {
	return QuestToken{
		Kind:    TkEndOfFile,
		WordPos: len(p.Data),
	}
}

func (p *QuestParser) TokenInteger(v int) QuestToken {
	return QuestToken{
		Kind:    TkInteger,
		Text:    strconv.Itoa(v),
		Value:   v,
		WordPos: p.LastWordPos,
	}
}

func (p *QuestParser) TokenIdentifier(s string) QuestToken {
	return QuestToken{
		Kind:    TkIdentifier,
		Text:    s,
		WordPos: p.LastWordPos,
	}
}

func (p *QuestParser) TokenSeparator(sep int) QuestToken {
	return QuestToken{
		Kind:    TkSeparator,
		Value:   sep,
		WordPos: p.LastWordPos,
	}
}
