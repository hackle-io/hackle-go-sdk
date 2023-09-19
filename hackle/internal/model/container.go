package model

type Container struct {
	ID       int64
	BucketID int64
	Groups   []ContainerGroup
}

type ContainerGroup struct {
	ID          int64
	Experiments []int64
}

func (c *Container) GetGroup(containerGroupID int64) (ContainerGroup, bool) {
	for _, group := range c.Groups {
		if group.ID == containerGroupID {
			return group, true
		}
	}
	return ContainerGroup{}, false
}
