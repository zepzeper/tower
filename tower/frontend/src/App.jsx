import { Routes, Route } from 'react-router-dom';
import LanguageWrapper from './components/LanguageWrapper';
import Home from './pages/Home';
import Connections from './pages/Connections';
import NotFound from './pages/NotFound';
import Layout from './components/Layout';

const App = () => {
  return (
    <Routes>
      {/* Handle language detection and initial redirect */}
      <Route path="/" element={<LanguageWrapper />} />
      
      {/* All language-specific routes */}
      <Route path="/:lang" element={<Layout />}>
        <Route index element={<Home />} />
        <Route path="product" element={<Connections />} />
        <Route path="integrations" element={<Connections />} />
        <Route path="pricing" element={<Connections />} />
        <Route path="*" element={<NotFound />} />
      </Route>
    </Routes>
  );
};

export default App;
