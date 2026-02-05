import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css'
import Welcome from './pages/Welcome';
import NotFound from './pages/NotFound';
import { ProblemsPage } from './pages/Problems';
import NavBar from './components/NavBar/NavBar';
import { Theme } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import { ViewProblemPage } from './pages/ViewProblem';
import { CreateProblemPage } from './pages/CreateProblem';
import { ToastProvider } from './components/Toast/ToastContext';
import { QuizPage } from './pages/QuizPage';

function App() {

  return (
    <Router>
      <ToastProvider>
        <div className='App'>
          <Theme>
            <NavBar />
            <div className="container mt-4">
              <div className='App-header'>
                <h1>Quizgame</h1>
              </div>

              <Routes>
                <Route path="/" element={<Welcome />} />
                <Route path="/*" element={<NotFound />} />
                <Route path="/problems" element={<ProblemsPage />} />
                <Route path="/problem/:id" element={<ViewProblemPage />} />
                <Route path="/problem/new" element={<CreateProblemPage />} />
                <Route path="/quiz" element={<QuizPage />} />
              </Routes>
            </div>
          </Theme>
        </div>
      </ToastProvider>
    </Router>
  )
}

export default App
