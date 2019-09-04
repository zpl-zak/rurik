package core

// NewNotNull not null
func (o *Object) NewNotNull() {
	o.Size = []int32{8, 8}

	o.GetAABB = GetSpriteAABB
}
