package setup

type LegalQuestion struct {
	Key   string
	Title string
}

var LegalQuestions = []LegalQuestion{
	{Key: "storage", Title: "Страна вуза и место хранения данных"},
	{Key: "consent", Title: "Согласие на обработку данных и согласие родителей несовершеннолетних"},
	{Key: "red_button", Title: "Красная кнопка навигатора: кому сообщает и на каком основании"},
	{Key: "recording", Title: "Запись защит и звонков"},
	{Key: "retention", Title: "Сроки хранения работ, оценок и личных данных"},
	{Key: "personal_access", Title: "Процедура доступа к личным данным"},
	{Key: "sandbox_watch", Title: "Надзор в песочнице и публичный хостинг демо"},
	{Key: "data_officer", Title: "Ответственный за данные и контакт для жалоб"},
}

func knownQuestion(key string) bool {
	for _, q := range LegalQuestions {
		if q.Key == key {
			return true
		}
	}
	return false
}
