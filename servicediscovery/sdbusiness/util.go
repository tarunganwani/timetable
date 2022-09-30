package sdbusiness

import "servicediscovery/sdmodel"

// Remove from slice without caring about order
func Remove(s []sdmodel.ServiceData, i int) []sdmodel.ServiceData {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
