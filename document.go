package dom

import "github.com/PuerkitoBio/goquery"

type Document struct {
	docEl *goquery.Document
}

func (d *Document) SelectElements(selector string) []*Element {
	var elements []*Element

	d.docEl.Find(selector).Each(func(i int, selection *goquery.Selection) {
		elements = append(elements, &Element{el: selection})
	})

	return elements
}

func (d *Document) Select(selector string) *Element {
	return &Element{
		el: d.docEl.Find(selector),
	}
}
