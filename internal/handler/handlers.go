package handler

type Handlers struct {
	Auth        *AuthHandler
	Questions   *QuestionsHandler
	Answers     *AnswersHandler
	Comments    *CommentsHandler
	Attachments *AttachmentsHandler
	Users       *UsersHandler
	Search      *SearchHandler
	Admin       *AdminHandler
}
