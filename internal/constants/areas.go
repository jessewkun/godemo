package constants

import (
	_ "embed"
	"encoding/json"
)

//go:embed areas.json
var areasJSON []byte

type Area struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Children []Area `json:"children"`
}

type AreaItem struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	CityCode     string `json:"cityCode"`
	ProvinceCode string `json:"provinceCode"`
}

// 解析JSON数据并构建三级结构
func parseAreasData() map[string]Area {
	var items []AreaItem
	if err := json.Unmarshal(areasJSON, &items); err != nil {
		panic("解析areas.json失败: " + err.Error())
	}

	// 构建三级结构
	areaMap := make(map[string]Area)

	// 按省份分组
	provinceMap := make(map[string]map[string][]AreaItem)
	for _, item := range items {
		if provinceMap[item.ProvinceCode] == nil {
			provinceMap[item.ProvinceCode] = make(map[string][]AreaItem)
		}
		provinceMap[item.ProvinceCode][item.CityCode] = append(provinceMap[item.ProvinceCode][item.CityCode], item)
	}

	// 构建省份结构
	for provinceCode, cities := range provinceMap {
		province := Area{
			Code:     provinceCode,
			Name:     getProvinceName(provinceCode),
			Children: make([]Area, 0),
		}

		// 构建城市结构
		for cityCode, districts := range cities {
			city := Area{
				Code:     cityCode,
				Name:     getCityName(cityCode),
				Children: make([]Area, 0),
			}

			// 构建区县结构
			for _, district := range districts {
				districtArea := Area{
					Code:     district.Code,
					Name:     district.Name,
					Children: []Area{},
				}
				city.Children = append(city.Children, districtArea)
			}

			province.Children = append(province.Children, city)
		}

		areaMap[provinceCode] = province
	}

	return areaMap
}

// ProvinceMap 省份名称
var ProvinceMap = map[string]string{
	"11": "北京市",
	"12": "天津市",
	"13": "河北省",
	"14": "山西省",
	"15": "内蒙古自治区",
	"21": "辽宁省",
	"22": "吉林省",
	"23": "黑龙江省",
	"31": "上海市",
	"32": "江苏省",
	"33": "浙江省",
	"34": "安徽省",
	"35": "福建省",
	"36": "江西省",
	"37": "山东省",
	"41": "河南省",
	"42": "湖北省",
	"43": "湖南省",
	"44": "广东省",
	"45": "广西壮族自治区",
	"46": "海南省",
	"50": "重庆市",
	"51": "四川省",
	"52": "贵州省",
	"53": "云南省",
	"54": "西藏自治区",
	"61": "陕西省",
	"62": "甘肃省",
	"63": "青海省",
	"64": "宁夏回族自治区",
	"65": "新疆维吾尔自治区",
	"71": "台湾省",
	"81": "香港特别行政区",
	"82": "澳门特别行政区",
}

// 获取省份名称
func getProvinceName(code string) string {
	return ProvinceMap[code]
}

