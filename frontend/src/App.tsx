import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css'
import Welcome from './pages/Welcome';
import NotFound from './pages/NotFound';
import { ProblemsPage } from './pages/Problems';
import NavBar from './components/NavBar/NavBar';

function App() {

  return (
    <Router>
      <div className='App'>
        <NavBar />
        <div className="container mt-4">
          <div className='App-header'>
            <h1>Quizgame</h1>
          </div>

          <Routes>
            <Route path="/" element={<Welcome />} />
            <Route path="/*" element={<NotFound />} />
            <Route path="/problems" element={<ProblemsPage />} />
          </Routes>
        </div>
      </div>
    </Router>
  )
}

export default App
