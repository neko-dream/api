package db

import (
	"context"
	"log"

	"github.com/neko-dream/server/internal/domain/model/opinion"
	"github.com/neko-dream/server/internal/domain/model/shared"
	"github.com/neko-dream/server/internal/domain/model/shared/time"
	talksession "github.com/neko-dream/server/internal/domain/model/talk_session"
	"github.com/neko-dream/server/internal/domain/model/user"
	"github.com/neko-dream/server/internal/domain/model/vote"
	"github.com/samber/lo"
)

type DummyInitializer struct {
	*DBManager
	UserRepo        user.UserRepository
	TalkSessionRepo talksession.TalkSessionRepository
	OpinionRepo     opinion.OpinionRepository
	VoteRepo        vote.VoteRepository

	TalkSessions []*talksession.TalkSession
	Users        []*user.User
	Opinions     []*opinion.Opinion
	Votes        []*vote.Vote
}

func NewDummyInitializer(
	dbManager *DBManager,
	userRepo user.UserRepository,
	talkSessionRepo talksession.TalkSessionRepository,
	opinionRepo opinion.OpinionRepository,
	voteRepo vote.VoteRepository,
) *DummyInitializer {
	return &DummyInitializer{
		DBManager:       dbManager,
		UserRepo:        userRepo,
		TalkSessionRepo: talkSessionRepo,
		OpinionRepo:     opinionRepo,
		VoteRepo:        voteRepo,
	}
}

func (i *DummyInitializer) Initialize() {
	log.Println("-------- Start DummyInitializer Initialize --------")
	_ = i.User()
	_ = i.TalkSession()
	_ = i.Opinion()
	log.Println("-------- End DummyInitializer Initialize --------")
}

func (d *DummyInitializer) User() error {
	users := []user.User{
		// 否定派閥
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user1"),
			lo.ToPtr("オブジェクト指向大好きマン"),
			"user1",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user2"),
			lo.ToPtr("手続き型よかまし"),
			"user2",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user3"),
			lo.ToPtr("<script>alert('test')</script>"),
			"user3",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user4"),
			lo.ToPtr("hogehoge' SELECT * FROM users; --"),
			"user4",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user5"),
			lo.ToPtr("関数型至上主義"),
			"user5",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user6"),
			lo.ToPtr("user6"),
			"user6",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user7"),
			lo.ToPtr("user7"),
			"user7",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user8"),
			lo.ToPtr("user8"),
			"user8",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user9"),
			lo.ToPtr("user9"),
			"user9",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user10"),
			lo.ToPtr("user10"),
			"user10",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user11"),
			lo.ToPtr("user11"),
			"user11",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user12"),
			lo.ToPtr("user12"),
			"user12",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user13"),
			lo.ToPtr("user13"),
			"user13",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user14"),
			lo.ToPtr("user14"),
			"user14",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user15"),
			lo.ToPtr("user15"),
			"user15",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user16"),
			lo.ToPtr("user16"),
			"user16",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user17"),
			lo.ToPtr("user17"),
			"user17",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user18"),
			lo.ToPtr("user18"),
			"user18",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user19"),
			lo.ToPtr("user19"),
			"user19",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user20"),
			lo.ToPtr("user20"),
			"user20",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user21"),
			lo.ToPtr("user21"),
			"user21",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user22"),
			lo.ToPtr("user22"),
			"user22",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user23"),
			lo.ToPtr("user23"),
			"user23",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user24"),
			lo.ToPtr("user24"),
			"user24",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user25"),
			lo.ToPtr("user25"),
			"user25",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user126"),
			lo.ToPtr("user126"),
			"user126",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user127"),
			lo.ToPtr("user127"),
			"user127",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user128"),
			lo.ToPtr("user128"),
			"user128",
			"GOOGLE",
			user.NewProfileIcon(
				lo.ToPtr("https://images.kotohiro.com/users/0192521b-136d-7543-81f9-fc38cd16023f/profile_icon/1728037600.jpg"),
			),
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user129"),
			lo.ToPtr("user129"),
			"user129",
			"GOOGLE",
			nil,
		),
		user.NewUser(
			shared.NewUUID[user.User](),
			lo.ToPtr("user220"),
			lo.ToPtr("user220"),
			"user220",
			"GOOGLE",
			nil,
		),
	}

	userDemographics := []*user.UserDemographics{
		nil,
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)),
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			user.NewYearOfBirth(lo.ToPtr(1990)),
			user.NewOccupation(lo.ToPtr("会社員")),
			lo.ToPtr(user.NewGender(lo.ToPtr("男性"))),
			user.NewMunicipality(lo.ToPtr("中野区")),
			user.NewHouseholdSize(lo.ToPtr(1)),
			lo.ToPtr("東京都"),
		)),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)),
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			user.NewYearOfBirth(lo.ToPtr(1990)),
			user.NewOccupation(lo.ToPtr("会社員")),
			lo.ToPtr(user.NewGender(lo.ToPtr("男性"))),
			user.NewMunicipality(lo.ToPtr("中野区")),
			user.NewHouseholdSize(lo.ToPtr(1)),
			lo.ToPtr("東京都"),
		)),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		)),
		lo.ToPtr(user.NewUserDemographics(
			shared.NewUUID[user.UserDemographics](),
			user.NewYearOfBirth(lo.ToPtr(1990)),
			user.NewOccupation(lo.ToPtr("会社員")),
			lo.ToPtr(user.NewGender(lo.ToPtr("男性"))),
			user.NewMunicipality(lo.ToPtr("中野区")),
			user.NewHouseholdSize(lo.ToPtr(1)),
			lo.ToPtr("東京都"),
		)),
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	}
	var resultUsers []*user.User
	ctx := context.Background()
	for i, u := range users {
		err := d.UserRepo.Create(ctx, u)
		if err != nil {
			return err
		}
		demographics := userDemographics[i]
		if demographics != nil {
			u.SetDemographics(*demographics)
		}
		err = d.UserRepo.Update(ctx, u)
		if err != nil {
			return err
		}
		resultUsers = append(resultUsers, &u)
	}
	d.Users = resultUsers

	return nil
}

