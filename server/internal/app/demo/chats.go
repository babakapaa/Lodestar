package demo

import (
	"context"
	"fmt"

	"github.com/babakapaa/lodestar/server/internal/app"
	"github.com/babakapaa/lodestar/server/internal/core/auth"
	"github.com/babakapaa/lodestar/server/internal/modules/chat"
	"github.com/babakapaa/lodestar/server/internal/modules/study"
)

type line struct {
	from auth.Principal
	text string
}

func seedChats(ctx context.Context, core app.Core, people map[string]auth.Principal, courseID string) error {
	chats := chat.NewService(core.Kernel)
	messages := chat.NewMessageService(core.Kernel)
	student, peer, teacher := people[Student.Subject], people[Peer.Subject], people[Teacher.Subject]

	group, _, err := chats.Create(ctx, teacher, chat.ChatInput{
		Kind:      chat.KindGroup,
		Title:     GroupName + ", биология",
		MemberIDs: []string{student.UserID, peer.UserID},
	})
	if err != nil {
		return err
	}
	news, _, err := chats.Create(ctx, teacher, chat.ChatInput{
		Kind:   chat.KindChannel,
		Title:  "Новости курса «Основы биологии»",
		About:  "Объявления преподавателя",
		Origin: &chat.Origin{Kind: study.RosterCourse, ID: courseID},
	})
	if err != nil {
		return err
	}
	direct, _, err := chats.Direct(ctx, student, chat.DirectInput{UserID: peer.UserID})
	if err != nil {
		return err
	}
	if _, err := chats.Accept(ctx, peer, direct.ID); err != nil {
		return err
	}
	for chatID, lines := range map[string][]line{
		news.ID: {
			{teacher, "Добро пожаловать на курс! Материалы лекций уже доступны навигатору, задавайте ему вопросы."},
			{teacher, "Эссе о клетке сдаём через 5 дней, после сдачи будет короткая защита."},
		},
		group.ID: {
			{teacher, "Здесь обсуждаем курс. Пишите, если что-то непонятно."},
			{student, "Ерлан Маратович, карту поиска можно собирать в заметках?"},
			{teacher, "Да, в заметках и в разделе «Источники». Нужно минимум три источника, один научный."},
		},
		direct.ID: {
			{student, "Привет! Пойдёшь со мной в команду на проект по модели клетки?"},
			{peer, "Давай! Могу взять митохондрии и ядро."},
			{student, "Отлично, тогда я соберу мембрану и ЭПС."},
		},
	} {
		for i, l := range lines {
			if _, _, err := messages.Send(ctx, l.from, chatID, chat.MessageInput{
				ClientID: fmt.Sprintf("demo-%s-%d", chatID, i),
				Text:     l.text,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}
