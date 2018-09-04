package processor

import (
	"../base"
	"../config"
)

type DayTrip struct {
	hotel        base.DbPojo
	scenes       []base.DbPojo
	restaurants  []base.DbPojo
	days         int
	singleBudget float64
}

type DbPojoItem struct {
	items  []base.DbPojo
	scores float64
}

// 返回笛卡尔积的列表 [scenes, restaurants, hotels, scenes, restaurants, hotels]
func CombineDays(hotels []base.DbPojo, restaurants []base.DbPojo, scenes []base.DbPojo, days int) []DayTrip {
	defaultConfig := config.GetDefaultConfig(days)

	// 取得要组合的各个类型
	combineHotels := base.Combine(hotels, defaultConfig.HotelCnt)
	combineRestaurants := base.Combine(restaurants, defaultConfig.RestaurantCnt)
	combineScenics := base.Combine(scenes, defaultConfig.ScenicCnt)

	// 排序
	sortScenes := sliceByScore(combineScenics, 10)
	sortRestaurants := sliceByScore(combineRestaurants, 10)
	sortHotels := sliceByScore(combineHotels, 10)

	// 组合成结果
	return cartesianType(sortHotels, sortRestaurants, sortScenes, days)
}

func sliceByScore(srcArray [][]base.DbPojo, size int) [][]base.DbPojo {
	if len(srcArray) <= size {
		return srcArray
	}

	pojos := make([][]base.DbPojo, size)
	items := make([]DbPojoItem, len(srcArray))
	var score float64
	for index, srcItem := range srcArray {
		score = 0
		for _, item := range srcItem {
			score = score + item.Score
		}
		items[index] = DbPojoItem{
			srcItem,
			score,
		}
	}

	// 排序
	SortItem(items, func(p, q *DbPojoItem) bool {
		return p.scores > q.scores
	})

	for i := 0; i < size; i++ {
		pojos[i] = items[i].items
	}

	return pojos
}

// 对得到的结果集组成以天为单位的对象
func cartesianType(hotels [][]base.DbPojo, restaurants [][]base.DbPojo, scenes [][]base.DbPojo, days int) []DayTrip {
	var dayTrips = make([]DayTrip, len(hotels)*len(restaurants)*len(scenes))
	index := 0
	var singleBudget float64
	for _, hotel := range hotels {
		for _, restaurant := range restaurants {
			for _, scenic := range scenes {
				singleBudget = calculateBudget(hotel[0], restaurant, scenic, days)
				dayTrips[index] = DayTrip{
					hotel[0],
					scenic,
					restaurant,
					days,
					singleBudget,
				}
				index++
			}
		}
	}

	return dayTrips
}

func calculateBudget(hotel base.DbPojo, restaurants []base.DbPojo, scenes []base.DbPojo, days int) float64 {
	var singleBudget float64
	singleBudget = 0
	// 加酒店的钱
	singleBudget = hotel.MiddlePrice * float64(days)
	// 加餐馆的钱
	for _, restaurant := range restaurants {
		singleBudget += restaurant.MiddlePrice
	}
	// 加景点的钱
	for _, scenic := range scenes {
		singleBudget += scenic.MiddlePrice
	}

	return singleBudget
}
