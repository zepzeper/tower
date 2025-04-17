import { Header } from './Header';
import { Footer } from './Footer';
import { Outlet } from 'react-router-dom';

const Layout = () => (
  <div className="min-h-screen flex flex-col">
    <Header />
    <main className="flex-grow">
      <Outlet />
    </main>
    <Footer />
  </div>
);

export default Layout;
