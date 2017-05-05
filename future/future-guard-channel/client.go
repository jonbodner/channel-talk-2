package main

import (
	"time"
	"errors"
	"net/http"
)

type User struct {
	Id string
	Name string
	State string
}

type Branch struct {
	Id string
	Street string
	City string
	State string
	Zip string
}

//not a real implementation
func GetUserFromMicroservice(id string) (User, error) {
	return User{Id: id, Name: "Bob Bobson", State: "ME"}, nil
}

//not a real implementation
func GetBranchesInState(state string) ([]Branch, error) {
	return []Branch{}, nil
}

//not a real implementation
func doSomethingElseForABit(id string) {

}

func FindBranchListForUser(id string) ([]Branch, error) {
	//create my future
	future := New(func() (interface{}, error) {
		return GetUserFromMicroservice(id)
	}).Then(func(val interface{}) (interface{}, error) {
		return GetBranchesInState(val.(User).State)
	})

	doSomethingElseForABit(id)
	//wait for it to return, but not too long
	branches, err := future.GetUntil(2*time.Second)
	if err == FUTURE_TIMEOUT {
		//never mind, don't care about results
		future.Cancel()
		//alert that the system is running slowly
		return nil, errors.New("Timed out getting branch list for user")
	}
	if err != nil {
		return nil, err
	}
	return branches.([]Branch), nil
}
