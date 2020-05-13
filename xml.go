package gofakeit

import (
	"bytes"
	"encoding/xml"
	"errors"
	"reflect"
)

// XMLOptions defines values needed for json generation
type XMLOptions struct {
	Type          string  `json:"type" xml:"type"` // single or multiple
	RootElement   string  `json:"root_element" xml:"root_element"`
	RecordElement string  `json:"record_element" xml:"record_element"`
	RowCount      int     `json:"row_count" xml:"row_count"`
	Fields        []Field `json:"fields" xml:"fields"`
	Indent        bool    `json:"indent" xml:"indent"`
}

type xmlArray struct {
	XMLName xml.Name
	Array   []xmlMap
}

type xmlMap struct {
	XMLName xml.Name
	Map     map[string]interface{} `xml:",chardata"`
}

type xmlEntry struct {
	XMLName xml.Name
	Value   interface{} `xml:",chardata"`
}

func (m xmlMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m.Map) == 0 {
		return nil
	}

	start.Name = m.XMLName

	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	err = xmlMapLoop(e, &m)
	if err != nil {
		return err
	}

	return e.EncodeToken(start.End())
}

func xmlMapLoop(e *xml.Encoder, m *xmlMap) error {
	var err error
	for key, value := range m.Map {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Bool,
			reflect.String,
			reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			err = e.Encode(xmlEntry{XMLName: xml.Name{Local: key}, Value: value})
			if err != nil {
				return err
			}
		case reflect.Slice:
			e.EncodeToken(xml.StartElement{Name: xml.Name{Local: key}})
			for i := 0; i < v.Len(); i++ {
				err = e.Encode(xmlEntry{XMLName: xml.Name{Local: "value"}, Value: v.Index(i).String()})
				if err != nil {
					return err
				}
			}
			e.EncodeToken(xml.EndElement{Name: xml.Name{Local: key}})
		case reflect.Map:
			err = e.Encode(xmlMap{
				XMLName: xml.Name{Local: key},
				Map:     value.(map[string]interface{}),
			})
			if err != nil {
				return err
			}
		default:
			err = e.Encode(value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// XML generates an object or an array of objects in json format
func XML(xo *XMLOptions) ([]byte, error) {
	// Check to make sure they passed in a type
	if xo.Type != "single" && xo.Type != "array" {
		return nil, errors.New("Invalid type, must be array or object")
	}

	if xo.Fields == nil || len(xo.Fields) <= 0 {
		return nil, errors.New("Must pass fields in order to build json object(s)")
	}

	if xo.RootElement == "" {
		xo.RecordElement = "xml"
	}

	if xo.RecordElement == "" {
		xo.RecordElement = "record"
	}

	if xo.Type == "single" {
		v := xmlMap{
			XMLName: xml.Name{Local: xo.RootElement},
			Map:     make(map[string]interface{}),
		}

		// Loop through fields and add to them to map[string]interface{}
		for _, field := range xo.Fields {
			// Get function info
			funcInfo := GetFuncLookup(field.Function)
			if funcInfo == nil {
				return nil, errors.New("Invalid function, " + field.Function + " does not exist")
			}

			value, err := funcInfo.Call(&field.Params, funcInfo)
			if err != nil {
				return nil, err
			}

			v.Map[field.Name] = value
		}

		// Marshal into bytes
		var b bytes.Buffer
		x := xml.NewEncoder(&b)
		if xo.Indent {
			x.Indent("", "    ")
		}
		err := x.Encode(v)
		if err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	}

	if xo.Type == "array" {
		// Make sure you set a row count
		if xo.RowCount <= 0 {
			return nil, errors.New("Must have row count")
		}

		xa := xmlArray{
			XMLName: xml.Name{Local: xo.RootElement},
			Array:   make([]xmlMap, xo.RowCount),
		}

		for i := 1; i <= int(xo.RowCount); i++ {
			v := xmlMap{
				XMLName: xml.Name{Local: xo.RecordElement},
				Map:     make(map[string]interface{}),
			}

			// Loop through fields and add to them to map[string]interface{}
			for _, field := range xo.Fields {
				if field.Function == "autoincrement" {
					v.Map[field.Name] = i
					continue
				}

				// Get function info
				funcInfo := GetFuncLookup(field.Function)
				if funcInfo == nil {
					return nil, errors.New("Invalid function, " + field.Function + " does not exist")
				}

				value, err := funcInfo.Call(&field.Params, funcInfo)
				if err != nil {
					return nil, err
				}

				v.Map[field.Name] = value
			}

			xa.Array = append(xa.Array, v)
		}

		// Marshal into bytes
		var b bytes.Buffer
		x := xml.NewEncoder(&b)
		if xo.Indent {
			x.Indent("", "    ")
		}
		err := x.Encode(xa)
		if err != nil {
			return nil, err
		}

		return b.Bytes(), nil
	}

	return nil, errors.New("Invalid type, must be array or object")
}
