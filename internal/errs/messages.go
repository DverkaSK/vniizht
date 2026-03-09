package errs

const (
	InternalError  = "внутренняя ошибка сервера"
	InvalidBody    = "некорректное тело запроса"
	NotImplemented = "не реализовано"
)

const (
	InvalidCredentials = "неверный логин или пароль"
	AccountDisabled    = "учётная запись отключена"
	Unauthorized       = "необходима авторизация"
	Forbidden          = "недостаточно прав"
	SessionExpired     = "сессия истекла"
)

const (
	CommentNotFound     = "комментарий не найден"
	CommentForbidden    = "нет доступа к этому комментарию"
	CommentBodyRequired = "текст комментария обязателен"
	CommentInvalidID    = "некорректный идентификатор комментария"
	CommentCreateError  = "ошибка создания комментария"
	CommentUpdateError  = "ошибка обновления комментария"
	CommentDeleteError  = "ошибка удаления комментария"
)

const (
	AnswerNotFound      = "ответ не найден"
	AnswerForbidden     = "нет доступа к этому ответу"
	AnswerBodyRequired  = "текст ответа обязателен"
	AnswerInvalidID     = "некорректный идентификатор ответа"
	AnswerListError     = "ошибка получения ответов"
	AnswerCreateError   = "ошибка создания ответа"
	AnswerUpdateError   = "ошибка обновления ответа"
	AnswerDeleteError   = "ошибка удаления ответа"
	AnswerVerifyError   = "ошибка подтверждения ответа"
	AnswerUnverifyError = "ошибка снятия подтверждения ответа"
	AnswerInvalidVote   = "недопустимое значение голоса"
	AnswerSelfVote      = "нельзя голосовать за собственный ответ"
	AnswerVoteExists    = "оценка уже существует"
	AnswerVoteError     = "ошибка голосования за ответ"
)

const (
	QuestionNotFound       = "вопрос не найден"
	QuestionForbidden      = "нет доступа к этому вопросу"
	QuestionFieldsRequired = "заголовок и текст вопроса обязательны"
	QuestionInvalidID      = "некорректный идентификатор вопроса"
	QuestionListError      = "ошибка получения списка вопросов"
	QuestionCreateError    = "ошибка создания вопроса"
	QuestionUpdateError    = "ошибка обновления вопроса"
	QuestionCloseError     = "ошибка закрытия вопроса"
	QuestionAlreadyClosed  = "вопрос уже закрыт"
	QuestionInvalidStatus  = "недопустимый статус вопроса"
	QuestionInvalidPage    = "некорректный параметр страницы"
)

const (
	CategoryNotFound     = "категория не найдена"
	CategoryDuplicate    = "категория с таким названием уже существует"
	CategoryNameRequired = "название категории обязательно"
	CategoryListError    = "ошибка получения списка категорий"
	CategoryCreateError  = "ошибка создания категории"
	CategoryUpdateError  = "ошибка обновления категории"
	CategoryDeleteError  = "ошибка удаления категории"
	CategoryInvalidID    = "некорректный идентификатор категории"
)

const (
	TagNotFound     = "тег не найден"
	TagDuplicate    = "тег с таким названием уже существует"
	TagNameRequired = "название тега обязательно"
	TagListError    = "ошибка получения списка тегов"
	TagCreateError  = "ошибка создания тега"
	TagUpdateError  = "ошибка обновления тега"
	TagDeleteError  = "ошибка удаления тега"
	TagInvalidID    = "некорректный идентификатор тега"
)

const (
	AttachmentNotFound      = "вложение не найдено"
	AttachmentForbidden     = "нет доступа к этому вложению"
	AttachmentTooLarge      = "файл превышает максимальный размер 50 МБ"
	AttachmentNoFile        = "файл не передан"
	AttachmentInvalidID     = "некорректный идентификатор вложения"
	AttachmentInvalidTarget = "некорректный тип или идентификатор объекта"
	AttachmentUploadError   = "ошибка загрузки файла"
	AttachmentDeleteError   = "ошибка удаления файла"
)

const (
	SearchQueryRequired = "поисковый запрос не может быть пустым"
	SearchError         = "ошибка выполнения поиска"
)

const (
	ProfileError = "ошибка получения профиля"
)

const (
	UserNotFound        = "пользователь не найден"
	UserDuplicate       = "пользователь с таким логином или email уже существует"
	UserFieldsRequired  = "логин, email и пароль обязательны"
	UserInvalidRole     = "недопустимая роль"
	UserInvalidID       = "некорректный идентификатор пользователя"
	UserHashError       = "ошибка хеширования пароля"
	UserListError       = "ошибка получения списка пользователей"
	UserCreateError     = "ошибка создания пользователя"
	UserUpdateRoleError = "ошибка изменения роли"
)
