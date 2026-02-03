import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css'
import Welcome from './pages/Welcome';
import NotFound from './pages/NotFound';
import { ProblemsPage } from './pages/Problems';

function App() {

  return (
    <Router>
      <div className='App'>
        {/* TODO: NavBar */}
        <span>Quizgame NavBar</span>
        <Routes>
          <Route path="/" element={<Welcome />} />
          <Route path="/*" element={<NotFound />} />
          <Route path="/problems" element={<ProblemsPage />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App
