package strings_pipeline

import (
	"sort"
	"strconv"
	"strings"
)

// Item 表示原始商品数据结构
type Item struct {
	RawName string // 原始商品名
	Name    string // 处理后的商品名
	Price   int    // 商品价格
}

// Product 表示合并后的商品结构
type Product struct {
	Name   string
	Prices []int
}

// Pipeline 流水线结构
type Pipeline struct {
	steps []func([]Item) []Item
}

// AddStep 添加处理步骤
func (p *Pipeline) AddStep(step func([]Item) []Item) {
	p.steps = append(p.steps, step)
}

// Run 执行流水线处理
func (p *Pipeline) Run(items []Item) []Product {
	// 执行所有处理步骤
	for _, step := range p.steps {
		items = step(items)
	}

	// 最终合并处理
	return mergeItems(items)
}

// 去重处理函数
func Deduplicate(items []Item) []Item {
	seen := make(map[string]struct{})
	result := make([]Item, 0)

	for _, item := range items {
		key := item.RawName + "|" + strconv.Itoa(item.Price)
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// 商品名标准化处理函数
func NormalizeNames(items []Item) []Item {
	for i := range items {
		// 示例处理：去除第一个修饰词
		parts := strings.Split(items[i].RawName, " ")
		if len(parts) > 1 {
			items[i].Name = strings.Join(parts[1:], " ")
		} else {
			items[i].Name = items[i].RawName
		}
	}
	return items
}

// 合并商品并排序价格
func mergeItems(items []Item) []Product {
	productMap := make(map[string][]int)

	for _, item := range items {
		productMap[item.Name] = append(productMap[item.Name], item.Price)
	}

	products := make([]Product, 0, len(productMap))
	for name, prices := range productMap {
		sort.Ints(prices)
		products = append(products, Product{Name: name, Prices: prices})
	}

	return products
}

// 字符串解析为Item的预处理函数
func ParseInput(input []string) []Item {
	items := make([]Item, 0, len(input))

	for _, s := range input {
		parts := strings.SplitN(s, " ", 2)
		if len(parts) != 2 {
			continue
		}

		price, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		items = append(items, Item{
			RawName: parts[0],
			Price:   price,
		})
	}
	return items
}