// CityMap 城市名称
var CityMap = map[string]string{
	"1101": "市辖区",
	"1201": "市辖区",
	"1301": "石家庄市",
	"1302": "唐山市",
	"1303": "秦皇岛市",
	"1304": "邯郸市",
	"1305": "邢台市",
	"1306": "保定市",
	"1307": "张家口市",
	"1308": "承德市",
	"1309": "沧州市",
	"1310": "廊坊市",
	"1311": "衡水市",
	"1401": "太原市",
	"1402": "大同市",
	"1403": "阳泉市",
	"1404": "长治市",
	"1405": "晋城市",
	"1406": "朔州市",
	"1407": "晋中市",
	"1408": "运城市",
	"1409": "忻州市",
	"1410": "临汾市",
	"1411": "吕梁市",
	"1501": "呼和浩特市",
	"1502": "包头市",
	"1503": "乌海市",
	"1504": "赤峰市",
	"1505": "通辽市",
	"1506": "鄂尔多斯市",
	"1507": "呼伦贝尔市",
	"1508": "巴彦淖尔市",
	"1509": "乌兰察布市",
	"1522": "兴安盟",
	"1525": "锡林郭勒盟",
	"1529": "阿拉善盟",
	"2101": "沈阳市",
	"2102": "大连市",
	"2103": "鞍山市",
	"2104": "抚顺市",
	"2105": "本溪市",
	"2106": "丹东市",
	"2107": "锦州市",
	"2108": "营口市",
	"2109": "阜新市",
	"2110": "辽阳市",
	"2111": "盘锦市",
	"2112": "铁岭市",
	"2113": "朝阳市",
	"2114": "葫芦岛市",
	"2201": "长春市",
	"2202": "吉林市",
	"2203": "四平市",
	"2204": "辽源市",
	"2205": "通化市",
	"2206": "白山市",
	"2207": "松原市",
	"2208": "白城市",
	"2224": "延边朝鲜族自治州",
	"2301": "哈尔滨市",
	"2302": "齐齐哈尔市",
	"2303": "鸡西市",
	"2304": "鹤岗市",
	"2305": "双鸭山市",
	"2306": "大庆市",
	"2307": "伊春市",
	"2308": "佳木斯市",
	"2309": "七台河市",
	"2310": "牡丹江市",
	"2311": "黑河市",
	"2312": "绥化市",
	"2327": "大兴安岭地区",
	"3101": "市辖区",
	"3201": "南京市",
	"3202": "无锡市",
	"3203": "徐州市",
	"3204": "常州市",
	"3205": "苏州市",
	"3206": "南通市",
	"3207": "连云港市",
	"3208": "淮安市",
	"3209": "盐城市",
	"3210": "扬州市",
	"3211": "镇江市",
	"3212": "泰州市",
	"3213": "宿迁市",
	"3301": "杭州市",
	"3302": "宁波市",
	"3303": "温州市",
	"3304": "嘉兴市",
	"3305": "湖州市",
	"3306": "绍兴市",
	"3307": "金华市",
	"3308": "衢州市",
	"3309": "舟山市",
	"3310": "台州市",
	"3311": "丽水市",
	"3401": "合肥市",
	"3402": "芜湖市",
	"3403": "蚌埠市",
	"3404": "淮南市",
	"3405": "马鞍山市",
	"3406": "淮北市",
	"3407": "铜陵市",
	"3408": "安庆市",
	"3410": "黄山市",
	"3411": "滁州市",
	"3412": "阜阳市",
	"3413": "宿州市",
	"3415": "六安市",
	"3416": "亳州市",
	"3417": "池州市",
	"3418": "宣城市",
	"3501": "福州市",
	"3502": "厦门市",
	"3503": "莆田市",
	"3504": "三明市",
	"3505": "泉州市",
	"3506": "漳州市",
	"3507": "南平市",
	"3508": "龙岩市",
	"3509": "宁德市",
	"3601": "南昌市",
	"3602": "景德镇市",
	"3603": "萍乡市",
	"3604": "九江市",
	"3605": "新余市",
	"3606": "鹰潭市",
	"3607": "赣州市",
	"3608": "吉安市",
	"3609": "宜春市",
	"3610": "抚州市",
	"3611": "上饶市",
	"3701": "济南市",
	"3702": "青岛市",
	"3703": "淄博市",
	"3704": "枣庄市",
	"3705": "东营市",
	"3706": "烟台市",
	"3707": "潍坊市",
	"3708": "济宁市",
	"3709": "泰安市",
	"3710": "威海市",
	"3711": "日照市",
	"3713": "临沂市",
	"3714": "德州市",
	"3715": "聊城市",
	"3716": "滨州市",
	"3717": "菏泽市",
	"4101": "郑州市",
	"4102": "开封市",
	"4103": "洛阳市",
}

// 获取城市名称
func getCityName(code string) string {
	return CityMap[code]
}

// AreaMap 使用embed嵌入的JSON数据构建的三级省市区结构
var AreaMap = parseAreasData()
