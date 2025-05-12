package QQMusicDecoder

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"

	"github.com/beevik/etree"
)

var (
	// 使用 regexp2 编译支持前瞻断言的正则
	AmpRegex  = regexp2.MustCompile(`&(?![a-zA-Z]{2,6};|#[0-9]{2,4};)`, regexp2.None)
	QuotRegex = regexp2.MustCompile(
		`(\s+[\w:.-]+\s*=\s*")(([^"]*)((")((?![^"]*\s+[\w:.-]+\s*=\s*"|[\s/]*>))[^"]*)*)"`,
		regexp2.None)
)

/*
public static XmlDocument Create(string content)

	{
	    content = RemoveIllegalContent(content);

	    content = ReplaceAmp(content);

	    var _content = ReplaceQuot(content);

	    var doc = new XmlDocument();

	    try
	    {
	        doc.LoadXml(_content);
	    }
	    catch
	    {
	        doc.LoadXml(content);
	    }

	    return doc;
	}
*/
func CreateXmlDom(content string) *etree.Document {
	content = RemoveIllegalContent(content)

	content = ReplaceAmp(content)

	//content = ReplaceQuot(content)
	doc := etree.NewDocument()

	if err := doc.ReadFromString(content); err != nil {
		fmt.Println("Error loading XML:", err)
	}

	return doc
}

func RemoveIllegalContent(content string) string {
	left := 0
	i := 0
	for i < len(content) {
		if content[i] == '<' {
			left = i
		}

		// 闭区间
		if i > 0 && content[i] == '>' && content[i-1] == '/' {
			part := content[left : i+1]

			// 存在有且只有一个等号
			eqIdx := strings.Index(part, "=")
			if eqIdx != -1 && eqIdx == strings.LastIndex(part, "=") {
				eqPos := left + eqIdx
				part1 := content[left:eqPos]

				if !strings.Contains(strings.TrimSpace(part1), " ") {
					content = content[:left] + content[i+1:]
					i = 0
					continue
				}
			}
		}
		i++
	}
	return strings.TrimSpace(content)
}

func ReplaceAmp(content string) string {
	// 使用 regexp2 的替换方法
	result, _ := AmpRegex.Replace(content, "&amp;", -1, -1)
	return result
}

func ReplaceQuot(content string) string {
	sb := strings.Builder{}
	currentPos := 0

	// 使用 regexp2 的匹配循环
	for {
		match, _ := QuotRegex.FindStringMatchStartingAt(content, currentPos)
		if match == nil {
			break
		}

		group := match.GroupByNumber(2) // 获取第二个捕获组
		_, end := group.Index, group.Index+group.Length

		sb.WriteString(content[currentPos:match.Index])
		sb.WriteString(content[match.Index:group.Index])

		// 替换双引号
		replaced := strings.ReplaceAll(group.String(), `"`, "&quot;")
		sb.WriteString(replaced)

		currentPos = end
	}
	sb.WriteString(content[currentPos:])
	return sb.String()
}

func RecursionFindElement(elem *etree.Element,
	mappingDict map[string]string,
	resDict map[string]*etree.Element) {

	// 检查当前元素是否在映射表中
	if targetKey, exists := mappingDict[elem.Tag]; exists {
		resDict[targetKey] = elem
	}
	/*
	   if (!xmlNode.HasChildNodes)
	   {
	       return;
	   }*/

	if elem.ChildElements() == nil {
		return
	}
	// 递归处理子元素
	for _, child := range elem.ChildElements() {
		RecursionFindElement(child, mappingDict, resDict)
	}
}
