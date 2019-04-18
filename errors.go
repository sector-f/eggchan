package eggchan

type UnauthorizedError struct{}

func (e UnauthorizedError) Error() string {
	return "Unauthorized"
}

type CategoryNotFoundError struct{}

func (e CategoryNotFoundError) Error() string {
	return "Category not found"
}

type BoardNotFoundError struct{}

func (e BoardNotFoundError) Error() string {
	return "Board not found"
}

type ThreadNotFoundError struct{}

func (e ThreadNotFoundError) Error() string {
	return "Thread not found"
}

type UserNotFoundError struct{}

func (e UserNotFoundError) Error() string {
	return "User not found"
}

type DatabaseError struct{}

func (e DatabaseError) Error() string {
	return "Database error"
}

type PermissionDeniedError struct{}

func (e PermissionDeniedError) Error() string {
	return "Permission denied"
}

type UnimplementedError struct{}

func (e UnimplementedError) Error() string {
	return "Unimplemented"
}
