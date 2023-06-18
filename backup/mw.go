package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func process_lemma(oinl string) (stringList []string) {
	inl := strings.Replace(oinl, " or ", "/", -1)
	aa := []string{inl}
	conti := true
	sum := 0
	for conti {
		sum += 1
		conti = false
		bb := make([]string, 0, 20)
		for _, al := range aa {
			if strings.Contains(al, "/") {
				for _, abl := range strings.Split(al, " ") {
					if strings.Contains(abl, "/") {
						for _, abal := range strings.Split(abl, "/") {
							bb = append(bb, strings.Replace(al, abl, abal, -1))
						}
						break
					}
				}
			} else {
				bb = append(bb, al)
			}
		}
		aa = bb
		for _, al := range aa {
			if strings.Contains(al, "/") {
				if sum < 2 {
					conti = true
				}
			}
		}
	}

	rb := regexp.MustCompile(`\((.*?)\)`)
	sum = 0
	conti = true
	for conti {
		sum += 1
		conti = false
		bb := make([]string, 0, 20)
		for _, al := range aa {
			if strings.Contains(al, "(") && strings.Contains(al, ")") {
				for _, xxs := range rb.FindAllStringSubmatch(al, -1) {
					bb = append(bb, strings.Replace(al, xxs[0], xxs[1], -1))
					bb = append(bb, strings.Replace(al, xxs[0], "", -1))
				}
			} else {
				bb = append(bb, al)
			}
		}
		aa = bb
		for _, al := range aa {
			if strings.Contains(al, "(") && strings.Contains(al, ")") {
				if sum < 2 {
					conti = true
				}
			}
		}
	}
	cc := make([]string, 0, 20)
	for _, al := range aa {
		if al != oinl {
			regx := regexp.MustCompile(" +")
			ac := regx.ReplaceAllString(al, " ")
			ac2 := strings.Trim(ac, " ")
			cc = append(cc, ac2)
		}
	}
	return cc
}

func closego(wg *sync.WaitGroup, chct chan int, xx *int64, name *string) {
	*xx += 1
	wg.Done()
	<-chct
	fmt.Println("finished", *xx, *name)
}

