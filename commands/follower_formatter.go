package commands

import "github.com/dghubble/go-twitter/twitter"

func MapUserToListValue(user twitter.User) []ListValue {
	vsm := make([]ListValue, 1)
	vsm[0] = ListValue{Value:user.ScreenName}
	return vsm
}

func MapUsersToListValueList(users []twitter.User, f func(user twitter.User) []ListValue) [][]ListValue {
	vsm := make([][]ListValue, len(users))
	for i, v := range users {
		vsm[i] = f(v)
	}
	return vsm
}

func PrintUsersAsList(follower []twitter.User) {
	headerValues := []ListValue{
		{Value:"Name"},
	}

	realValues := MapUsersToListValueList(follower, MapUserToListValue)

	formatList(headerValues, realValues)
}