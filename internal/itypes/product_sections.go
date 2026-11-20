/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package itypes

import (
	"encoding/xml"
	"sort"
	"strconv"
	"strings"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go-v2/types"
)

const (
	// ovfNS is the OVF envelope namespace used by the ProductSection elements.
	ovfNS = "http://schemas.dmtf.org/ovf/envelope/1"

	// productSectionInfo is the OVF Info of a ProductSection created by the SDK.
	productSectionInfo = "Custom properties"
)

// xmlNode is a generic XML element: name, attributes, children and text are kept
// as they are so that nothing is lost when the document is read, modified and sent back.
type xmlNode struct {
	Name     xml.Name
	Attrs    []xml.Attr
	Children []*xmlNode
	Text     string
}

// isNamespaceDeclaration reports whether the attribute only declares a namespace.
// The encoder declares the namespaces of the elements and attributes it writes.
func isNamespaceDeclaration(a xml.Attr) bool {
	return a.Name.Space == "xmlns" || (a.Name.Space == "" && (a.Name.Local == "xmlns" || strings.HasPrefix(a.Name.Local, "xmlns:")))
}

func (n *xmlNode) unmarshal(d *xml.Decoder, start xml.StartElement) error {
	n.Name = start.Name
	for _, a := range start.Attr {
		if !isNamespaceDeclaration(a) {
			n.Attrs = append(n.Attrs, a)
		}
	}

	var text strings.Builder
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.StartElement:
			child := &xmlNode{}
			if err := child.unmarshal(d, t); err != nil {
				return err
			}
			n.Children = append(n.Children, child)
		case xml.CharData:
			text.Write(t)
		case xml.EndElement:
			// the text of an element that has children is only indentation
			if len(n.Children) == 0 {
				n.Text = text.String()
			}
			return nil
		}
	}
}

func (n *xmlNode) marshal(e *xml.Encoder) error {
	start := xml.StartElement{Name: n.Name, Attr: n.Attrs}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if len(n.Children) == 0 && n.Text != "" {
		if err := e.EncodeToken(xml.CharData(n.Text)); err != nil {
			return err
		}
	}
	for _, child := range n.Children {
		if err := child.marshal(e); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

func (n *xmlNode) is(local string) bool { return n.Name.Local == local }

func (n *xmlNode) attrValue(local string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == local {
			return a.Value
		}
	}
	return ""
}

func (n *xmlNode) setAttr(local, value string) {
	for i, a := range n.Attrs {
		if a.Name.Local == local {
			n.Attrs[i].Value = value
			return
		}
	}
	n.Attrs = append(n.Attrs, xml.Attr{Name: xml.Name{Space: ovfNS, Local: local}, Value: value})
}

func (n *xmlNode) child(local string) *xmlNode {
	for _, c := range n.Children {
		if c.is(local) {
			return c
		}
	}
	return nil
}

func ovfElement(local, text string) *xmlNode {
	return &xmlNode{Name: xml.Name{Space: ovfNS, Local: local}, Text: text}
}

// APIProductSectionList is the XML ProductSectionList of a vApp (GET and PUT of
// /api/vApp/{id}/productSections). The document is kept as a generic tree: an
// update reads the list, changes the properties it manages and sends everything
// else back unchanged (property attributes such as password and qualifiers,
// several Value elements, Category, Product..., Icon, extension elements, several
// ProductSection).
type APIProductSectionList struct {
	root xmlNode
}

// UnmarshalXML decodes the whole ProductSectionList document.
func (l *APIProductSectionList) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	l.root = xmlNode{}
	return l.root.unmarshal(d, start)
}

// MarshalXML encodes the whole ProductSectionList document.
func (l *APIProductSectionList) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	root := l.root
	if root.Name.Local == "" {
		root.Name = xml.Name{Space: vcloudXMLNS, Local: "ProductSectionList"}
	}
	root.Attrs = append([]xml.Attr{{Name: xml.Name{Local: "xmlns:ovf"}, Value: ovfNS}}, root.Attrs...)
	return root.marshal(e)
}

// ToModel converts the XML ProductSectionList to the public guest property model.
func (l *APIProductSectionList) ToModel() types.ModelVAppGuestProperties {
	model := types.ModelVAppGuestProperties{Sections: make([]types.ModelVAppGuestPropertySection, 0, 1)}
	if l == nil {
		return model
	}

	for _, s := range l.root.Children {
		if !s.is("ProductSection") {
			continue
		}

		section := types.ModelVAppGuestPropertySection{
			Class:      s.attrValue("class"),
			Instance:   s.attrValue("instance"),
			Properties: make(map[string]types.ModelVAppGuestProperty),
		}
		for _, p := range s.Children {
			if !p.is("Property") {
				continue
			}

			// The value set by the user is the Value element, the value attribute is the default.
			value := p.attrValue("value")
			if v := currentValue(p); v != nil {
				value = v.attrValue("value")
			}

			key := p.attrValue("key")
			section.Properties[key] = types.ModelVAppGuestProperty{
				Key:              key,
				Label:            childText(p, "Label"),
				Type:             p.attrValue("type"),
				Value:            value,
				UserConfigurable: p.attrValue("userConfigurable") == "true",
			}
		}
		model.Sections = append(model.Sections, section)
	}

	return model
}

func childText(n *xmlNode, local string) string {
	if c := n.child(local); c != nil {
		return c.Text
	}
	return ""
}

// currentValue returns the Value element that is not tied to a deployment configuration, or the first one.
func currentValue(property *xmlNode) *xmlNode {
	var first *xmlNode
	for _, c := range property.Children {
		if !c.is("Value") {
			continue
		}
		if first == nil {
			first = c
		}
		if c.attrValue("configuration") == "" {
			return c
		}
	}
	return first
}

