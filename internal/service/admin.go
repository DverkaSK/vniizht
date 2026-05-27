package service

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"vniizht/internal/errs"
	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type AdminService struct {
	users      *repository.UserRepo
	categories *repository.CategoryRepo
	tags       *repository.TagRepo
}

func NewAdminService(users *repository.UserRepo, categories *repository.CategoryRepo, tags *repository.TagRepo) *AdminService {
	return &AdminService{users: users, categories: categories, tags: tags}
}

func (s *AdminService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.users.List(ctx)
}

func (s *AdminService) CreateUser(ctx context.Context, username, email, password string, role model.Role) (*model.User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}
	if err := s.users.Create(ctx, u); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (s *AdminService) ChangeRole(ctx context.Context, id int64, role model.Role) error {
	return s.users.UpdateRole(ctx, id, role)
}

func (s *AdminService) SetUserActive(ctx context.Context, id int64, active bool) error {
	return s.users.SetActive(ctx, id, active)
}

func (s *AdminService) ListCategories(ctx context.Context) ([]*model.Category, error) {
	return s.categories.List(ctx)
}

func (s *AdminService) CreateCategory(ctx context.Context, name, description string, specialistID *int64) (*model.Category, error) {
	category := &model.Category{Name: name, Description: description, SpecialistID: specialistID}
	if err := s.categories.Create(ctx, category); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (s *AdminService) UpdateCategory(ctx context.Context, id int64, name, description string, specialistID *int64) error {
	category := &model.Category{ID: id, Name: name, Description: description, SpecialistID: specialistID}
	if err := s.categories.Update(ctx, category); err != nil {
		if isDuplicateAdminErr(err) {
			return errs.ErrConflict
		}
		return err
	}
	return nil
}

func (s *AdminService) DeleteCategory(ctx context.Context, id int64) error {
	return s.categories.Delete(ctx, id)
}

func (s *AdminService) ListTags(ctx context.Context) ([]*model.Tag, error) {
	return s.tags.List(ctx)
}

func (s *AdminService) CreateTag(ctx context.Context, name string) (*model.Tag, error) {
	tag := &model.Tag{Name: name}
	if err := s.tags.Create(ctx, tag); err != nil {
		if isDuplicateAdminErr(err) {
			return nil, errs.ErrConflict
		}
		return nil, fmt.Errorf("create tag: %w", err)
	}
	return tag, nil
}

func (s *AdminService) UpdateTag(ctx context.Context, id int64, name string) error {
	tag := &model.Tag{ID: id, Name: name}
	if err := s.tags.Update(ctx, tag); err != nil {
		if isDuplicateAdminErr(err) {
			return errs.ErrConflict
		}
		return err
	}
	return nil
}

func (s *AdminService) DeleteTag(ctx context.Context, id int64) error {
	return s.tags.Delete(ctx, id)
}

// SuggestResult содержит автоподобранные категорию и теги.
type SuggestResult struct {
	CategoryID *int64  `json:"category_id"`
	TagIDs     []int64 `json:"tag_ids"`
}

// Suggest возвращает наиболее подходящую категорию и теги по тексту заголовка вопроса.
// Использует посимвольное совпадение слов (без ML): для категорий берётся лучший результат,
// для тегов — все совпавшие.
func (s *AdminService) Suggest(ctx context.Context, query string) (*SuggestResult, error) {
	categories, err := s.categories.List(ctx)
	if err != nil {
		return nil, err
	}
	tags, err := s.tags.List(ctx)
	if err != nil {
		return nil, err
	}

	result := &SuggestResult{TagIDs: []int64{}}
	words := suggestTokenize(query)
	if len(words) == 0 {
		return result, nil
	}

	// Наилучшая категория по числу совпадений в name + description
	bestScore := 0
	for _, c := range categories {
		score := suggestScore(words, c.Name+" "+c.Description)
		if score > bestScore {
			bestScore = score
			id := c.ID
			result.CategoryID = &id
		}
	}

	// Все теги, у которых есть хотя бы одно совпадение
	for _, t := range tags {
		if suggestScore(words, t.Name) > 0 {
			result.TagIDs = append(result.TagIDs, t.ID)
		}
	}
	return result, nil
}

// suggestTokenize разбивает строку на слова длиной ≥ 3 символов (строчные).
func suggestTokenize(s string) []string {
	s = strings.ToLower(s)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := make([]string, 0, len(words))
	for _, w := range words {
		if len([]rune(w)) >= 3 {
			out = append(out, w)
		}
	}
	return out
}

// suggestScore считает, сколько слов из words совпадают с любым токеном в target.
// Совпадение: точная подстрока ИЛИ общий префикс ≥ 4 символов, покрывающий ≥ 70% более
// короткого слова (обрабатывает русские падежные окончания: «нитки»/«нитках», «окна»/«окнах»).
func suggestScore(words []string, target string) int {
	target = strings.ToLower(target)
	targetTokens := suggestTokenize(target)
	score := 0
	for _, w := range words {
		if suggestWordMatches(w, target, targetTokens) {
			score++
		}
	}
	return score
}

// suggestWordMatches проверяет совпадение одного слова с целевой строкой.
func suggestWordMatches(word, targetFull string, targetTokens []string) bool {
	// Быстрый путь: точная подстрока
	if strings.Contains(targetFull, word) {
		return true
	}
	// Нечёткое совпадение по префиксу с каждым токеном целевой строки
	wr := []rune(word)
	for _, tok := range targetTokens {
		tr := []rune(tok)
		pl := suggestCommonPrefixLen(wr, tr)
		if pl < 4 {
			continue
		}
		shorter := len(wr)
		if len(tr) < shorter {
			shorter = len(tr)
		}
		if float64(pl)/float64(shorter) >= 0.70 {
			return true
		}
	}
	return false
}

// suggestCommonPrefixLen возвращает длину общего префикса двух rune-слайсов.
func suggestCommonPrefixLen(a, b []rune) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func isDuplicateAdminErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "unique") || strings.Contains(msg, "duplicate")
}
