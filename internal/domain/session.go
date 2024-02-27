package domain

import "time"

type Session struct {
	Id        string
	UserId    int
	ExpiresAt time.Time
}

func SortSessionsByTime(sessions *[]Session) {
	for i := 0; i+1 < len(*sessions); i++ {
		for j := 0; j+1 < len(*sessions); j++ {
			if (*sessions)[j+1].ExpiresAt.Unix() < (*sessions)[j].ExpiresAt.Unix() {
				t := (*sessions)[j+1]
				(*sessions)[j+1] = (*sessions)[j]
				(*sessions)[j] = t
			}
		}
	}
}
