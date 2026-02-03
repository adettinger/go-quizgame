import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css'
import Welcome from './pages/Welcome';
import NotFound from './pages/NotFound';

function App() {

  return (
    <Router>
      <div className='App'>
        {/* TODO: NavBar */}
        <span>Quizgame NavBar</span>
        <Routes>
          <Route path="/" element={<Welcome />} />
          <Route path="/*" element={<NotFound />} />
        </Routes>
      </div>
    </Router>
  )
}

export default App
