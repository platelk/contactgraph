package usecases

import (
	"context"
	"fmt"
	"sort"

	"github.com/platelk/contactgraph/domain/models/contacts"
	"github.com/platelk/contactgraph/domain/models/users"
	"github.com/platelk/contactgraph/infra/logger"
)

type SuggestContactReq struct {
	User users.ID `json:"user_id"`
}

type SuggestContactResp struct {
	Contacts []*users.User
}

type SuggestContact func(ctx context.Context, req *SuggestContactReq) (*SuggestContactResp, error)

func SetupSuggestContact(log logger.Logger, maxReturned uint, lookuper contactLookuper, userStore userGetter) SuggestContact {
	sortFriendByAppearance := generateSortFriendByAppearance(lookuper)
	return func(ctx context.Context, req *SuggestContactReq) (*SuggestContactResp, error) {
		// Retrieve contacts
		contactIDs, err := lookuper.Lookup(ctx, req.User)
		if err != nil {
			return nil, fmt.Errorf("can't lookup user [%v]: %w", req.User, err)
		}
		// Rank friends of friends by number of appearance
		scoredUsers := sortFriendByAppearance(ctx, contactIDs)
		userContacts := make([]*users.User, 0, maxReturned)
		for i := 0; i < len(scoredUsers) && uint(i) < maxReturned; i++ {
			usr, err := userStore.Get(ctx, scoredUsers[i].User)
			if err != nil {
				return nil, fmt.Errorf("can't retrieve user [%v] information: %w", scoredUsers[i].User, err)
			}
			userContacts = append(userContacts, usr)
		}
		return &SuggestContactResp{
			Contacts: userContacts,
		}, nil
	}
}

type userScore struct {
	User  users.ID
	Score uint
}

func generateSortFriendByAppearance(lookuper contactLookuper) func(ctx context.Context, friends contacts.List) []userScore {
	return func(ctx context.Context, friends contacts.List) []userScore {
		topList := make(map[users.ID]uint)
		for contactID := range friends {
			friendContacts, err := lookuper.Lookup(ctx, contactID)
			if err != nil {
				continue
			}
			for friendFriend := range friendContacts {
				topList[friendFriend]++
			}
		}
		//
		var scoredList []userScore
		for suggestion, score := range topList {
			scoredList = append(scoredList, userScore{
				User:  suggestion,
				Score: score,
			})
		}
		sort.SliceStable(scoredList, func(i, j int) bool {
			return scoredList[i].Score < scoredList[j].Score
		})
		return scoredList
	}
}
