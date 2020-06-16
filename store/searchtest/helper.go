// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package searchtest

import (
	"fmt"
	"testing"

	"github.com/telluria-bruno-gouvea/mattermost-server/v5/model"
	"github.com/telluria-bruno-gouvea/mattermost-server/v5/store"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type SearchTestHelper struct {
	Store              store.Store
	Team               *model.Team
	AnotherTeam        *model.Team
	User               *model.User
	User2              *model.User
	UserAnotherTeam    *model.User
	ChannelBasic       *model.Channel
	ChannelPrivate     *model.Channel
	ChannelAnotherTeam *model.Channel
	ChannelDeleted     *model.Channel
}

func (th *SearchTestHelper) SetupBasicFixtures() error {
	// Create teams
	team, err := th.createTeam("searchtest-team", "Searchtest team", model.TEAM_OPEN)
	if err != nil {
		return err
	}
	anotherTeam, err := th.createTeam("another-searchtest-team", "Another Searchtest team", model.TEAM_OPEN)
	if err != nil {
		return err
	}

	// Create users
	user, err := th.createUser("basicusername1", "basicnickname1", "basicfirstname1", "basiclastname1")
	if err != nil {
		return err
	}
	user2, err := th.createUser("basicusername2", "basicnickname2", "basicfirstname2", "basiclastname2")
	if err != nil {
		return err
	}
	useranother, err := th.createUser("basicusername3", "basicnickname3", "basicfirstname3", "basiclastname3")
	if err != nil {
		return err
	}

	// Create channels
	channelBasic, err := th.createChannel(team.Id, "channel-a", "ChannelA", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}
	channelPrivate, err := th.createChannel(team.Id, "channel-private", "ChannelPrivate", "", model.CHANNEL_PRIVATE, false)
	if err != nil {
		return err
	}
	channelDeleted, err := th.createChannel(team.Id, "channel-deleted", "ChannelA (deleted)", "", model.CHANNEL_OPEN, true)
	if err != nil {
		return err
	}
	channelAnotherTeam, err := th.createChannel(anotherTeam.Id, "channel-a", "ChannelA", "", model.CHANNEL_OPEN, false)
	if err != nil {
		return err
	}

	err = th.addUserToTeams(user, []string{team.Id, anotherTeam.Id})
	if err != nil {
		return err
	}

	err = th.addUserToTeams(user2, []string{team.Id, anotherTeam.Id})
	if err != nil {
		return err
	}

	err = th.addUserToTeams(useranother, []string{anotherTeam.Id})
	if err != nil {
		return err
	}

	_, err = th.addUserToChannels(user, []string{channelBasic.Id, channelPrivate.Id, channelDeleted.Id})
	if err != nil {
		return err
	}

	_, err = th.addUserToChannels(user2, []string{channelPrivate.Id, channelDeleted.Id})
	if err != nil {
		return err
	}

	_, err = th.addUserToChannels(useranother, []string{channelAnotherTeam.Id})
	if err != nil {
		return err
	}

	th.Team = team
	th.AnotherTeam = anotherTeam
	th.User = user
	th.User2 = user2
	th.UserAnotherTeam = useranother
	th.ChannelBasic = channelBasic
	th.ChannelPrivate = channelPrivate
	th.ChannelAnotherTeam = channelAnotherTeam
	th.ChannelDeleted = channelDeleted

	return nil
}

func (th *SearchTestHelper) CleanFixtures() error {
	err := th.deleteChannels([]*model.Channel{
		th.ChannelBasic, th.ChannelPrivate, th.ChannelAnotherTeam, th.ChannelDeleted,
	})
	if err != nil {
		return err
	}

	err = th.deleteTeam(th.Team)
	if err != nil {
		return err
	}

	err = th.deleteTeam(th.AnotherTeam)
	if err != nil {
		return err
	}

	err = th.deleteUser(th.User)
	if err != nil {
		return err
	}

	err = th.deleteUser(th.User2)
	if err != nil {
		return err
	}

	err = th.deleteUser(th.UserAnotherTeam)
	if err != nil {
		return err
	}

	return nil
}

func (th *SearchTestHelper) createTeam(name, displayName, teamType string) (*model.Team, error) {
	team, appError := th.Store.Team().Save(&model.Team{
		Name:        name,
		DisplayName: displayName,
		Type:        teamType,
	})
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return team, nil
}

func (th *SearchTestHelper) deleteTeam(team *model.Team) error {
	appError := th.Store.Team().RemoveAllMembersByTeam(team.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}
	appError = th.Store.Team().PermanentDelete(team.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	return nil
}

func (th *SearchTestHelper) makeEmail() string {
	return "success_" + model.NewId() + "@simulator.amazon.com"
}

func (th *SearchTestHelper) createUser(username, nickname, firstName, lastName string) (*model.User, error) {
	user, appError := th.Store.User().Save(&model.User{
		Username:  username,
		Password:  username,
		Nickname:  nickname,
		FirstName: firstName,
		LastName:  lastName,
		Email:     th.makeEmail(),
	})
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return user, nil
}

func (th *SearchTestHelper) deleteUser(user *model.User) error {
	appError := th.Store.User().PermanentDelete(user.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	return nil
}

func (th *SearchTestHelper) createBot(username, displayName, ownerID string) (*model.Bot, error) {
	botModel := &model.Bot{
		Username:    username,
		DisplayName: displayName,
		OwnerId:     ownerID,
	}

	user, apperr := th.Store.User().Save(model.UserFromBot(botModel))
	if apperr != nil {
		return nil, errors.New(apperr.Error())
	}

	botModel.UserId = user.Id
	bot, err := th.Store.Bot().Save(botModel)
	if err != nil {
		th.Store.User().PermanentDelete(bot.UserId)
		return nil, errors.New(err.Error())
	}

	return bot, nil
}

func (th *SearchTestHelper) deleteBot(botID string) error {
	err := th.Store.Bot().PermanentDelete(botID)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (th *SearchTestHelper) createChannel(teamID, name, displayName, purpose, channelType string, deleted bool) (*model.Channel, error) {
	channel, appError := th.Store.Channel().Save(&model.Channel{
		TeamId:      teamID,
		DisplayName: displayName,
		Name:        name,
		Type:        channelType,
		Purpose:     purpose,
	}, 999)
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	if deleted {
		appError := th.Store.Channel().Delete(channel.Id, model.GetMillis())
		if appError != nil {
			return nil, errors.New(appError.Error())
		}
	}

	return channel, nil
}

func (th *SearchTestHelper) createDirectChannel(teamID, name, displayName string, users []*model.User) (*model.Channel, error) {
	channel := &model.Channel{
		TeamId:      teamID,
		Name:        name,
		DisplayName: displayName,
		Type:        model.CHANNEL_DIRECT,
	}

	m1 := &model.ChannelMember{}
	m1.ChannelId = channel.Id
	m1.UserId = users[0].Id
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()

	m2 := &model.ChannelMember{}
	m2.ChannelId = channel.Id
	m2.UserId = users[0].Id
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()

	channel, err := th.Store.Channel().SaveDirectChannel(channel, m1, m2)
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (th *SearchTestHelper) createGroupChannel(teamID, displayName string, users []*model.User) (*model.Channel, error) {
	userIDS := make([]string, len(users))
	for _, user := range users {
		userIDS = append(userIDS, user.Id)
	}

	group := &model.Channel{
		TeamId:      teamID,
		Name:        model.GetGroupNameFromUserIds(userIDS),
		DisplayName: displayName,
		Type:        model.CHANNEL_GROUP,
	}

	channel, apperr := th.Store.Channel().Save(group, 10000)
	if apperr != nil {
		return nil, errors.New(apperr.Error())
	}

	for _, user := range users {
		_, err := th.addUserToChannels(user, []string{channel.Id})
		if err != nil {
			return nil, err
		}
	}

	return channel, nil

}

func (th *SearchTestHelper) deleteChannel(channel *model.Channel) error {
	appError := th.Store.Channel().PermanentDeleteMembersByChannel(channel.Id)
	if appError != nil {
		return errors.New(appError.Error())
	}

	err := th.Store.Channel().PermanentDelete(channel.Id)
	if err != nil {
		return err
	}

	return nil
}

func (th *SearchTestHelper) deleteChannels(channels []*model.Channel) error {
	for _, channel := range channels {
		err := th.deleteChannel(channel)
		if err != nil {
			return err
		}
	}

	return nil
}

func (th *SearchTestHelper) createPostModel(userID, channelID, message, hashtags, postType string, createAt int64, pinned bool) *model.Post {
	return &model.Post{
		Message:       message,
		ChannelId:     channelID,
		PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
		UserId:        userID,
		Hashtags:      hashtags,
		IsPinned:      pinned,
		CreateAt:      createAt,
		Type:          postType,
	}
}

func (th *SearchTestHelper) createPost(userID, channelID, message, hashtags, postType string, createAt int64, pinned bool) (*model.Post, error) {
	var creationTime int64 = 1000000
	if createAt > 0 {
		creationTime = createAt
	}
	postModel := th.createPostModel(userID, channelID, message, hashtags, postType, creationTime, pinned)
	post, appError := th.Store.Post().Save(postModel)
	if appError != nil {
		return nil, errors.New(appError.Error())
	}

	return post, nil
}

func (th *SearchTestHelper) createReply(userID, message, hashtags string, parent *model.Post, createAt int64, pinned bool) (*model.Post, error) {
	replyModel := th.createPostModel(userID, parent.ChannelId, message, hashtags, parent.Type, createAt, pinned)
	replyModel.ParentId = parent.Id
	replyModel.RootId = parent.Id
	reply, appError := th.Store.Post().Save(replyModel)
	if appError != nil {
		return nil, errors.New(appError.Error())
	}
	return reply, nil
}

func (th *SearchTestHelper) deleteUserPosts(userID string) error {
	err := th.Store.Post().PermanentDeleteByUser(userID)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (th *SearchTestHelper) addUserToTeams(user *model.User, teamIDS []string) error {
	for _, teamID := range teamIDS {
		_, err := th.Store.Team().SaveMember(&model.TeamMember{TeamId: teamID, UserId: user.Id}, -1)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	return nil
}

func (th *SearchTestHelper) addUserToChannels(user *model.User, channelIDS []string) ([]*model.ChannelMember, error) {

	channelMembers := make([]*model.ChannelMember, len(channelIDS))
	for _, channelID := range channelIDS {
		cm, err := th.Store.Channel().SaveMember(&model.ChannelMember{
			ChannelId:   channelID,
			UserId:      user.Id,
			NotifyProps: model.GetDefaultChannelNotifyProps(),
		})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		channelMembers = append(channelMembers, cm)
	}

	return channelMembers, nil
}

func (th *SearchTestHelper) assertUsersMatchInAnyOrder(t *testing.T, expected, actual []*model.User) {
	expectedUsernames := make([]string, 0, len(expected))
	for _, user := range expected {
		user.Sanitize(map[string]bool{})
		expectedUsernames = append(expectedUsernames, user.Username)
	}

	actualUsernames := make([]string, 0, len(actual))
	for _, user := range actual {
		user.Sanitize(map[string]bool{})
		actualUsernames = append(actualUsernames, user.Username)
	}

	if assert.ElementsMatch(t, expectedUsernames, actualUsernames) {
		assert.ElementsMatch(t, expected, actual)
	}
}

func (th *SearchTestHelper) checkPostInSearchResults(t *testing.T, postID string, searchResults map[string]*model.Post) {
	t.Helper()
	postIDS := make([]string, len(searchResults))
	for ID := range searchResults {
		postIDS = append(postIDS, ID)
	}
	assert.Contains(t, postIDS, postID, "Did not find expected post in search results.")
}

func (th *SearchTestHelper) checkChannelIdsMatch(t *testing.T, expected []string, results *model.ChannelList) {
	t.Helper()
	channelIds := make([]string, len(*results))
	for i, channel := range *results {
		channelIds[i] = channel.Id
	}
	require.ElementsMatch(t, expected, channelIds)
}

type ByChannelDisplayName model.ChannelList

func (s ByChannelDisplayName) Len() int { return len(s) }
func (s ByChannelDisplayName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByChannelDisplayName) Less(i, j int) bool {
	if s[i].DisplayName != s[j].DisplayName {
		return s[i].DisplayName < s[j].DisplayName
	}

	return s[i].Id < s[j].Id
}
