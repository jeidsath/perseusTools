package parse

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"strings"
)

var entityMap = map[string]string{
	"responsibility":  "Responsibility",
	"fund.AnnCPB":     "fund.AnnCPB",
	"Perseus.publish": "Perseus.publish",
	"lsqb":            "&lsqb;",
	"rsqb":            "&rsqb;",
	"lpar":            "&lpar;",
	"rpar":            "&rpar;",
	"dagger":          "†",
	"mdash":           "—",
	"stigma":          "Ϛ",
	"ldquo":           "“",
	"rdquo":           "”",
	"lsquo":           "‘",
	"rsquo":           "’",
}

type Element struct {
	Hierarchy []string
	Attrs     []xml.Attr
	Content   string
	Children  []*Element
	Parent    *Element
}

func (*Element) ToHtml() []string {
	return []string{}
}

func (ee *Element) HierarchyString() string {
	return strings.Join(ee.Hierarchy, ">")
}

func (ee *Element) Name() string {
	if len(ee.Hierarchy) == 0 {
		return ""
	}
	return ee.Hierarchy[len(ee.Hierarchy)-1]
}

func addElement(ee *Element, tk *xml.StartElement) *Element {
	hier := make([]string, len(ee.Hierarchy)+1)
	copy(hier, ee.Hierarchy)
	hier[len(ee.Hierarchy)] = tk.Name.Local
	newElement := &Element{hier, tk.Attr, "", [](*Element){}, ee}
	ee.Children = append(ee.Children, newElement)
	return newElement
}

func endElement(ee *Element, tk *xml.EndElement) (*Element, error) {
	if ee.Name() != tk.Name.Local {
		return &Element{}, errors.New("endElement: unmatched tag")
	}
	return ee.Parent, nil
}

func addCharData(ee *Element, tk *xml.CharData) {
	hier := make([]string, len(ee.Hierarchy)+1)
	copy(hier, ee.Hierarchy)
	hier[len(ee.Hierarchy)] = "CHARDATA"
	newElement := &Element{hier, []xml.Attr{}, string(*tk), [](*Element){}, ee}
	ee.Children = append(ee.Children, newElement)
}

func GetTextables(input []byte) (*Element, error) {
	decoder := xml.NewDecoder(bytes.NewBuffer(input))
	decoder.Entity = entityMap

	rootHier := []string{"ROOT"}
	root := &Element{rootHier, []xml.Attr{}, "", [](*Element){}, &Element{}}
	tip := root

	for {
		token, err := decoder.RawToken()

		if err == io.EOF {
			break
		}

		if err != nil {
			return &Element{}, err
		}

		// StartElement, EndElement, CharData, Comment, ProcInst, or Directive
		switch token.(type) {
		default:
			return &Element{}, errors.New("getTextables: unexpected token")
		case xml.StartElement:
			tk := token.(xml.StartElement)
			tip = addElement(tip, &tk)
		case xml.EndElement:
			tk := token.(xml.EndElement)
			tip, err = endElement(tip, &tk)
			if err != nil {
				return &Element{}, err
			}
		case xml.CharData:
			tk := token.(xml.CharData)
			addCharData(tip, &tk)
		case xml.Comment:
		case xml.ProcInst:
		case xml.Directive:
		}

	}

	return root, nil
}
