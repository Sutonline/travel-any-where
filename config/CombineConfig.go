package config

type Config struct {
	HotelCnt            int
	RestaurantCnt       int
	ScenicCnt           int
	RestaurantCntPerDay int
	ScenicCntPerDay     int
}

func GetDefaultConfig(days int) Config {
	return Config{
		1,
		3 * days,
		3 * days,
		3,
		3,
	}
}
