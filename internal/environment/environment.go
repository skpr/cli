package environment

import "github.com/skpr/api/pb"

// Contains checks if an environment already exists.
func Contains(name string, list []*pb.Environment) bool {
	for _, item := range list {
		if item.Name == name {
			return true
		}
	}

	return false
}
