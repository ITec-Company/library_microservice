package model

// TestDevDirection ...
func TestAuthor() *Author {
	return &Author{
		ID:       1,
		FullName: "Ivanov Ivan Ivanovich",
	}
}

// TestDevDirection ...
func TestDevDirection() *DevDirection {
	return &DevDirection{
		ID:        1,
		Direction: Backend,
	}
}

// TestDevDirectionDTO ...
func TestDevDirectionDTO() *DevDirectionDTO {
	return &DevDirectionDTO{
		ID:        1,
		Direction: string(Backend),
	}
}

// TestDevSubDirection ...
func TestDevSubDirection() *DevSubDirection {
	return &DevSubDirection{
		ID:           1,
		SubDirection: Golang,
	}
}

// TestDevSubDirectionDTO ...
func TestDevSubDirectionDTO() *DevSubDirectionDTO {
	return &DevSubDirectionDTO{
		ID:           1,
		SubDirection: string(Golang),
	}
}
