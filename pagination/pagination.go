package pagination

import (
	"strconv"
)

type PageItem struct {
	Page	int
	Str		string
}

type Page struct {
	Items	[]PageItem
	MaxPage	int
}

func CreatePage(currentPage, maxPage int) Page {
	page := Page{}
	page.MaxPage = maxPage
	printMaxPage := 7

	// Cals printing page
	startPage := currentPage - (printMaxPage / 2)
	if startPage < 1 {
		startPage = 1
	}
	endPage := currentPage + (printMaxPage /2)
	if endPage > maxPage {
		endPage = maxPage
	}

	// Previos
	if currentPage <= 1 {
		page.Items = append(page.Items, PageItem{0, "<<"})
	} else {
		page.Items = append(page.Items, PageItem{currentPage - 1, "<<"})
	}
	if startPage >= 2 {
		page.Items = append(page.Items, PageItem{1, "1"})
		page.Items = append(page.Items, PageItem{0, "..."})
	}

	// Page
	for i := startPage; i <= endPage; i++ {
		if i == currentPage {
			page.Items = append(page.Items, PageItem{0, strconv.Itoa(i)})
		} else {
			page.Items = append(page.Items, PageItem{i, strconv.Itoa(i)})
		}
	}

	// Next
	if endPage <= (maxPage - (printMaxPage / 2) + 1) {
		page.Items = append(page.Items, PageItem{0, "..."})
		page.Items = append(page.Items, PageItem{maxPage, strconv.Itoa(maxPage)})
	}
	if currentPage >= maxPage {
		page.Items = append(page.Items, PageItem{0, ">>"})
	} else {
		page.Items = append(page.Items, PageItem{currentPage + 1, ">>"})
	}

	return page
}
