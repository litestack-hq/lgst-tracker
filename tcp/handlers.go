package tcp

import (
	"bufio"
	"math"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
	DEVICE_HEART_REGEX = `^\*HQ,(\d{10}),V1,(\d{6}),([VA]),(\d+\.\d+),([NS]),(\d+\.\d+),([EW]),(\d+\.\d{2}),(\d{3}),(\d{6}),([A-Z]{8}),(\d{3}),(\d{2}),(\d+),(\d+),([1-6])#$`
)

type HandlerOpts struct {
	Logger zerolog.Logger
}

type DeviceHeartBeat struct {
	DeviceId         string
	ValidGpsData     bool
	Latitude         float64
	Longitude        float64
	Speed            float32
	AzimuthTrueNorth int
	VehicleStatus    string
	CountryCode      string
	OperatorsCount   int
	BaseStationCount int
	DistrictId       string
	BatteryStatus    int16
	DeviceTime       time.Time
}

var heartBeatRegex *regexp.Regexp

func init() {
	var err error
	heartBeatRegex, err = regexp.Compile(DEVICE_HEART_REGEX)
	if err != nil {
		panic(err)
	}
}

func (h *HandlerOpts) handleHeartBeat(data string, conn net.Conn) {
	regexGroups := heartBeatRegex.FindStringSubmatch(data)

	utc, _ := time.LoadLocation("UTC")
	year, _ := strconv.ParseInt(regexGroups[10][4:6], 10, 64)
	month, _ := strconv.ParseInt(regexGroups[10][2:4], 10, 64)
	days, _ := strconv.ParseInt(regexGroups[10][:2], 10, 64)
	hours, _ := strconv.ParseInt(regexGroups[2][:2], 10, 64)
	minutes, _ := strconv.ParseInt(regexGroups[2][2:4], 10, 64)
	seconds, _ := strconv.ParseInt(regexGroups[2][4:6], 10, 64)

	if year >= 60 {
		year = 1900 + year
	} else {
		year = 2000 + year
	}

	deviceTime := time.Date(int(year), time.Month(month), int(days), int(hours), int(minutes), int(seconds), 0, utc)

	heartBeat := DeviceHeartBeat{
		DeviceId: regexGroups[1],
		ValidGpsData: func() bool {
			return strings.ToUpper(regexGroups[3]) == "A"
		}(),
		Latitude: func() float64 {
			latSlice := strings.Split(regexGroups[4], ".")
			degs, _ := strconv.ParseInt(latSlice[0][:len(latSlice[0])-2], 10, 64)
			mins, _ := strconv.ParseFloat(latSlice[0][len(latSlice[0])-2:]+"."+latSlice[1], 64)
			mins = mins / 60
			return math.Floor(((float64(degs) + mins) * 100000)) / 100000
		}(),
		Longitude: func() float64 {
			lgtSlice := strings.Split(regexGroups[6], ".")
			degs, _ := strconv.ParseInt(lgtSlice[0][:len(lgtSlice[0])-2], 10, 64)
			mins, _ := strconv.ParseFloat(lgtSlice[0][len(lgtSlice[0])-2:]+"."+lgtSlice[1], 64)
			mins = mins / 60
			return math.Floor(((float64(degs) + mins) * 100000)) / 100000
		}(),
		Speed: func() float32 {
			speed, err := strconv.ParseFloat(regexGroups[8], 32)
			if err != nil {
				h.Logger.Err(err).Msg("failed to parse speed value")
			}
			return float32(speed)
		}(),
		AzimuthTrueNorth: func() int {
			trueNorth, err := strconv.ParseInt(regexGroups[9], 10, 64)
			if err != nil {
				h.Logger.Err(err).Msg("failed to parse azimuth true north value")
			}
			return int(trueNorth)
		}(),
		VehicleStatus: func() string {
			return regexGroups[11]
		}(),
		CountryCode: func() string {
			return regexGroups[12]
		}(),
		OperatorsCount: func() int {
			operatorCount, err := strconv.ParseInt(regexGroups[13], 10, 64)
			if err != nil {
				h.Logger.Err(err).Msg("failed to parse operator count value")
			}
			return int(operatorCount)
		}(),
		BaseStationCount: func() int {
			baseStationCount, err := strconv.ParseInt(regexGroups[14], 10, 64)
			if err != nil {
				h.Logger.Err(err).Msg("failed to parse base station count value")
			}
			return int(baseStationCount)
		}(),
		DistrictId: func() string {
			return regexGroups[15]
		}(),
		BatteryStatus: func() int16 {
			bat, err := strconv.ParseInt(regexGroups[16], 10, 64)
			if err != nil {
				h.Logger.Err(err).Msg("failed to parse battery status value")
			}
			return int16(bat)
		}(),
		DeviceTime: deviceTime,
	}

	serverTime := time.Now().UTC().Format("20060102150405")
	response := "*HQ," + heartBeat.DeviceId + ",V4V1," + serverTime + "#"

	_, err := conn.Write([]byte(response))
	if err != nil {
		h.Logger.Err(err).Msg("TCP write failed")
	}

	h.Logger.Info().Interface("heart_beat", heartBeat).Interface("response", response).Msg("handle device heartbeat")
}

func New(logger zerolog.Logger) *HandlerOpts {
	return &HandlerOpts{
		Logger: logger,
	}
}

func (h *HandlerOpts) DevicePingHandler(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := bufio.NewReader(conn).ReadString('#')
		if heartBeatRegex.MatchString(data) {
			h.handleHeartBeat(data, conn)
			return
		}

		if err != nil {
			if err.Error() != "EOF" {
				h.Logger.Err(err).Msg("failed to read data")
			}
			return
		}

		h.Logger.Info().Str("data", data).Msg("device data")
	}
}
