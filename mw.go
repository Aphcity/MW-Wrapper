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
	// oinl = strings.Replace(oinl, " or ", "/", -1)
	aa := []string{oinl}
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
	if len(cc) == 0 {
		cc = []string{oinl}
	}
	return cc
}

func closego(wg *sync.WaitGroup, chct chan int, xx *int64, name *string) {
	*xx += 1
	wg.Done()
	<-chct
	fmt.Println("finished", *xx, *name)
}

func readXml(wg *sync.WaitGroup, f *os.File, sound_url *map[string]string, fname string, ch chan *string, meow chan string, xx *int64, entry_basic *map[string]bool, all_inflect *map[string]string, chct chan int) {
	fentry := strings.ToLower(strings.Replace(fname, "_", "/", -1))
	for _, vl := range process_lemma(fentry) {
		new_l := strings.ToLower(vl)
		if !strings.EqualFold(new_l, fentry) {
			if _, ok := (*entry_basic)[new_l]; !ok {
				meow <- ""
				(*all_inflect)[new_l] = fentry
				<-meow
			}
		}
	}
	wrapper := `<head><meta charset="utf-8"><link href="mw_now.css"  rel="stylesheet" type="text/css"/><link href="mdd1.css"  rel="stylesheet" type="text/css"/><script src="jquery_mw.js" charset="utf-8" type="text/javascript"></script><script src="mw_now.js" charset="utf-8" type="text/javascript"></script><script>prepare_mw();</script></head><body><karxmw></karxmw></body>`
	mydoc, _ := goquery.NewDocumentFromReader(strings.NewReader(wrapper))
	reg20 := regexp.MustCompile("\n|\r")
	reg21 := regexp.MustCompile(" {2,}")
	reg22 := regexp.MustCompile(`href\s?=\s?"/(dictionary|thesaurus)/`)
	reg225 := regexp.MustCompile(`href\s?=\s?'/(dictionary|thesaurus)/`)
	reg23 := regexp.MustCompile("<!--.*?-->")
	reg1 := regexp.MustCompile("<body.*</body>")

	name := ""
	defer closego(wg, chct, xx, &name)
	name = strings.Replace(fname, "_", "/", -1)
	ppp := ""
	meow <- ""
	var fname_raw string = ""
	if _, ok := (*entry_basic)[name]; !ok {
		// fmt.Println("yes")
		yy := name
		for {
			if xx, ok2 := (*all_inflect)[yy]; ok2 {

				if _, ok3 := (*entry_basic)[xx]; !ok3 {
					fmt.Println(yy, " | ", xx)
					yy = xx
				} else {
					fname_raw = xx
					break
				}
			} else {
				break
			}

		}
		if fname_raw == "" && name[len(name)-1] == 's' {
			fname_raw = name[0 : len(name)-1]
		}
	} else {
		fname_raw = name
	}
	zxz := strings.Replace(fname_raw, "/", "_", -1)

	path := "raws/" + zxz + ".html"
	ofile, err := os.Open(path)
	only_thsr := false
	if err == nil {
		fileinfo, _ := ofile.Stat()
		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)
		ofile.Read(buffer)
		alls := string(buffer)
		ofile.Close()
		<-meow
		alls22 := reg20.ReplaceAllString(alls, "")
		alls23 := reg21.ReplaceAllString(alls22, " ")
		alls24 := reg22.ReplaceAllString(alls23, `href="entry://`)
		alls245 := reg225.ReplaceAllString(alls24, `href='entry://`)
		alls25 := reg23.ReplaceAllString(alls245, "")
		result1 := reg1.FindAllStringSubmatch(alls25, -1)
		bodystring := result1[0][0]
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(bodystring))
		mainC := dom.Find("div.left-content.col")
		wrapper := `<karxdict></karxdict>`
		thsrwrap, _ := goquery.NewDocumentFromReader(strings.NewReader(wrapper))
		mydoc.Find("karxmw").AppendSelection(thsrwrap.Find("karxdict"))

		mainC.Find("h1.hword, h1.hword>.syl").Each(func(_ int, tag *goquery.Selection) {
			// tag.Find("span").Each(func(_ int, itag *goquery.Selection) {
			// 	itag.Remove()
			// })
			ppn := strings.Trim(tag.Text(), " ")
			if len(fname_raw) > 5 {
				if fname_raw[0:6] == "0error" {
					fmt.Println(ppn, "0error")
					fname_raw = strings.ToLower(ppn)
					name = strings.ToLower(ppn)
				}
			}
			if !strings.EqualFold(strings.ToLower(ppn), strings.ToLower(fname_raw)) {
				if _, ok := (*entry_basic)[strings.ToLower(ppn)]; !ok {
					meow <- ""
					(*all_inflect)[strings.ToLower(ppn)] = strings.ToLower(fname_raw)
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
				if jname == fname_raw {
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
		dom.Find("div.mw-def-2020-ad-container, div.mw-mobile-def-2020-ad-container, .disclaimer, .see-more, .ul-must-login-def, div.ref-interlink, .time-travel-content-section").Each(func(_ int, tag *goquery.Selection) {
			tag.Remove()
		})
		mainC = dom.Find("div.left-content.col")

		mainC.Find(".if, .drp, .ure, .va, .mw_t_phrase").Each(func(_ int, tag *goquery.Selection) {
			aurelike := strings.ToLower(tag.Text())
			for _, vl := range process_lemma(aurelike) {
				if !strings.EqualFold(vl, strings.ToLower(fname_raw)) {
					if _, ok := (*entry_basic)[vl]; !ok {
						meow <- ""
						(*all_inflect)[vl] = strings.ToLower(fname_raw)
						// (*sound_url)[vl] = strings.ToLower(fname_raw) //used to retrieve mw_inflections
						<-meow
					}
				}
			}
		})

		mainC.Find("img").Each(func(_ int, tag *goquery.Selection) {
			datasrc, exist := tag.Attr("src")
			if !exist {
				datasrc, _ = tag.Attr("data-src")
			}
			pp := strings.Split(datasrc, "/")
			ppl := pp[len(pp)-1]
			tag.SetAttr("src", "/"+ppl)
		})
		mainC.Find("#synonyms .mw-btn-outline-orange").Each(func(_ int, tag *goquery.Selection) {
			tag.AddClass("myxx")
			tag.RemoveAttr("href")
			tag.SetHtml(`See all Synonyms & Antonyms <img class="ps-1" src="/arrow-right-orange.svg" alt="">`)
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
		idList := []string{"word-history", "examples", "phrases", "related-phrases", "little-gems", "synonyms"}
		for _, value := range idList {
			classMap[value] = true
		}
		mainC.Children().Each(func(_ int, tag *goquery.Selection) {
			theclass, _ := tag.Attr("class")
			theid, exist0 := tag.Attr("id")
			if exist0 {
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
			} else {
				if theclass == "widget more_defs" {
					mydoc.Find("karxdict").AppendSelection(tag)
				}
			}

		})
	} else {
		only_thsr = true
		<-meow
	}
	meow <- ""
	path2 := "thsr/" + fname + ".html"
	ofile2, err := os.Open(path2)
	if err == nil {
		fileinfo, _ := ofile2.Stat()
		filesize2 := fileinfo.Size()
		buffer2 := make([]byte, filesize2)
		ofile2.Read(buffer2)
		alls_thsr := string(buffer2)
		ofile2.Close()
		<-meow
		alls22 := reg20.ReplaceAllString(alls_thsr, "")
		alls23 := reg21.ReplaceAllString(alls22, " ")
		alls24 := reg22.ReplaceAllString(alls23, `href="entry://`)
		alls245 := reg22.ReplaceAllString(alls24, `href='entry://`)
		alls25 := reg23.ReplaceAllString(alls245, "")
		result1 := reg1.FindAllStringSubmatch(alls25, -1)
		bodystring := result1[0][0]
		dom, _ := goquery.NewDocumentFromReader(strings.NewReader(bodystring))
		dom.Find("div.ref-interlink, div.opp-list-scored-container, #faqs .function-label").Each(func(_ int, tag *goquery.Selection) {
			tag.Remove()
		})
		mainT := dom.Find("div.left-content.col")
		mainT.Find("img").Each(func(_ int, tag *goquery.Selection) {
			datasrc, exist := tag.Attr("src")
			if !exist {
				datasrc, _ = tag.Attr("data-src")
			}
			pp := strings.Split(datasrc, "/")
			ppl := pp[len(pp)-1]
			tag.SetAttr("src", "/"+ppl)
		})
		var wrapper string
		if only_thsr {
			wrapper = `<karxthsr></karxthsr>`
		} else {
			wrapper = `<karxthsr style="display: none;"></karxthsr>`
		}
		thsrwrap, _ := goquery.NewDocumentFromReader(strings.NewReader(wrapper))
		mydoc.Find("karxmw").AppendSelection(thsrwrap.Find("karxthsr"))
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
	} else {
		<-meow
	}

	s3, _ := mydoc.Html()
	ppp += name + "\r\n" + s3 + "\r\n</>\r\n"
	ch <- &ppp
	f.WriteString(ppp)
	<-ch
}

func main() {
	var wg sync.WaitGroup
	filePath := "finalOut.html"
	ch := make(chan *string, 1)
	chct := make(chan int, 200)
	meow := make(chan string, 1)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int64 = 0
	entry_basic := make(map[string]bool)
	doneDict := make(map[string]bool)
	sound_url := make(map[string]string)
	all_inflect := make(map[string]string)
	files_raws, _ := os.ReadDir(`E:\Golang\mw\raws`)
	for _, file := range files_raws {
		kk := file.Name()
		kkpure := strings.Replace(kk[:len(kk)-5], "_", "/", -1)
		entry_basic[kkpure] = true
	}
	file, err := os.OpenFile("forms-EN.txt", os.O_RDWR, 0666)
	if err == nil {
		buf := bufio.NewReader(file)
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			aabb := strings.Split(line, ": ")
			origin := strings.ToLower(aabb[0])

			inflections := strings.Split(aabb[1], ", ")
			for _, infl := range inflections {
				infl = strings.ToLower(infl)
				if _, ok := entry_basic[infl]; !ok {
					if infl != origin {
						all_inflect[infl] = origin
					}
				}
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

	file2, err2 := os.OpenFile("mw_inflect.txt", os.O_RDWR, 0666)
	if err2 == nil {
		buf := bufio.NewReader(file2)
		for {
			line, err := buf.ReadString('\n')
			line = strings.TrimSpace(line)
			aabb := strings.Split(line, "_")
			// fmt.Println(aabb)
			if !strings.EqualFold(aabb[0], aabb[1]) {
				if _, ok := entry_basic[strings.ToLower(aabb[0])]; !ok {
					all_inflect[strings.ToLower(aabb[0])] = strings.ToLower(aabb[1])
				}
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
		file2.Close()
	}
	files_thsr, _ := os.ReadDir(`E:\Golang\mw\thsr`)
	for _, file := range files_thsr {
		kk := file.Name()
		kkpure := kk[:len(kk)-5]
		if true {
			// if _, ok := entry_basic[kkpure]; !ok {
			doneDict[kkpure] = true
			chct <- 1
			wg.Add(1)
			go readXml(&wg, f, &sound_url, kkpure, ch, meow, &count, &entry_basic, &all_inflect, chct)
		}
	}
	for _, file := range files_raws {
		kk := file.Name()
		kkpure := kk[:len(kk)-5]
		// if kk == "monster.html" || kk == "0error737.html" {
		if _, ok := doneDict[kkpure]; !ok {
			doneDict[kkpure] = true
			chct <- 1
			wg.Add(1)
			go readXml(&wg, f, &sound_url, kkpure, ch, meow, &count, &entry_basic, &all_inflect, chct)
		}
	}
	wg.Wait()
	fmt.Println("final processed", count)
	for name_from, name_to := range all_inflect {
		if _, ok := doneDict[name_from]; !ok {
			zzz := name_from + "\r\n@@@LINK=" + name_to + "\r\n</>\r\n"
			f.WriteString(zzz)
		}
	}
	f.Close()
	fr, _ := os.OpenFile("sound_Url.txt", os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0777)
	for k, v := range sound_url {
		fr.WriteString(k + "_" + v + "\r\n")
	}
	fr.Close()
}