func readXml(wg *sync.WaitGroup, f *os.File, sound_url *map[string]string, fname string, ch chan *string, meow chan string, xx *int64, entryAll *map[string]bool, doneDict *map[string]string, chct chan int) {
	wrapper := `<head><meta charset="utf-8"><link href="mw_now.css"  rel="stylesheet" type="text/css"/><link href="mdd1.css"  rel="stylesheet" type="text/css"/><script src="jquery_mw.js" charset="utf-8" type="text/javascript"></script><script src="mw_now.js" charset="utf-8" type="text/javascript"></script><script>prepare_mw();</script></head><body><karxmw><karxdict></karxdict></karxmw></body>`
	mydoc, _ := goquery.NewDocumentFromReader(strings.NewReader(wrapper))
	reg20 := regexp.MustCompile("\n|\r|/dist-cross-dungarees/images/svg|https://www.merriam-webster.com/dist-cross-dungarees/images/svg")
	reg21 := regexp.MustCompile(" {2,}")
	reg22 := regexp.MustCompile("href ?= ?\"/dictionary/")
	reg23 := regexp.MustCompile("<!--.*?-->")
	name := ""
	defer closego(wg, chct, xx, &name)
	ppp := ""

	path := "raws/" + fname + ".html"
	ofile, err := os.Open(path)
	if err == nil {
		fileinfo, _ := ofile.Stat()
		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)
		ofile.Read(buffer)
		alls := string(buffer)
		ofile.Close()

		alls22 := reg20.ReplaceAllString(alls, "")
		alls23 := reg21.ReplaceAllString(alls22, " ")
		alls24 := reg22.ReplaceAllString(alls23, `href="entry://`)
		alls25 := reg23.ReplaceAllString(alls24, "")
		reg1 := regexp.MustCompile("<body.*</body>")
		result1 := reg1.FindAllStringSubmatch(alls25, -1)
		bodystring := result1[0][0]
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(bodystring))
		mainC := dom.Find("div.left-content.col")
		mainC.Find(".widget more_defs").Each(func(_ int, tag *goquery.Selection) {
			tag.Remove()
		})
		mainC.Find(".time-travel-content-section").Each(func(_ int, tag *goquery.Selection) {
			tag.Remove()
		})
		mainC.Find(".hword").Each(func(_ int, tag *goquery.Selection) {
			tag.Find("span").Each(func(_ int, itag *goquery.Selection) {
				itag.Remove()
			})
			ppn := strings.Trim(tag.Text(), " ")
			if name == "" {
				name = ppn
			} else {
				if !strings.EqualFold(ppn, name) {
					meow <- ""
					if _, ok := (*entryAll)[ppn]; !ok {
						(*entryAll)[ppn] = true
						ppp += ppn + "\r\n@@@LINK=" + name + "\r\n</>\r\n"
						fmt.Println(ppn, "refers to", name)
					}
					<-meow
				}
			}
		})

		regX := regexp.MustCompile(`(href="entry://.*?)"`)
		resultX := regX.FindAllStringSubmatch(bodystring, -1)
		regN := regexp.MustCompile(`href="entry://(.*?)(#.*)`)
		for _, ja := range resultX {
			amatch := ja[1]
			resultN := regN.FindAllStringSubmatch(amatch, -1)
			if len(resultN) > 0 {
				jname := resultN[0][1]
				jpart := resultN[0][2]
				if jname == name {
					bodystring = strings.Replace(bodystring, amatch+"\"", `href="`+jpart+`-anchor"`, 1)
				} else {
					bodystring = strings.Replace(bodystring, amatch+"\"", amatch+"-anchor\"", 1)
				}
			}
		}
		regR := regexp.MustCompile(`(href="#h\d+)"`)
		resultR := regR.FindAllStringSubmatch(bodystring, -1)
		for _, jR := range resultR {
			amatch := jR[1]
			bodystring = strings.Replace(bodystring, amatch+"\"", amatch+"-anchor\"", 1)
		}

		dom, _ = goquery.NewDocumentFromReader(strings.NewReader(bodystring))
		dom.Find("div.mw-def-2020-ad-container, div.mw-mobile-def-2020-ad-container, div.widget.more_defs, .disclaimer, .see-more, .ul-must-login-def").Each(func(_ int, tag *goquery.Selection) {
			tag.Remove()
		})
		mainC = dom.Find("div.left-content.col")
		ddd := mainC.Find(".hword").Slice(0, 1)
		name = ddd.Text()

		meow <- ""
		if hadpath, ok := (*doneDict)[name]; !ok {
			(*doneDict)[name] = path
			<-meow
		} else {
			fmt.Println(name, "repeat! old and new:", hadpath, path)
			<-meow
			return
		}

		for _, vl := range process_lemma(name) {
			meow <- ""
			if _, ok := (*entryAll)[vl]; !ok {
				(*entryAll)[vl] = true
				ppp += vl + "\r\n@@@LINK=" + name + "\r\n</>\r\n"
			}
			<-meow
		}
		mainC.Find(".if, .drp, .ure, .va").Each(func(_ int, tag *goquery.Selection) {
			aurelike := tag.Text()
			meow <- ""
			// (*sound_url)[aurelike] = name
			if _, ok := (*entryAll)[aurelike]; !ok {
				(*entryAll)[aurelike] = true
				ppp += aurelike + "\r\n@@@LINK=" + name + "\r\n</>\r\n"
			}
			<-meow
		})
		mainC.Find("img.lazyload").Each(func(_ int, tag *goquery.Selection) {
			datasrc, _ := tag.Attr("data-src")
			pp := strings.Split(datasrc, "/")
			ppl := pp[len(pp)-1]
			tag.RemoveAttr("data-src")
			tag.RemoveAttr("data-loaded")
			tag.SetAttr("src", "/"+ppl)
		})
		mainC.Find(".play-pron-v2").Each(func(_ int, tag *goquery.Selection) {
			datalang, _ := tag.Attr("data-lang")
			datafile, _ := tag.Attr("data-file")
			datadir, _ := tag.Attr("data-dir")
			tag.RemoveAttr("data-lang")
			tag.RemoveAttr("data-file")
			tag.RemoveAttr("data-dir")
			tag.RemoveAttr("data-url")
			tag.RemoveAttr("data-title")
			tag.RemoveAttr("title")
			voiceurl := "https://media.merriam-webster.com/audio/prons/" + strings.Replace(datalang, "_", "/", -1) + "/mp3/" + datadir + "/" + datafile + ".mp3"
			tag.SetAttr("onclick", `if(!mdd1Existmw()){new Audio("`+voiceurl+`").play();return false;}`)
			tag.SetAttr("href", "sound://sound/"+datafile+".mp3")
			meow <- ""
			(*sound_url)[datafile+".mp3"] = voiceurl
			<-meow
		})

		classMap := make(map[string]bool)
		idList := [6]string{"word-history", "examples", "phrases", "related-phrases", "little-gems", "synonyms"}
		for _, value := range idList {
			classMap[value] = true
		}
		mainC.Children().Each(func(_ int, tag *goquery.Selection) {
			theclass, _ := tag.Attr("class")
			theid, _ := tag.Attr("id")
			if _, ok := classMap[theid]; ok {
				mydoc.Find("karxdict").AppendSelection(tag)
			} else {
				if strings.Split(theclass, " ")[0] == "usage_notes" {
					mydoc.Find("karxdict").AppendSelection(tag)
				} else {
					if len(theid) > 17 {
						if theid[0:17] == "dictionary-entry-" {
							mydoc.Find("karxdict").AppendSelection(tag)
						}
					}
				}
			}
		})
	}
	path2 := "thsr/" + fname + ".html"
	ofile2, err := os.Open(path2)
	if err == nil {
		fileinfo, _ := ofile2.Stat()
		filesize2 := fileinfo.Size()
		buffer2 := make([]byte, filesize2)
		ofile2.Read(buffer2)
		alls_thsr := string(buffer2)
		ofile2.Close()
		alls22 := reg20.ReplaceAllString(alls_thsr, "")
		alls23 := reg21.ReplaceAllString(alls22, " ")
		alls24 := reg22.ReplaceAllString(alls23, `href="entry://`)
		alls25 := reg23.ReplaceAllString(alls24, "")
		reg1 := regexp.MustCompile("<body.*</body>")
		result1 := reg1.FindAllStringSubmatch(alls25, -1)
		bodystring := result1[0][0]
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(bodystring))
		mainT := dom.Find("div.left-content.col")
		wrapper := `<karxthsr></karxthsr>`
		thsrwrap, _ := goquery.NewDocumentFromReader(strings.NewReader(wrapper))
		mydoc.Find("karxmw").AppendSelection(thsrwrap.Find("karxthsr"))
		mainT.Find(".hword").Each(func(_ int, tag *goquery.Selection) {
			ppn := strings.Trim(tag.Text(), " ")
			if name == "" {
				name = ppn
			}
		})

		mainT.Children().Each(func(_ int, tag *goquery.Selection) {
			theid, _ := tag.Attr("id")
			if theid == "faqs" || theid == "related-phrases" {
				mydoc.Find("karxthsr").AppendSelection(tag)
			} else {
				if len(theid) > 16 {
					if theid[0:16] == "thesaurus-entry-" {
						mydoc.Find("karxthsr").AppendSelection(tag)
					}
				}
			}
		})
	}

	s3, _ := mydoc.Html()
	ppp += name + "\r\n" + s3 + "\r\n</>\r\n"
	ch <- &ppp
	f.WriteString(ppp)
	<-ch
}

