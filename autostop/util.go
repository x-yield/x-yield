package autostop

import (
	"errors"
	"strconv"
	"unicode"
)

// "1m", "4d5h", "1s", "3h141m5s", "55"
func ParseDuration(duration string) (seconds int32, err error) {
	var stack = ""
	if len(duration) == 0 {
		err = errors.New("empty duration string")
	}
	for _, s := range duration {
		if unicode.IsDigit(s) {
			stack += string(s)
		} else {
			if len(stack) == 0 {
				return 0, errors.New("invalid autostop duration format")
			}
			switch s {
			case 'd':
				days, _ := strconv.Atoi(stack)
				seconds += int32(days) * 24 * 3600
				stack = ""
			case 'h':
				hours, _ := strconv.Atoi(stack)
				seconds += int32(hours) * 3600
				stack = ""
			case 'm':
				minutes, _ := strconv.Atoi(stack)
				seconds += int32(minutes) * 60
				stack = ""
			case 's':
				secs, _ := strconv.Atoi(stack)
				seconds += int32(secs)
				stack = ""
			default:
				return 0, errors.New("invalid autostop duration format")
			}
		}
	}
	if len(stack) > 0 {
		secs, _ := strconv.Atoi(stack)
		seconds += int32(secs)
		stack = ""
	}
	return
}

// "1s500ms"
func ParseMilliseconds(duration string) (milliseconds int32, err error) {
	var stack = ""
	if len(duration) == 0 {
		err = errors.New("empty milliseconds string")
	}
	for i := range duration {
		s := rune(duration[i])
		if unicode.IsDigit(s) {
			stack += string(s)
		} else {
			if len(stack) == 0 {
				return 0, errors.New("invalid autostop milliseconds format")
			}
			switch s {
			case 'd':
				days, _ := strconv.Atoi(stack)
				milliseconds += int32(days) * 24 * 3600 * 1000
				stack = ""
			case 'h':
				hours, _ := strconv.Atoi(stack)
				milliseconds += int32(hours) * 3600 * 1000
				stack = ""
			case 'm':
				if i < len(duration) - 1 && rune(duration[i+1]) == 's' {
					continue
				}
				minutes, _ := strconv.Atoi(stack)
				milliseconds += int32(minutes) * 60 * 1000
				stack = ""
			case 's':
				secs, _ := strconv.Atoi(stack)
				if rune(duration[i-1]) != 'm' {
					secs *= 1000
				}
				milliseconds += int32(secs)
				stack = ""
			default:
				return 0, errors.New("invalid autostop milliseconds format")
			}
		}
	}
	if len(stack) > 0 {
		ms, _ := strconv.Atoi(stack)
		milliseconds += int32(ms)
		stack = ""
	}
	return
}
