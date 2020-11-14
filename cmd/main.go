package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

const (
	FlagNameSrcDir       = "srcDir"
	FlagNameMatchPattern = "matchPattern"
	FlagNameOutDir       = "outDir"
	FlagFileOrderBy      = "orderBy"
	FlagFileOrder        = "order"

	OrderByName         = "name"
	OrderByDateModified = "date_modified"
	OrderAsc            = "asc"
	OrderDesc           = "desc"

	FlagDefaultOutDir      = "out"
	FlagDefaultFileOrderBy = OrderByName
	FlagDefaultFileOrder   = OrderAsc
	FlagDefaultSrcDir      = "."
)

var srcDirPath = flag.String(FlagNameSrcDir, FlagDefaultSrcDir, "The directory containing the videos to be processed")
var matchFilePattern = flag.String(FlagNameMatchPattern, "", "The regex matching pattern for wanted files")
var outDirPath = flag.String(FlagNameOutDir, FlagDefaultOutDir, "The output directory for processed videos")
var orderBy = flag.String(FlagFileOrderBy, FlagDefaultFileOrderBy,
	fmt.Sprintf("The order criteria of files to concatenate [%s|%s]", OrderByName, OrderByDateModified))
var ordering = flag.String(FlagFileOrder, FlagDefaultFileOrder,
	fmt.Sprintf("The order of files to concatenate [%s|%s]", OrderAsc, OrderDesc))

func main() {

	flag.Parse()

	fmt.Printf("Source dir is: '%s'\nOutput dir is: '%s'\nMatching pattern is: '%s'\n",
		*srcDirPath, *outDirPath, *matchFilePattern)

	reFileName, err := regexp.Compile(*matchFilePattern)
	if err != nil {
		log.Fatalf("invalid regular expression for %S: %v", FlagNameMatchPattern, err)
	}

	files, err := ReadDir(*srcDirPath, *orderBy, *ordering, reFileName)
	if err != nil {
		log.Fatalf("Unable to enumarate %s directory (%s) contents: %v", FlagNameSrcDir, *srcDirPath, err)
	}

	for _, file := range files {
		fmt.Printf("%s\t%v\t%t\n", file.Name(), file.ModTime(), file.IsDir())
	}
}

// ReadDir reads the directory named by dirname and returns
// a list of files sorted by orderBy.
func ReadDir(dirname, orderBy, order string, match *regexp.Regexp) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	//goland:noinspection GoUnhandledErrorResult
	f.Close()
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, errors.New("no nodes in directory")
	}

	var orderFunc func(i, j os.FileInfo, reverse bool) bool
	switch orderBy {
	case OrderByDateModified:
		orderFunc = compareModTimes
	case OrderByName:
		orderFunc = compareNames
	default:
		return nil, errors.New(fmt.Sprintf("unknown %s", FlagFileOrderBy))
	}

	var orderDesc bool
	switch order {
	case OrderAsc:
		orderDesc = false
	case OrderDesc:
		orderDesc = true
	default:
		return nil, errors.New(fmt.Sprintf("unknown %s", FlagFileOrder))
	}

	// remove node if it is a dir or does not match expected name
	// use old school looping because we may reduce size of list as we go
	for i := 0; i < len(list); i++ {
		if list[i].IsDir() || (match != nil && !match.MatchString(list[i].Name())) {
			list[i] = list[len(list)-1]
			list = list[:len(list)-1]
			i-- // force recheck of current i because it has just been swaped by the previously last item
		}
	}

	sort.Slice(list, func(i, j int) bool { return orderFunc(list[i], list[j], orderDesc) })
	return list, nil
}

func compareNames(i, j os.FileInfo, reverse bool) bool {
	if reverse {
		return i.Name() > j.Name()
	}
	return i.Name() < j.Name()
}

func compareModTimes(i, j os.FileInfo, reverse bool) bool {
	if reverse {
		return i.ModTime().After(j.ModTime())
	}
	return i.ModTime().Before(j.ModTime())
}
