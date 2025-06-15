package strings_pipeline

import (
	"testing"
)

func TestPipeline(t *testing.T) {
	// 示例输入数据
	input := []string{
		"汇灵盏 1599999",
		"五色旗盒 1599999",
		"碧玉葫芦 1599999",
		"碧玉葫芦 1599999",
		"汇灵盏 1599999",
		"金甲仙衣 1799999",
		"金甲仙衣 1799999",
		"嗜血幡 1899999",
		"嗜血幡 1899999",
		"拭剑石 1899999",
		"拭剑石 1899999",
		"飞剑 1899999",
		"九黎战鼓 2199999",
		"九黎战鼓 2199999",
		"风袋 2399999",
		"风袋 2399999",
		"异域风情 2399999",
		"风袋 2399999",
		"鬼谷子 2999999",
		"法宝任务书 2999999",
		"舍利子 880000",
		"月亮石 1104000",
		"太阳石 1312000",
		"太阳石 1312000",
		"光芒石 1440000",
		"红玛瑙 1632000",
		"舍利子 1760000",
		"月亮石 2208000",
		"太阳石 2624000",
		"光芒石 2880000",
		"红玛瑙 3264000",
		"月亮石 4416000",
		"星辉石 4698000",
		"太阳石 5248000",
	}

	// 初始化流水线
	pipeline := &Pipeline{}
	pipeline.AddStep(Deduplicate)
	pipeline.AddStep(NormalizeNames)

	// 解析输入数据
	items := ParseInput(input)

	// 运行流水线处理
	result := pipeline.Run(items)

	// 输出结果
	for _, product := range result {
		t.Logf("{ Name: %q, Prices: %v }\n", product.Name, product.Prices)
	}
}
