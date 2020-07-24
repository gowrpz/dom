package dom

import "github.com/PuerkitoBio/goquery"

type Element struct {
	el *goquery.Selection
}

func (e *Element) SelectElements(selector string) []*Element {
	var elements []*Element

	e.el.Find(selector).Each(func(i int, selection *goquery.Selection) {
		elements = append(elements, &Element{el: selection})
	})

	return elements
}

func (e *Element) Select(selector string) *Element {
	return &Element{
		el: e.el.Find(selector),
	}
}

func (e *Element) GetAttr(attrName string) (string, bool) {
	return e.el.Attr(attrName)
}

func (e *Element) GetText() string {
	return e.el.Text()
}