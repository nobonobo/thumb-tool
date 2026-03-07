package main

import (
	"bytes"
	"encoding/xml"
	"log"
	"os"
	"strings"
)

func patchSVG(input, output string) error {
	outputFile, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	template, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	patched, _, err := replaceTextSpansByID(template, config.params)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := outputFile.Write(patched); err != nil {
		log.Fatal(err)
	}
	if err := outputFile.Sync(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// XPathライク: id="xxx"配下のspanテキストを置換（regexpなし）
func replaceTextSpansByID(data []byte, replacements map[string]string) ([]byte, int, error) {
	dec := xml.NewDecoder(bytes.NewReader(data))
	buf := bytes.NewBuffer(nil)
	enc := xml.NewEncoder(buf)

	var stack []xml.StartElement
	targetID := ""
	processed := 0

	for {
		t, err := dec.Token()
		if err != nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			stack = append(stack, se)

			if se.Name.Local == "svg" {
				fix := []xml.Attr{}
				for _, attr := range se.Attr {
					if attr.Name.Local == "xmlns" && attr.Name.Space == "" {
						continue
					}
					fix = append(fix, attr)
				}
				se.Attr = fix
			}
			// id属性でtarget設定
			if id := getAttr(se.Attr, "id"); id != "" {
				if _, exists := replacements[id]; exists {
					targetID = id
				}
			}

			// targetID下のtspan/span開始 → 内容スキップ
			if targetID != "" && (isTextSpan(se.Name.Local)) {
				enc.EncodeToken(se) // 開始タグ出力

				// 終了までスキップし、置換テキスト挿入
				skipToEndElement(dec, enc, se.Name, replacements[targetID])
				processed++
				continue
			}

			enc.EncodeToken(se)

		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}

			// targetID終了でフラグリセット
			if targetID != "" {
				targetID = ""
			}

			enc.EncodeToken(se)

		default:
			enc.EncodeToken(t)
		}
	}

	enc.Flush()
	return buf.Bytes(), processed, nil
}

// tspanまたはspanか判定
func isTextSpan(localName string) bool {
	return strings.ToLower(localName) == "tspan" ||
		strings.ToLower(localName) == "span"
}

// 要素終了までスキップし、置換テキストを挿入
func skipToEndElement(dec *xml.Decoder, enc *xml.Encoder, startName xml.Name, replacementText string) {
	// 置換テキスト出力
	enc.EncodeToken(xml.CharData([]byte(replacementText)))

	// 終了タグまでスキップ
	for {
		t, err := dec.Token()
		if err != nil {
			return
		}

		if ee, ok := t.(xml.EndElement); ok && ee.Name == startName {
			enc.EncodeToken(ee)
			return
		}
	}
}

func getAttr(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return ""
}
