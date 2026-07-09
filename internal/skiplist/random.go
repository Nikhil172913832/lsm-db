package skiplist

import "math/rand/v2"

func RandomLevel(maxLevel int) int {
	level := 1
	temp := rand.Float64()
	for {
		if level == maxLevel || temp < 0.5 {
			return level
		}
		level++
		temp = rand.Float64()
	}
}
