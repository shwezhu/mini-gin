package main

import "strings"

type node struct {
	pattern  string // 只有叶子节点才有pattern, 即最终路径
	part     string // 当前节点的 part
	children []*node
	isWild   bool // part 是否含有 : 或 *
}

func (n *node) searchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part {
			return child
		}
	}
	return nil
}

func (n *node) searchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 函数目标: 把 parts 数组中的元素, 从左到右插入到前缀树中, 然后给叶子节点(parts 最后一个元素)添加 pattern 信息
func (n *node) insertChild(pattern string, parts []string, depth int) {
	// 只有叶子节点才有 pattern, 即最终路径
	if len(parts) == depth {
		n.pattern = pattern
		return
	}

	// 插入前缀 可以理解成从根节点到叶子节点的路径
	// 若根节点为a, 大致如下: /a/b/c => a -> b -> c, /a/b/d => a -> b -> d
	// 我们要插入的路径为 /a/b/r, parts = [a, b, r]
	// 递归插入, a, b, 已经存在 所以不需要创建新的节点
	// r 不存在, 需要创建新的节点, c, d, r 都都有相同的前缀 a -> b, 这就是前缀树
	part := parts[depth]
	child := n.searchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	// parts[depth] 已被创建或已经存在, 接下来插入下一个节点 parts[depth+1],
	// 直到 depth == len(parts), 也就最后一个节点, 注意 depth 是从 0 开始的
	// 所以 depth == len(parts) 时, 是最后一个节点, 注意函数的第一个 if 语句, 就是判断这个的
	child.insertChild(pattern, parts, depth+1)
}

// 插入节点的时候 只有叶子节点的 pattern 不为空
// 此函数的本质就是去搜索所有叶子节点的 pattern, 别忘了 parts 就是 pattern 拆出来的数组
func (n *node) matchPattern(parts []string, depth int) *node {
	// ‘*’ 只会出现在 pattern 的最后, 即 /api/chat/v1/*
	// 代表含有 /api/chat/v1/ 前缀的 pattern, 都会匹配到此 pattern
	if len(parts) == depth || strings.HasPrefix(n.part, "*") {
		// pattern 为空, 说明此节点不是叶子节点, 也就是说此节点不是最终路径
		// 比如前缀树为 /a/b/c, /a/b/d,
		// 我们的路径为 /a/b, 此时的 depth == 2, 也就是说我们的路径已经到了 /a/b,
		// 此时 len(parts) == depth 但 /a/b 不是最终路径
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[depth]
	// 这里可能会匹配多个节点, 主要是因为节点树的某个节点的值为 :
	// 比如前缀树为 /a/:id/chat, /a/:name/doc, /a/api/v2,  此时我们的路径为 /a/api/v2
	// 那就会匹配到 /a/:id, /a/:name, /a/api, 三个节点, 然后继续递归匹配到 /a/api/v2 而不是前两个
	children := n.searchChildren(part)
	for _, child := range children {
		result := child.matchPattern(parts, depth+1)
		// 已经找到 pattern, 直接返回
		if result != nil {
			return result
		}
	}

	return nil
}
