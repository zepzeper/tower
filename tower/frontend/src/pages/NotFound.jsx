import { useTranslation } from 'react-i18next';
import Layout from '../components/Layout';

const NotFound = () => {
  const { t } = useTranslation('routing');
  
  return (
    <Layout>
      <div className="...">
        <h1>{t('404.title')}</h1>
        <p>{t('404.message')}</p>
      </div>
    </Layout>
  );
};

export default NotFound;
