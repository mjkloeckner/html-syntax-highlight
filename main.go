package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/styles"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/formatters/html"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"fmt"
)

// print CSS needed for colorizing the code parts
func printCSS(w io.Writer, formatter, style string) {
	f := html.New(html.WithClasses(true), html.TabWidth(4))

	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	println(f.WriteCSS(w, s))
}

// Highlight some text.
// Lexer, formatter and style may be empty, in which case a best-effort is made.
func Highlight(w io.Writer, source, lexer, formatter, style string) error {
	// Determine lexer.
	l := lexers.Get(lexer)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	f := html.New(html.WithClasses(true), html.TabWidth(4))

	// Determine style.
	s := styles.Get(style)
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return f.Format(w, s, it)
}

func replaceCodeParts(mdFile []byte) (string, error) {
	byteReader := bytes.NewReader(mdFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile("(<pre><code class=\"language-).*?(\">)")
	rp := strings.NewReplacer("<pre><code class=\"language-", "", "\">", "")

	langs := re.FindAllString(string(mdFile), -1)

	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("pre").Each(func(i int, s *goquery.Selection) {
		lang := rp.Replace(string(langs[i]))

		if lang == "console" {
			lang = "bash"
		}

		buf := new(bytes.Buffer)

		err = Highlight(buf, s.Text(), lang, "html", "monokai")
		if err != nil {
			log.Fatal(err)
		}

		s.ReplaceWithHtml(string(buf.String()))
	})
	new, err := doc.Html()
	if err != nil {
		return "", err
	}

	// remove unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)

	return new, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("No file name provided")
	}

	args := os.Args[1:]
	html, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	replaced, err := replaceCodeParts(html)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(replaced)
	// printCSS(os.Stdout, "html", "gruvbox")
}
