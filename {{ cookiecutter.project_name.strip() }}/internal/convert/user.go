package convert

import (
	pb "{{ cookiecutter.go_module_path.strip('/') }}/api/go/{{ cookiecutter.project_name.strip() }}/v1"
	"{{ cookiecutter.go_module_path.strip('/') }}/internal/repos/users"
)

// UserDBToProto converts a database User model to a protobuf User message.
func UserDBToProto(user *users.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Age:       user.Age,
		CreatedAt: user.CreatedAt.Unix(),
		UpdatedAt: user.UpdatedAt.Unix(),
	}
}
