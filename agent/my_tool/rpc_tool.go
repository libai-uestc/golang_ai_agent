package mytool

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

var (
	GAODE_KEY = os.Getenv("GAODE_KEY")
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http GET %w", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http response StatusCode %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}
	return data, nil
}

func GetOutboundIP() string {
	data, err := HttpGet("https://httpbin.org/ip")
	if err != nil {
		log.Printf("获取对外ip失败: %s", err)
		return ""
	}
	mp := make(map[string]string, 1)
	sonic.Unmarshal(data, &mp)
	return mp["origin"]
}

type Location struct {
	Status    string
	Province  string
	City      string
	Rectangle string
	Longitude float64
	Latitude  float64
	AdCode    string
}

func GetMyLocation() (*Location, error) {
	ip := GetOutboundIP()
	url := "https://restapi.amap.com/v3/ip?key=" + GAODE_KEY + "&ip=" + ip
	data, err := HttpGet(url)
	if err != nil {
		return nil, err
	}
	var location Location
	err = sonic.Unmarshal(data, &location)
	if err != nil {
		return nil, err
	}
	if location.Status != "1" {
		return nil, fmt.Errorf("GetMyLocation failed, status is: %s", location.Status)
	}
	arr := strings.Split(location.Rectangle, ";")
	if len(arr) == 2 {
		var Longitude1, Longitude2, Latitude1, Latitude2 float64
		brr := strings.Split(arr[0], ",")
		if len(brr) == 2 {
			if v, err := strconv.ParseFloat(brr[0], 64); err == nil {
				Longitude1 = v
			}
			if v, err := strconv.ParseFloat(brr[1], 64); err == nil {
				Latitude1 = v
			}
		}
		brr = strings.Split(arr[1], ",")
		if len(brr) == 2 {
			if v, err := strconv.ParseFloat(brr[0], 64); err == nil {
				Longitude2 = v
			}
			if v, err := strconv.ParseFloat(brr[1], 64); err == nil {
				Latitude2 = v
			}
		}
		location.Longitude = (Longitude1 + Longitude2) / 2
		location.Latitude = (Latitude1 + Latitude2) / 2
	}
	return &location, nil
}