func newProperty(key string, p types.ModelVAppGuestProperty) *xmlNode {
	n := &xmlNode{Name: xml.Name{Space: ovfNS, Local: "Property"}}
	n.setAttr("key", key)
	n.setAttr("value", "")
	n.setAttr("type", propertyType(p.Type))
	n.setAttr("userConfigurable", strconv.FormatBool(p.UserConfigurable))
	n.Children = append(n.Children, ovfElement("Label", labelOr(p.Label, key)))
	value := ovfElement("Value", "")
	value.setAttr("value", p.Value)
	n.Children = append(n.Children, value)
	return n
}

func propertyType(t string) string {
	if t == "" {
		return "string"
	}
	return t
}

func labelOr(label, fallback string) string {
	if label != "" {
		return label
	}
	return fallback
}

// updateProperty applies the modeled fields of a property and leaves the others untouched.
func updateProperty(n *xmlNode, p types.ModelVAppGuestProperty) {
	if p.Type != "" {
		n.setAttr("type", p.Type)
	}
	n.setAttr("userConfigurable", strconv.FormatBool(p.UserConfigurable))

	if p.Label != "" {
		if label := n.child("Label"); label != nil {
			label.Text = p.Label
		} else {
			n.Children = append([]*xmlNode{ovfElement("Label", p.Label)}, n.Children...)
		}
	}

	if v := currentValue(n); v != nil {
		v.setAttr("value", p.Value)
		return
	}
	value := ovfElement("Value", "")
	value.setAttr("value", p.Value)
	n.Children = append(n.Children, value)
}

// ToAPIProductSectionList builds a new ProductSectionList from the public model.
func ToAPIProductSectionList(model types.ModelVAppGuestProperties) *APIProductSectionList {
	return MergeProductSectionList(nil, model)
}

// MergeProductSectionList applies the public model on the current ProductSectionList: the
// properties of the model replace the existing ones (matched by key, in whichever
// section they are), the properties that are not in the model are removed, new
// ones are added to the section with the same class and instance (created when
// missing). Everything the model does not carry is kept as it was read.
func MergeProductSectionList(current *APIProductSectionList, model types.ModelVAppGuestProperties) *APIProductSectionList {
	list := APIProductSectionList{}
	if current != nil {
		list = *current
	}
	if list.root.Name.Local == "" {
		list.root.Name = xml.Name{Space: vcloudXMLNS, Local: "ProductSectionList"}
	}

	type desiredProperty struct {
		property        types.ModelVAppGuestProperty
		class, instance string
	}
	desired := make(map[string]desiredProperty)
	for _, s := range model.Sections {
		for key, p := range s.Properties {
			desired[key] = desiredProperty{property: p, class: s.Class, instance: s.Instance}
		}
	}

	// Existing properties: update the ones that are kept, drop the others.
	seen := make(map[string]bool)
	for _, s := range list.root.Children {
		if !s.is("ProductSection") {
			continue
		}
		kept := s.Children[:0:0]
		for _, c := range s.Children {
			if !c.is("Property") {
				kept = append(kept, c)
				continue
			}
			key := c.attrValue("key")
			d, ok := desired[key]
			if !ok {
				continue
			}
			updateProperty(c, d.property)
			seen[key] = true
			kept = append(kept, c)
		}
		s.Children = kept
	}

	// New properties (sorted by key for a deterministic body).
	keys := make([]string, 0, len(desired))
	for key := range desired {
		if !seen[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	for _, key := range keys {
		d := desired[key]
		section := findSection(&list.root, d.class, d.instance)
		if section == nil {
			section = newSection(d.class, d.instance)
			list.root.Children = append(list.root.Children, section)
		}
		insertProperty(section, newProperty(key, d.property))
	}

	// A model with an empty section and no existing section still sends the section.
	if len(model.Sections) > 0 && findSection(&list.root, model.Sections[0].Class, model.Sections[0].Instance) == nil {
		list.root.Children = append(list.root.Children, newSection(model.Sections[0].Class, model.Sections[0].Instance))
	}

	return &list
}

func findSection(root *xmlNode, class, instance string) *xmlNode {
	for _, s := range root.Children {
		if s.is("ProductSection") && s.attrValue("class") == class && s.attrValue("instance") == instance {
			return s
		}
	}
	return nil
}

func newSection(class, instance string) *xmlNode {
	s := &xmlNode{Name: xml.Name{Space: ovfNS, Local: "ProductSection"}}
	if class != "" {
		s.setAttr("class", class)
	}
	if instance != "" {
		s.setAttr("instance", instance)
	}
	s.Children = append(s.Children, ovfElement("Info", productSectionInfo))
	return s
}

// insertProperty adds the property after the last Property or Category, which keeps
// the extension elements (they come last in the schema) after the properties.
func insertProperty(section, property *xmlNode) {
	at := len(section.Children)
	last := -1
	for i, c := range section.Children {
		if c.is("Property") || c.is("Category") {
			last = i
		}
	}
	if last >= 0 {
		at = last + 1
	} else {
		// no property yet: after the product description elements, before extension elements
		for i, c := range section.Children {
			if !isProductDescription(c.Name.Local) {
				at = i
				break
			}
		}
	}
	section.Children = append(section.Children[:at], append([]*xmlNode{property}, section.Children[at:]...)...)
}

func isProductDescription(local string) bool {
	switch local {
	case "Info", "Product", "Vendor", "Version", "FullVersion", "ProductUrl", "VendorUrl", "AppUrl", "Icon":
		return true
	default:
		return false
	}
}
