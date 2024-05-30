package model

import pb "project/proto/user"

type User struct {
	Name         string `bson:"name,omitempty"`
	Age          int32  `bson:"age,omitempty"`
	MobileNumber int32  `bson:"phoneNumber,omitempty"`
	EmailId      string `bson:"email,omitempty"`
}

func (u *User) ConvertToSchema(user *pb.User) {
	u.Name = user.GetName()
	u.Age = user.GetAge()
	u.MobileNumber = user.GetPhoneNumber()
	u.EmailId = user.GetEmail()
}

func (u *User) ConvertToProto() *pb.User {
	return &pb.User{
		Name:        u.Name,
		Age:         u.Age,
		PhoneNumber: u.MobileNumber,
		Email:       u.EmailId,
	}
}
