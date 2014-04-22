package parse

import (
	"errors"
	"github.com/jeidsath/unigreek"
	"strings"
)

func (root *Element) findInChildren(name string) (*Element, error) {
	for _, ee := range root.Children {
		if ee.Name() == name {
			return ee, nil

		}
	}
	return &Element{}, errors.New("findInChildren: " + name)
}

func (root *Element) lookup(loc string) (*Element, error) {
	locSplit := strings.Split(loc, ">")
	if len(locSplit) > 1 {
		ee, err := root.findInChildren(locSplit[0])
		if err != nil {
			return &Element{}, err
		}
		return ee.lookup(strings.Join(locSplit[1:], ">"))
	} else {
		return root.findInChildren(locSplit[0])
	}
}

func (ee *Element) lookupStrings(locations map[string]string) (map[string]string, error) {
	output := map[string]string{}

	for kk, vv := range locations {
		cd, err := ee.lookup(vv + ">CHARDATA")
		if err != nil {
			return map[string]string{}, err
		}
		output[kk] = cd.Content
	}
	return output, nil
}

func (root *Element) rootToText() ([]string, error) {
	output := []string{}

	monoloc := "TEI.2>teiHeader>fileDesc>sourceDesc>biblStruct>monogr"
	monogr, err := root.lookup(monoloc)
	if err != nil {
		return []string{}, err
	}

	lookupMap := map[string]string{
		"author":    "author",
		"title":     "title",
		"publisher": "imprint>publisher",
		"date":      "imprint>date",
	}

	monoMap, err := monogr.lookupStrings(lookupMap)
	if err != nil {
		return []string{}, err
	}

	body, err := root.lookup("TEI.2>text>body")
	if err != nil {
		return []string{}, err
	}

	bodyStrings, err := body.bodyToText()
	if err != nil {
		return []string{}, err
	}

	output = append(output, monoMap["title"])
	output = append(output, "By "+monoMap["author"])
	output = append(output, "\n")
	output = append(output, monoMap["publisher"])
	output = append(output, monoMap["date"])
	output = append(output, bodyStrings...)

	return output, nil
}

func (body *Element) bodyToText() ([]string, error) {
	output := []string{}
	line := ""
	bps, err := body.bodyChildren()
	if err != nil {
		return output, err
	}
	for _, ll := range bps {
		switch ll := ll.(type) {
		case string:
			if line != "" {
				line += " "
			}
			line += ll
		case milestone:
                        // Paragraph
			if ll.Ed == "P" {
				output = append(output, line)
				line = "    "
			}
                case newline:
			output = append(output, line)
                        line = ""
		}
	}
	if line != "" {
		output = append(output, line)
	}
	return output, nil
}

//string, newline
type bodyParts interface {
}

type milestone struct {
	N    string
	Unit string
	Ed   string
}

type newline struct {
}

func replaceEscapes(input string) string {
	escapes := map[string]string{
		"&lsqb;": "[",
		"&rsqb;": "]",
		"&lpar;": "(",
		"&rpar;": ")",
	}

	output := input
	for kk, vv := range escapes {
		output = strings.Replace(output, kk, vv, -1)
	}
	return output
}

func trim(input string) string {
	return strings.Trim(input, " \n")
}

func prepareGreek(input string) (string, error) {
	grk, err := unigreek.Convert(input)
	if err != nil {
		return "", err
	}
	grk = replaceEscapes(grk)
	grk = trim(grk)
	return grk, nil
}

func (ee *Element) childParts() ([]bodyParts, error) {
	output := []bodyParts{}
	for _, child := range ee.Children {
		out, err := child.bodyChildren()
		if err != nil {
			return []bodyParts{}, err
		}
		output = append(output, out...)
	}
	return output, nil
}

func (ee *Element) bodyChildren() ([]bodyParts, error) {
	switch ee.Name() {
	case "CHARDATA":
		grk, err := prepareGreek(ee.Content)
		return []bodyParts{grk}, err
	case "milestone":
		var n, unit, ed string
		for _, attr := range ee.Attrs {
			switch attr.Name.Local {
			case "n":
				n = attr.Value
			case "unit":
				unit = attr.Value
			case "ed":
				ed = attr.Value
			default:
				return []bodyParts{}, errors.New("unknown milestone attribute")
			}
		}
		return []bodyParts{milestone{n, unit, ed}}, nil
        case "head":
                out := []bodyParts{newline{}, newline{}}
                additional, err := ee.childParts()
                out = append(out, additional...)
                out = append(out, newline{})
                return out, err
	default:
		return ee.childParts()
	}
}

func (ee *Element) Document() ([]string, error) {
	return ee.rootToText()
}

func (ee *Element) ToText() ([]string, error) {
	output := []string{}

	// Itself
	switch ee.Name() {
	default:
		output = append(output, ee.sAttrs())
		texts, err := ee.childrenToText()
		if err != nil {
			return []string{}, err
		}
		output = append(output, texts...)
	case "CHARDATA":
		output = append(output, ee.sAttrs())
		ss, err := ee.charDataText()
		if err != nil {
			return []string{}, err
		}
		if ss != "" {
			output = append(output, ss)
		}
	}

	return output, nil
}

func (ee *Element) childrenToText() ([]string, error) {
	output := []string{}
	for _, cc := range ee.Children {
		texts, err := cc.ToText()
		if err != nil {
			return []string{}, err
		}
		for _, ss := range texts {
			output = append(output, "  "+ss)
		}
	}
	return output, nil
}

func (ee *Element) sAttrs() string {
	output := ee.HierarchyString()
	for _, aa := range ee.Attrs {
		output += " " + aa.Name.Local + "=" + aa.Value
	}
	return output
}

func (ee *Element) charDataText() (string, error) {
	ss := trim(ee.Content)
	if ss == "" {
		return ss, nil
	}

	switch ee.Parent.Name() {
	default:
		return "", errors.New("Unknown language for <" + ee.Parent.Name() + ">")
	case "title", "author", "titleStmt", "extent", "fileDesc", "biblStruct", "monogr", "publisher", "date", "language", "name", "resp", "item":
		out := trim(ee.Content)
		out = replaceEscapes(out)
		return out, nil
	case "p", "head", "quote":
		out, err := unigreek.Convert(ee.Content)
		if err != nil {
			return "", err
		}
		out = replaceEscapes(out)
		return out, nil
	}
}
