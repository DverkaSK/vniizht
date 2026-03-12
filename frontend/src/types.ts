export type Role = 'GUEST' | 'USER' | 'SPECIALIST' | 'ADMIN'
export type QuestionStatus = 'OPEN' | 'CLOSED' | 'DUPLICATE'
export type VoteValue = 'UP' | 'DOWN'

export interface User {
  id: number
  username: string
  email: string
  role: Role
  reputation: number
  is_active: boolean
  created_at: string
  question_count: number
  answer_count: number
}

export interface PublicUser {
  id: number
  username: string
  role: Role
  reputation: number
  created_at: string
  question_count: number
  answer_count: number
  recent_questions: RecentQuestion[]
  recent_answers: RecentAnswer[]
}

export interface RecentQuestion {
  id: number
  title: string
  status: QuestionStatus
  created_at: string
}

export interface RecentAnswer {
  id: number
  question_id: number
  snippet: string
  created_at: string
}

export interface Tag {
  id: number
  name: string
  created_at: string
}

export interface Category {
  id: number
  name: string
  description: string
  specialist_id?: number
  created_at: string
}

export interface Question {
  id: number
  author_id: number
  author_username: string
  category_id?: number
  specialist_id?: number
  title: string
  body: string
  status: QuestionStatus
  view_count: number
  answer_count: number
  has_verified: boolean
  tags: Tag[]
  created_at: string
  updated_at: string
}

export interface QuestionListResponse {
  items: Question[]
  total: number
  limit: number
  offset: number
}

export interface Answer {
  id: number
  question_id: number
  author_id: number
  author_username: string
  body: string
  is_verified: boolean
  vote_score: number
  current_user_vote?: VoteValue
  created_at: string
  updated_at: string
}

export interface Comment {
  id: number
  answer_id: number
  author_id: number
  author_username: string
  body: string
  created_at: string
  updated_at: string
}

export interface Attachment {
  id: number
  uploader_id: number
  target_type: string
  target_id: number
  filename: string
  mime_type: string
  size_bytes: number
  url: string
}

export interface SearchResult {
  type: 'question' | 'answer'
  answer_id?: number
  question_id: number
  question_title: string
  snippet: string
  status: QuestionStatus
  author_id: number
  category_id?: number
  created_at: string
  rank: number
}

export interface SearchResponse {
  items: SearchResult[]
  total: number
}
