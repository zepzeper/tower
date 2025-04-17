import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  LineChart, 
  BarChart3, 
  Link as LinkIcon, 
  Settings,
  Mail,
  FileText,
  BarChart
} from 'lucide-react';

export const FeatureSection = () => {
  const { t } = useTranslation('pages');
  const [activeTab, setActiveTab] = useState('api');
  
  // Define the tabs and their content
  const tabs = [
    {
      id: 'api',
      title: t('homepage.features.tabs.api.title'),
      description: t('homepage.features.tabs.api.description'),
      image: '/landing/screenshot-api-connections.webp',
      imageAlt: t('homepage.features.tabs.api.imageAlt')
    },
    {
      id: 'dashboard',
      title: t('homepage.features.tabs.dashboard.title'),
      description: t('homepage.features.tabs.dashboard.description'),
      image: '/landing/screenshot-dashboard.webp',
      imageAlt: t('homepage.features.tabs.dashboard.imageAlt')
    },
    {
      id: 'automation',
      title: t('homepage.features.tabs.automation.title'),
      description: t('homepage.features.tabs.automation.description'),
      image: '/landing/screenshot-automation.webp',
      imageAlt: t('homepage.features.tabs.automation.imageAlt')
    },
    {
      id: 'analytics',
      title: t('homepage.features.tabs.analytics.title'),
      description: t('homepage.features.tabs.analytics.description'),
      image: '/landing/screenshot-analytics.webp',
      imageAlt: t('homepage.features.tabs.analytics.imageAlt')
    }
  ];
  
  return (
    <section id="features" aria-label="Features for connecting APIs" className="relative overflow-hidden bg-green-600">
      <img 
        src="/landing/background-features.webp" 
        width="2245" 
        height="1636" 
        alt="Features background" 
        className="absolute left-1/2 top-1/2 max-w-none translate-x-[-44%] translate-y-[-42%]" 
      />
      
      <div className="pb-28 pt-20 sm:py-32">
        <div className="relative w-full mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="max-w-2xl md:mx-auto md:text-center xl:max-w-none">
            <h2 className="font-semibold text-3xl tracking-tight text-white sm:text-4xl md:text-5xl">
              {t('homepage.features.title')}
            </h2>
            <p className="mt-6 text-lg tracking-tight text-green-100">
              {t('homepage.features.subtitle')}
            </p>
          </div>
          
          <div className="mt-16 grid grid-cols-1 items-center gap-y-2 pt-10 sm:gap-y-6 md:mt-20 lg:grid-cols-12 lg:pt-0">
            {/* Tabs navigation */}
            <div className="-mx-4 flex overflow-x-auto pb-4 sm:mx-0 sm:overflow-visible sm:pb-0 lg:col-span-5">
              <div role="tablist" aria-orientation="vertical" className="relative z-10 flex gap-x-4 whitespace-nowrap px-4 sm:mx-auto sm:px-0 lg:mx-0 lg:block lg:gap-x-0 lg:gap-y-1 lg:whitespace-normal">
                {tabs.map((tab) => (
                  <div 
                    key={tab.id}
                    className={`${
                      activeTab === tab.id 
                        ? 'bg-white lg:bg-white/10 lg:ring-1 lg:ring-inset lg:ring-white/10' 
                        : 'hover:bg-white/10 lg:hover:bg-white/5'
                    } group relative rounded-full px-4 py-1 lg:rounded-l-xl lg:rounded-r-none lg:p-6`}
                  >
                    <h3>
                      <button
                        role="tab"
                        type="button" 
                        aria-selected={activeTab === tab.id}
                        className={`font-semibold text-lg focus:outline-none ${
                          activeTab === tab.id 
                            ? 'text-green-600 lg:text-white' 
                            : 'text-green-100 hover:text-white lg:text-white'
                        }`}
                        onClick={() => setActiveTab(tab.id)}
                      >
                        <span className="absolute inset-0 rounded-full lg:rounded-l-xl lg:rounded-r-none"></span>
                        {tab.title}
                      </button>
                    </h3>
                    <p className={`mt-2 hidden text-sm lg:block ${
                      activeTab === tab.id 
                        ? 'text-white' 
                        : 'text-green-100 group-hover:text-white'
                    }`}>
                      {tab.description}
                    </p>
                  </div>
                ))}
              </div>
            </div>
            
            {/* Tab content */}
            <div className="lg:col-span-7">
              {tabs.map((tab) => (
                <div 
                  key={tab.id}
                  role="tabpanel" 
                  aria-hidden={activeTab !== tab.id}
                  style={{ display: activeTab === tab.id ? 'block' : 'none' }}
                >
                  {activeTab === tab.id && (
                    <>
                      <div className="relative sm:px-6 lg:hidden">
                        <div className="absolute -inset-x-4 bottom-[-4.25rem] top-[-6.5rem] bg-white/10 ring-1 ring-inset ring-white/10 sm:inset-x-0 sm:rounded-t-xl"></div>
                        <p className="relative mx-auto max-w-2xl text-base text-white sm:text-center">
                          {tab.description}
                        </p>
                      </div>
                      <div className="mt-10 w-[45rem] overflow-hidden rounded-xl bg-slate-50 shadow-xl shadow-green-900/20 sm:w-auto lg:mt-0 lg:w-[67.8125rem]">
                        <img 
                          className="w-full" 
                          src={tab.image}
                          alt={tab.imageAlt}
                          sizes="(min-width: 1024px) 67.8125rem, (min-width: 640px) 100vw, 45rem"
                        />
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};
