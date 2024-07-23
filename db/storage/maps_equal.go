package storage

func MapsEqual(map1, map2 map[string]interface{}) bool {
	if len(map1) != len(map2) {
		return false
	}
	for key, value1 := range map1 {
		if value2, exists := map2[key]; !exists || value1 != value2 {
			return false
		}
	}
	return true
}
