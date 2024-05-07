package web

import "strings"

// 实现匹配动态路由
// url要么为空，要么存放的是已经注册的路由
type node struct {
	url      string  // 待匹配的路由， 例如/p/:lang, 非子节点为空， 真实的url
	part     string  // 路由中的一部分，例如 :lang。  可能为* 或者：，此时会出现模糊匹配
	children []*node // 子节点
	isWild   bool    // 是否精确匹配, 如果含* 或者：，true 则为模糊匹配，存在* 和 ：
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 遍历到height层, parts去除‘/’后的 路径
func (n *node) insert(url string, parts []string, height int) {
	if len(parts) == height {
		n.url = url
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(url, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.url == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