func (d *DummyInitializer) TalkSession() error {
	ctx := context.Background()
	talkSessions := []*talksession.TalkSession{
		talksession.NewTalkSession(
			shared.NewUUID[talksession.TalkSession](),
			"オブジェクト指向は悪",
			d.Users[0].UserID(),
			time.Now(ctx),
			time.Now(ctx).Add(ctx, 3*time.Month),
			nil,
		),
		talksession.NewTalkSession(
			shared.NewUUID[talksession.TalkSession](),
			"オブジェクト指向は良",
			d.Users[1].UserID(),
			time.Now(ctx),
			time.Now(ctx).Add(ctx, 3*time.Month),
			nil,
		),
	}

	var resultTalkSessions []*talksession.TalkSession
	for _, ts := range talkSessions {
		err := d.TalkSessionRepo.Create(ctx, ts)
		if err != nil {
			return err
		}
		resultTalkSessions = append(resultTalkSessions, ts)
	}
	d.TalkSessions = resultTalkSessions
	return nil
}

func (d *DummyInitializer) Opinion() error {
	ctx := context.Background()
	// 1つ目に集める
	ts := d.TalkSessions[0]

	var objectGroup []*user.User
	var functionalGroup []*user.User
	var proceduralGroup []*user.User

	objectGroup = append(objectGroup, d.Users[0], d.Users[1], d.Users[3], d.Users[6], d.Users[8])
	functionalGroup = append(functionalGroup, d.Users[1], d.Users[9], d.Users[5], d.Users[9])
	proceduralGroup = append(proceduralGroup, d.Users[2], d.Users[4], d.Users[7])
	objectGroup = append(objectGroup, d.Users[10:17]...)
	functionalGroup = append(functionalGroup, d.Users[17:25]...)
	proceduralGroup = append(proceduralGroup, d.Users[25:30]...)

	var opinions []*opinion.Opinion
	var votes []*vote.Vote
	o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[0].UserID(),
		nil,
		lo.ToPtr("オブジェクト指向は最高！"),
		"オブジェクト指向は現実世界をモデル化できる最高の方法！効率的だし、どんな規模でも対応できる！",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o1)
	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o1o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[2].UserID(),
		lo.ToPtr(o1.OpinionID()),
		nil,
		"現実をモデル化したとて複雑化するだけ。オブジェクト指向は時代遅れ。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o1o1)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o1o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o2, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[1].UserID(),
		nil,
		lo.ToPtr("手続型よりマシ"),
		"クラスとオブジェクトの概念がなかったら、大規模システムなんて絶対崩壊してるよ。手続き型で管理できるわけない。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o2)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o2o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[4].UserID(),
		lo.ToPtr(o2.OpinionID()),
		nil,
		"別にオブジェクト指向がなくても、大規模システムは作れる。大規模システムでオブジェクト指向を使っても崩壊することはある。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o2o1)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o2o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o3, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[2].UserID(),
		nil,
		nil,
		"オブジェクト指向はクラスの継承でコードがカオスになる。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o3)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o3.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o3.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o3.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o4, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[3].UserID(),
		nil,
		nil,
		"オブジェクト指向は分かりやすいし、チーム開発でもコミュニケーションがスムーズになるから最適。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o4)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o4o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[4].UserID(),
		lo.ToPtr(o4.OpinionID()),
		nil,
		"オブジェクト指向は正しく設計されていないとあまりにもわかりにくい。設計が重要。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o4o1)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o4o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o5, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[4].UserID(),
		nil,
		lo.ToPtr("オブジェクト指向は時代遅れ"),
		"オブジェクト指向なんて時代遅れだよ。状態管理が複雑すぎるし、バグの温床。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o5)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o5.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o5.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o5.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o6, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[5].UserID(),
		nil,
		nil,
		"関数型だけが正義。オブジェクト指向は状態管理が難しすぎる。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o6)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o6.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o6.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o6.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o7, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[6].UserID(),
		nil,
		nil,
		"理論的に見ても、オブジェクト指向は現実世界のシミュレーションに最も近い。これを捨てるなんて非合理的。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o7)

	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o7o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[7].UserID(),
		lo.ToPtr(o7.OpinionID()),
		nil,
		"理想論でしかない。現実世界をシミュレーションできて何が嬉しいのか。そもそもコードが現実世界をシミュレーションする必要があるのか。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o7o1)
	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o7o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o8, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[7].UserID(),
		nil,
		nil,
		"オブジェクト指向使ってるけど、正直メンテコストばかり増える気がする。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o8)
	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Pass,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o8o1, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[8].UserID(),
		lo.ToPtr(o8.OpinionID()),
		nil,
		"適切に設計されていればメンテコストは増えない。設計を正しくできない人間が文句言ってるだけ。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o8o1)
	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o8o1.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o9, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[8].UserID(),
		nil,
		nil,
		"オブジェクト指向は、継承の概念があるから、コードの再利用がしやすい。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o9)
	for _, u := range objectGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o9.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range functionalGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o9.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}
	for _, u := range proceduralGroup {
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o9.OpinionID(),
			ts.TalkSessionID(),
			u.UserID(),
			vote.Disagreed,
			time.Now(ctx).Time,
		)
		votes = append(votes, vs)
	}

	o10, _ := opinion.NewOpinion(
		shared.NewUUID[opinion.Opinion](),
		ts.TalkSessionID(),
		d.Users[9].UserID(),
		nil,
		nil,
		"関数型言語の方が、オブジェクト指向よりも再利用性が高い。",
		time.Now(ctx).Time,
	)
	opinions = append(opinions, o10)

	for _, o := range opinions {
		err := d.OpinionRepo.Create(ctx, *o)
		if err != nil {
			return err
		}
		// 自分の意見には必ず投票を紐付ける
		vs, _ := vote.NewVote(
			shared.NewUUID[vote.Vote](),
			o.OpinionID(),
			ts.TalkSessionID(),
			o.UserID(),
			vote.Agreed,
			time.Now(ctx).Time,
		)
		err = d.VoteRepo.Create(ctx, *vs)
		if err != nil {
			return err
		}

		d.Opinions = append(d.Opinions, o)
		d.Votes = append(d.Votes, vs)
	}
	for _, vs := range votes {
		_ = d.VoteRepo.Create(ctx, *vs)
	}

	return nil
}
