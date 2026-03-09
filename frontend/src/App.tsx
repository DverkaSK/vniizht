import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { AuthProvider } from '@/context/AuthContext'
import { Navbar } from '@/components/layout/Navbar'
import { HomePage } from '@/pages/HomePage'
import { QuestionPage } from '@/pages/QuestionPage'
import { LoginPage } from '@/pages/LoginPage'
import { ProfilePage } from '@/pages/ProfilePage'
import { SearchPage } from '@/pages/SearchPage'
import { CreateQuestionPage } from '@/pages/CreateQuestionPage'

function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <main>{children}</main>
    </div>
  )
}

export default function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/"
            element={
              <Layout>
                <HomePage />
              </Layout>
            }
          />
          <Route
            path="/questions/new"
            element={
              <Layout>
                <CreateQuestionPage />
              </Layout>
            }
          />
          <Route
            path="/questions/:id"
            element={
              <Layout>
                <QuestionPage />
              </Layout>
            }
          />
          <Route
            path="/search"
            element={
              <Layout>
                <SearchPage />
              </Layout>
            }
          />
          <Route
            path="/users/:id"
            element={
              <Layout>
                <ProfilePage />
              </Layout>
            }
          />
          <Route
            path="*"
            element={
              <Layout>
                <div className="flex items-center justify-center min-h-[60vh] text-muted-foreground">
                  <div className="text-center">
                    <p className="text-4xl font-bold mb-2">404</p>
                    <p>Страница не найдена</p>
                  </div>
                </div>
              </Layout>
            }
          />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  )
}
