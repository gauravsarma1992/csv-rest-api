package csvrestapi

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"strings"
)

type (
	CsvFinder struct {
		csvReq   map[string]string `json:"request"`
		Files    []string          `json:"files"`
		Folders  []string          `json:"folders"`
		Response []interface{}     `json:"response"`
	}
)

func NewCsvReader(files, folders []string, csvReq map[string]string) (finder *CsvFinder, err error) {
	finder = &CsvFinder{
		Files:   files,
		Folders: folders,
		csvReq:  csvReq,
	}
	if err = finder.findFilesInFolders(); err != nil {
		return
	}
	return
}

func (finder *CsvFinder) findFilesInFolders() (err error) {
	for _, folder := range finder.Folders {
		var (
			files []fs.FileInfo
		)
		if files, err = ioutil.ReadDir(folder); err != nil {
			log.Println(err)
			continue
		}
		for _, file := range files {
			fileName := file.Name()
			if !strings.HasSuffix(fileName, ".csv") {
				continue
			}
			finder.Files = append(finder.Files, fmt.Sprintf("%s/%s", folder, fileName))
		}
	}
	return
}

func (finder *CsvFinder) GetLines(fileName string) (spLines []string, err error) {
	var (
		fileB       []byte
		fileContent string
	)
	if fileB, err = ioutil.ReadFile(fileName); err != nil {
		return
	}
	fileContent = string(fileB)
	spLines = strings.Split(fileContent, "\n")
	return
}

func (finder *CsvFinder) SearchFile(fileName string) (err error) {
	var (
		headerIdxMap map[string]int
		spLines      []string
	)
	if spLines, err = finder.GetLines(fileName); err != nil {
		return
	}
	if headerIdxMap, err = finder.getHeaderIdxMap(spLines[0]); err != nil {
		return
	}
	if err = finder.compareLines(spLines, headerIdxMap); err != nil {
		return
	}
	return
}

func (finder *CsvFinder) didParamMatch(spElems []string, headerIdxMap map[string]int) (matched bool, err error) {
	for paramKey, paramVal := range finder.csvReq {
		var (
			convertedParam string
			headerIdx      int
			isPresent      bool
		)
		if convertedParam, err = finder.convertToKey(paramKey); err != nil {
			return
		}
		if headerIdx, isPresent = headerIdxMap[convertedParam]; !isPresent {
			err = errors.New(fmt.Sprintf("Header did not match param %s", convertedParam))
			return
		}
		// Encountered a line which doesn't follow the csv format
		if headerIdx >= len(spElems) {
			return
		}
		if paramVal != spElems[headerIdx] {
			return
		}
	}
	matched = true
	return
}

func (finder *CsvFinder) compareLines(spLines []string, headerIdxMap map[string]int) (err error) {
	for idx, line := range spLines {
		var (
			spElems []string
			resp    map[string]string
			matched bool
		)
		if idx == 0 {
			continue
		}
		spElems = strings.Split(line, ",")
		if matched, err = finder.didParamMatch(spElems, headerIdxMap); err != nil {
			log.Println(err)
			continue
		}
		if !matched {
			continue
		}
		if resp, err = finder.convertElemsToJson(spElems, headerIdxMap); err != nil {
			log.Println(err)
			continue
		}
		finder.Response = append(finder.Response, resp)
	}
	return
}

func (finder *CsvFinder) convertElemsToJson(spElems []string, headerIdxMap map[string]int) (resp map[string]string, err error) {
	resp = make(map[string]string)
	for headerKey, headerIdx := range headerIdxMap {
		resp[headerKey] = spElems[headerIdx]
	}
	return
}

func (finder *CsvFinder) getHeaderIdxMap(header string) (headerIdxMap map[string]int, err error) {
	var (
		spHeader []string
	)
	headerIdxMap = make(map[string]int)
	spHeader = strings.Split(header, ",")

	for idx, headerElem := range spHeader {
		var (
			convertedHeaderElem string
		)
		if convertedHeaderElem, err = finder.convertToKey(headerElem); err != nil {
			log.Println(err)
			continue
		}
		headerIdxMap[convertedHeaderElem] = idx
	}
	return
}

func (finder *CsvFinder) convertToKey(name string) (key string, err error) {
	key = strings.TrimSpace(name)
	key = strings.Replace(key, " ", "", -1)
	key = strings.ToLower(key)
	return
}

func (finder *CsvFinder) Search() (resp []interface{}, err error) {
	for _, fileName := range finder.Files {
		if err = finder.SearchFile(fileName); err != nil {
			log.Println(err)
			continue
		}
	}
	resp = finder.Response
	return
}
