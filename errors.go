package eggchan

type UnauthorizedError struct{}

func (e UnauthorizedError) Error() string {
	return "Unauthorized"
}

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "Not found"
}

type DatabaseError struct{}

func (e DatabaseError) Error() string {
	return "Database error"
}

type UnimplementedError struct{}

func (e UnimplementedError) Error() string {
	return "Unimplemented"
}
