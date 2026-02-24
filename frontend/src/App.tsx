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
import { GameHostPage } from './pages/GameHost';
import { GamePlayerPage } from './pages/GamePlayerPage';
import { EditProblemPage } from './pages/EditProblemPage';
import Banner from './components/Banner/Banner';
import { LegalText } from './components/LegalText';

function App() {

  return (
    <Router>
      <ToastProvider>
        <div className='App'>
          <Theme>
            <NavBar />
            <div className="container mt-4">
              <Banner />

              <Routes>
                <Route path="/" element={<Welcome />} />
                <Route path="/*" element={<NotFound />} />
                <Route path="/problems" element={<ProblemsPage />} />
                <Route path="/problem/:id" element={<ViewProblemPage />} />
                <Route path="/problem/edit/:id" element={<EditProblemPage />} />
                <Route path="/problem/new" element={<CreateProblemPage />} />
                <Route path="/quiz" element={<QuizPage />} />
                <Route path="/game/host" element={<GameHostPage />} />
                <Route path="/game/player" element={<GamePlayerPage />} />
              </Routes>
            </div>
            <LegalText />
          </Theme>
        </div>
      </ToastProvider>
    </Router>
  )
}

export default App