func CreateLocationTool() tool.InvokableTool {
	locationTool, err := utils.InferTool("location_tool", "获取当前的地理位置，包含省、城市（含城市名称和城市编码）、经纬度等信息", func(ctx context.Context, input struct{}) (output string, err error) {
		location, err := GetMyLocation()
		if err != nil {
			return "", fmt.Errorf("获取当前地理位置失败：%w", err)
		}
		return fmt.Sprintf("当前城市是%s，经度 %.8f，纬度 %.8f，城市编码为 %s，属于 %s", location.City, location.Longitude, location.Latitude, location.AdCode, location.Province), nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return locationTool
}

type WeatherCast struct {
	Date             string
	Week             string
	DayWeather       string
	DayTemperature   string `json:"daytemp"`
	DayWindDirection string `json:"daywind"`  // 风向
	DayWindPowerr    string `json:"daypower"` // 风力
	// 晚上
	NightWeather       string
	NightTemperature   string `json:"nighttemp"`
	NightWindDirection string `json:"nightwind"`
	NightWindPowerr    string `json:"nightpower"`
}

type LiveWeather struct {
	Weather       string
	Temperature   string
	WindDirection string
	WindPower     string
	Humidity      string
}

type WeatherInfo struct {
	Status    string        `json:"status"`
	Lives     []LiveWeather `json:"lives"` // 实况天气
	Forecasts []struct {
		Casts []WeatherCast //预报数据 list 结构，元素 cast,按顺序为当天、第二天、第三天、第四天的预报数据
	} // 预报天气
}

func WeatherForcast(GaodeCityCode string, live bool) *WeatherInfo {
	url := "https://restapi.amap.com/v3/weather/weatherInfo?key=" + GAODE_KEY + "&city=" + GaodeCityCode
	if !live {
		url += "&extensions=all"
	}
	data, err := HttpGet(url)
	if err != nil {
		log.Printf("get future weather failed: %s", err)
		return nil
	}
	var weather WeatherInfo
	err = sonic.Unmarshal(data, &weather)
	if err != nil {
		log.Printf("get future weather failed,json Unmarshal response failed: %s", err)
		return nil
	}
	if weather.Status != "1" {
		log.Printf("get future weather failed, status is: %s", weather.Status)
		return nil
	}
	return &weather
}

func GetFutureWeather(GaodeCityCode string, day int) *WeatherCast {
	weather := WeatherForcast(GaodeCityCode, false)
	if weather == nil {
		return nil
	}
	if len(weather.Forecasts) == 0 {
		return nil
	}
	if day < len(weather.Forecasts[0].Casts) {
		return &weather.Forecasts[0].Casts[day]
	} else {
		return nil
	}
}

func GetCurrentWeather(GaodeCityCode string) *LiveWeather {
	weather := WeatherForcast(GaodeCityCode, true)
	if weather == nil {
		return nil
	}
	if len(weather.Lives) == 0 {
		return nil
	}
	return &weather.Lives[0]
}

type WeatherParams struct {
	CityCode string `json:"city_code" jsonschema:"required,description=城市的编码"`
	Day      string `json:"day" jsonschema:"required,enum=0,enum=1,enum=2,enum=3,enum=4" jsonschema_description:"获取现在和未来的气象数据。0表示当前实时的气象数据，1表示今天的气象数据，2表示明天的气象数据，3表示后天的气象数据，4表示大后天的气象数据，最大只能到4"` // description中如果包含,还需要转义（,在tag中是分隔的意思），此时建议用jsonschema_description
}

func CreateWeatherTool() tool.InvokableTool {
	weatherTool, err := utils.InferTool("weather_tool", "获取气象数据，包含天气、温度、风力、风向等，既可以获得当前的实时气象数据，也可以获取未来3天的气象数据", func(ctx context.Context, params WeatherParams) (output string, err error) {
		day, err := strconv.Atoi(params.Day)
		if err != nil {
			return "", errors.New("参数day非法")
		}
		if day == 0 {
			weather := GetCurrentWeather(params.CityCode)
			if weather != nil {
				return fmt.Sprintf("现在的天气情况是: %s，温度 %s℃，湿度 %s%%，风力 %s级，风向 %s", weather.Weather, weather.Temperature, weather.Humidity, weather.WindPower, weather.WindDirection), nil
			} else {
				return "", errors.New("获取不到天气数据")
			}
		} else if day >= 1 && day <= 4 {
			weather := GetFutureWeather(params.CityCode, day-1)
			if weather != nil {
				return fmt.Sprintf("%s 白天的天气情况是：%s，温度 %s℃，风力 %s级，风向 %s；晚上的天气情况是：%s，温度 %s℃，风力 %s级，风向 %s", weather.Date, weather.DayWeather, weather.DayTemperature, weather.DayWindPowerr, weather.DayWindDirection, weather.NightWeather, weather.NightTemperature, weather.NightWindPowerr, weather.NightWindDirection), nil
			} else {
				return "", errors.New("获取不到天气数据")
			}
		} else {
			return "", errors.New("参数day非法")
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	return weatherTool
}

func UseLocationTool() {
	ctx := context.Background()
	tool := CreateLocationTool()

	resp, err := tool.InvokableRun(ctx, "{}")
	if err != nil {
		log.Fatalf("get location failed: %s", err)
	}

	fmt.Println(resp)
}

func UseWeatherTool(day int) {
	ctx := context.Background()
	tool := CreateWeatherTool()

	weatherRequest := WeatherParams{
		CityCode: "510100",
		Day:      strconv.Itoa(day),
	}
	jsonReq, err := sonic.Marshal(weatherRequest)
	if err != nil {
		log.Fatalf("Marshal of weather request failed: %s", err)
	}
	resp, err := tool.InvokableRun(ctx, string(jsonReq))
	if err != nil {
		log.Fatalf("get weather failed: %s", err)
	}
	fmt.Println(resp)
}
