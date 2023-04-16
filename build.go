package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
)

// From: https://stackoverflow.com/questions/51779243/copy-a-folder-in-go

func CopyDirectory(scrDir, dest string) error {
	entries, err := os.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if err := Copy(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Copy(srcFile, dstFile string) error {
	out, err := os.Create(dstFile)
	if err != nil {
		return err
	}

	defer out.Close()

	in, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}

var (
	build  = flag.Bool("build", false, "Build a static version of the site for Github Pages.")
	output = flag.String("output", "_site", "Output folder for static site version.")
	listen = flag.String("listen", ":3000", "Hostname to listen for dynamic version.")
)

type Scenario struct {
	Title         string   `yaml:"title"`
	Date          string   `yaml:"date"`
	Description   string   `yaml:"description"`
	Url           string   `yaml:"url"`
	Tags          []string `yaml:"tags"`
	Content       template.HTML
	Code          string
	FormattedCode template.HTML
	Page          string
}

func buildPage(filename string, obj any, out io.Writer) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	tpl := template.New(filename)

	tpl, err = tpl.Parse(string(content))
	if err != nil {
		return err
	}

	return tpl.Execute(out, obj)
}

func formatCode(code []byte) (string, error) {
	style := styles.Get("dracula")
	if style == nil {
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(false))

	lexer := lexers.Get("python")

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)

	err = formatter.Format(buf, style, iterator)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func buildScenario(filename string) (Scenario, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return Scenario{}, nil
	}

	var matter Scenario

	rest, err := frontmatter.Parse(bytes.NewReader(content), &matter)
	if err != nil {
		return Scenario{}, nil
	}

	code, err := os.ReadFile(strings.TrimSuffix(filename, ".md") + ".star")
	if err != nil {
		return Scenario{}, nil
	}

	matter.Code = string(code)

	formattedCode, err := formatCode(code)
	if err != nil {
		return Scenario{}, nil
	}

	matter.FormattedCode = template.HTML(formattedCode)

	matter.Content = template.HTML(markdown.ToHTML(rest, nil, nil))

	buf := new(bytes.Buffer)

	err = buildPage("pages/scenario.tpl.html", matter, buf)
	if err != nil {
		return Scenario{}, nil
	}

	matter.Page = buf.String()

	return matter, nil
}

func buildScenarioADay(basePath string) ([]Scenario, error) {
	var ret []Scenario

	years, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	for _, year := range years {
		months, err := os.ReadDir(filepath.Join(basePath, year.Name()))
		if err != nil {
			return nil, err
		}
		for _, month := range months {
			days, err := os.ReadDir(filepath.Join(
				basePath, year.Name(), month.Name()))
			if err != nil {
				return nil, err
			}
			for _, day := range days {
				if strings.HasSuffix(day.Name(), ".md") {
					scenario, err := buildScenario(filepath.Join(
						basePath, year.Name(), month.Name(), day.Name()))
					if err != nil {
						return nil, err
					}

					ret = append(ret, scenario)
				}
			}
		}
	}

	return ret, nil
}

func main() {
	flag.Parse()

	if *build {
		scenarios, err := buildScenarioADay("scenarioADay")
		if err != nil {
			log.Fatal(err)
		}

		indexFile, err := os.Create(filepath.Join(*output, "index.html"))
		if err != nil {
			log.Fatal(err)
		}
		defer indexFile.Close()

		err = buildPage("pages/index.tpl.html", struct{ Scenarios []Scenario }{
			Scenarios: scenarios,
		}, indexFile)
		if err != nil {
			log.Fatal(err)
		}

		for _, scenario := range scenarios {
			outputPath := filepath.Join(*output, scenario.Url)

			err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			f, err := os.Create(outputPath)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			_, err = f.Write([]byte(scenario.Page))
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.MkdirAll(filepath.Join(*output, "static"), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		err = CopyDirectory("static/", filepath.Join(*output, "static"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		mux := http.NewServeMux()

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			scenarios, err := buildScenarioADay("scenarioADay")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Internal Error: %s", err)
				return
			}

			err = buildPage("pages/index.tpl.html", struct{ Scenarios []Scenario }{
				Scenarios: scenarios,
			}, w)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Internal Error: %s", err)
				return
			}
		})

		mux.HandleFunc("/scenarioaday/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("getting scenario with name: %s", r.URL.Path)

			scenarios, err := buildScenarioADay("scenarioADay")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Internal Error: %s", err)
				return
			}

			for _, scenario := range scenarios {
				if r.URL.Path == scenario.Url {
					_, err := w.Write([]byte(scenario.Page))
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						fmt.Fprintf(w, "Internal Error: %s", err)
						return
					}

					return
				}
			}

			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "not found")
		})

		mux.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
			filename := strings.TrimPrefix(r.URL.Path, "/static/")
			log.Printf("got static request for %s", filename)
			http.ServeFile(w, r, filepath.Join("static", filename))
		})

		log.Printf("Listening on %s", *listen)

		err := http.ListenAndServe(*listen, mux)
		if err != nil {
			log.Fatal(err)
		}
	}
}
