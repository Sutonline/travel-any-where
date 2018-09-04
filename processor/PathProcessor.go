package processor

import (
	"../base"
	"../config"
)

type DayPath struct {
	Hotel       base.DbPojo
	Restaurants []base.DbPojo
	Scenes      []base.DbPojo
	Path        []base.DbPojo
}

type TripPath struct {
	Paths        []DayPath
	SingleBudget float64
}

func GenTripPath(city string, days int, budget float64) []TripPath {
	hotels := base.SelectHotelByCity(city)
	restaurants := base.SelectRestaurantByCity(city)
	scenes := base.SelectScenic(city)
	dayTrips := CombineDays(hotels, restaurants, scenes, days)

	paths := GetPath(dayTrips)

	tripPaths := make([]TripPath, 10)
	index := 0
	for i := 0; i < len(paths) && index < 10; i++ {
		if paths[i].SingleBudget <= budget {
			tripPaths[index] = paths[i]
			index++
		}
	}

	return tripPaths
}

func GetPath(trips []DayTrip) []TripPath {
	paths := make([]TripPath, len(trips))
	for index, trip := range trips {
		paths[index] = genPath(trip)
	}

	return paths
}

func genPath(trip DayTrip) TripPath {
	dayPaths := make([]DayPath, trip.days)
	defaultConfig := config.GetDefaultConfig(trip.days)

	restaurantArray := sliceArray(trip.restaurants, defaultConfig.RestaurantCntPerDay, trip.days)
	scenicArray := sliceArray(trip.scenes, defaultConfig.ScenicCntPerDay, trip.days)

	for i := 0; i < trip.days; i++ {
		path := calculatePath(trip.hotel, restaurantArray[i], scenicArray[i])
		dayPaths[i] = DayPath{
			trip.hotel,
			restaurantArray[i],
			scenicArray[i],
			path,
		}
	}

	return TripPath{dayPaths, trip.singleBudget}
}

func sliceArray(pojos []base.DbPojo, cntPerDay int, days int) [][]base.DbPojo {
	resultArray := make([][]base.DbPojo, days)

	for i := 0; i < days; i++ {
		resultArray[i] = pojos[i*cntPerDay : i*cntPerDay+cntPerDay]
	}

	return resultArray
}

// 通过经纬度的计算得到一个路线
func calculatePath(hotel base.DbPojo, restaurants []base.DbPojo, scenes []base.DbPojo) []base.DbPojo {
	// 加2是因为最后还要回到酒店
	path := make([]base.DbPojo, 2+len(restaurants)+len(scenes))

	// 第一步酒店出发
	index := 0
	path[index] = hotel

	selectorMap := make(map[uint]base.DbPojo, len(restaurants)+len(scenes))
	// 寻找最近的餐厅
	index++
	breakfast := getNearestDbPojo(hotel, restaurants, selectorMap)
	path[index] = breakfast
	// 寻找最近的一个景点节点
	index++
	firstScenic := getNearestDbPojo(breakfast, scenes, selectorMap)
	path[index] = firstScenic
	// 寻找最近的餐厅
	index++
	lunch := getNearestDbPojo(firstScenic, restaurants, selectorMap)
	path[index] = lunch

	// 寻找最近的一个景点
	index++
	secondScenic := getNearestDbPojo(lunch, scenes, selectorMap)
	path[index] = secondScenic
	// 寻找最近的一个景点
	index++
	thirdScenic := getNearestDbPojo(secondScenic, scenes, selectorMap)
	path[index] = thirdScenic
	// 寻找晚餐
	index++
	supper := getNearestDbPojo(thirdScenic, restaurants, selectorMap)
	path[index] = supper
	// 返回酒店
	index++
	path[index] = hotel

	return path
}

func getNearestDbPojo(from base.DbPojo, targets []base.DbPojo, selectorMap map[uint]base.DbPojo) base.DbPojo {
	var tempDbPojo base.DbPojo
	var tempDistance float64
	tempDistance = 0
	for _, value := range targets {
		_, ok := selectorMap[value.Id]
		if !ok {
			distance := base.EarthDistance(from.Latitude, from.Longitude, value.Latitude, value.Longitude)
			if tempDistance == 0 || tempDistance > distance {
				tempDistance = distance
				tempDbPojo = value
			}
		}

	}

	selectorMap[tempDbPojo.Id] = tempDbPojo

	return tempDbPojo
}