func main() {
	var wg sync.WaitGroup
	filePath := "test/mw1.html"
	ch := make(chan *string, 1)
	chct := make(chan int, 200)
	meow := make(chan string, 1)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int64 = 0
	entryAll := make(map[string]bool)
	doneDict := make(map[string]string)
	sound_url := make(map[string]string)
	all_inflect := make(map[string]string)

	file, err := os.OpenFile("test/forms-EN.txt", os.O_RDWR, 0666)
	if err == nil {
		buf := bufio.NewReader(file)
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			aabb := strings.Split(line, ": ")
			origin := aabb[0]
			inflections := strings.Split(aabb[1], ", ")
			for _, infl := range inflections {
				all_inflect[infl] = origin
			}
			if err != nil {
				if err == io.EOF {
					fmt.Println("File read ok!")
					break
				} else {
					fmt.Println("Read file error!", err)
					return
				}
			}
		}
		file.Close()
	}

	file, err = os.OpenFile("test/entryrecord.txt", os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		entryAll[line] = true
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
	}
	file.Close()
	files, _ := os.ReadDir(`E:\Golang\mw\raws`)
	for _, file := range files {
		kk := file.Name()

		if kk == "heavy.html" || kk == "happy.html" {
			chct <- 1
			wg.Add(1)
			go readXml(&wg, f, &sound_url, kk[0:len(kk)-5], ch, meow, &count, &entryAll, &doneDict, chct)
		}
	}
	files2, _ := os.ReadDir(`E:\Golang\mw\thsr`)
	for _, file := range files2 {
		kk := file.Name()
		kkpure := kk[:len(kk)-5]

		// if true {
		if _, ok := entryAll[kkpure]; !ok {
			fmt.Println(kkpure)
			chct <- 1
			wg.Add(1)
			go readXml(&wg, f, &sound_url, kk[0:len(kk)-5], ch, meow, &count, &entryAll, &doneDict, chct)
		}
	}
	wg.Wait()
	fmt.Println("final processed", count)
	f.Close()
	fr, _ := os.OpenFile("sound_url.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0777)
	for k, v := range sound_url {
		fr.WriteString(k + "_" + v + "\r\n")
	}
	fr.Close()
}
