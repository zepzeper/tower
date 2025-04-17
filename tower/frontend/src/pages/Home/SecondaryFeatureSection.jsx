import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { 
  Mail, 
  RefreshCw, 
  BarChart,
  Webhook,
  ArrowDownToLine
} from 'lucide-react';

export const SecondaryFeatureSection = () => {
  const { t } = useTranslation('pages');
  const [activeFeature, setActiveFeature] = useState('email');
  
  // Define the secondary features
  const features = [
    {
      id: 'email',
      icon: <Mail className="h-6 w-6 text-white" />,
      title: t('homepage.secondaryFeatures.tabs.email.title'),
      headline: t('homepage.secondaryFeatures.tabs.email.headline'),
      description: t('homepage.secondaryFeatures.tabs.email.description'),
      image: '/landing/feature-email-automation.webp',
      imageAlt: t('homepage.secondaryFeatures.tabs.email.imageAlt')
    },
    {
      id: 'sync',
      icon: <RefreshCw className="h-6 w-6 text-white" />,
      title: t('homepage.secondaryFeatures.tabs.sync.title'),
      headline: t('homepage.secondaryFeatures.tabs.sync.headline'),
      description: t('homepage.secondaryFeatures.tabs.sync.description'),
      image: '/landing/feature-data-sync.webp',
      imageAlt: t('homepage.secondaryFeatures.tabs.sync.imageAlt')
    },
    {
      id: 'data',
      icon: <BarChart className="h-6 w-6 text-white" />,
      title: t('homepage.secondaryFeatures.tabs.data.title'),
      headline: t('secondaryFeatures.tabs.data.headline'),
      description: t('homepage.secondaryFeatures.tabs.data.description'),
      image: '/landing/feature-data-collection.webp',
      imageAlt: t('homepage.secondaryFeatures.tabs.data.imageAlt')
    }
  ];
  
  return (
    <section id="secondary-features" aria-label="Features for simplifying business tasks" className="pb-14 pt-20 sm:pb-20 sm:pt-32 lg:pb-32">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="mx-auto max-w-2xl md:text-center">
          <h2 className="font-semibold text-3xl tracking-tight text-slate-900 sm:text-4xl">
            {t('homepage.secondaryFeatures.title')}
          </h2>
          <p className="mt-4 text-lg tracking-tight text-slate-700">
            {t('homepage.secondaryFeatures.subtitle')}
          </p>
        </div>
        
        {/* Mobile view (stacked) */}
        <div className="-mx-4 mt-20 flex flex-col gap-y-10 overflow-hidden px-4 sm:-mx-6 sm:px-6 lg:hidden">
          {features.map((feature) => (
            <div key={feature.id}>
              <div className="mx-auto max-w-2xl">
                <div className="h-9 w-9 rounded-lg flex items-center justify-center bg-green-500">
                  {feature.icon}
                </div>
                <h3 className="mt-6 text-sm font-medium text-green-600">
                  {feature.title}
                </h3>
                <p className="mt-2 font-semibold text-xl text-slate-900">
                  {feature.headline}
                </p>
                <p className="mt-4 text-sm text-slate-600">
                  {feature.description}
                </p>
              </div>
              <div className="relative mt-10 pb-10">
                <div className="absolute -inset-x-4 bottom-0 top-8 bg-slate-200 sm:-inset-x-6"></div>
                <div className="relative mx-auto w-[52.75rem] overflow-hidden rounded-xl bg-white shadow-lg shadow-slate-900/5 ring-1 ring-slate-500/10">
                  <img 
                    className="w-full" 
                    src={feature.image} 
                    alt={feature.imageAlt} 
                    sizes="52.75rem" 
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
        
        {/* Desktop view (tabs) */}
        <div className="hidden lg:mt-20 lg:block">
          <div role="tablist" aria-orientation="horizontal" className="grid grid-cols-3 gap-x-8">
            {features.map((feature) => (
              <div 
                key={feature.id} 
                className={`${activeFeature === feature.id ? '' : 'opacity-75 hover:opacity-100'} relative`}
                onClick={() => setActiveFeature(feature.id)}
              >
                <div className={`h-9 w-9 rounded-lg flex items-center justify-center ${
                  activeFeature === feature.id ? 'bg-green-500' : 'bg-slate-500'
                }`}>
                  {feature.icon}
                </div>
                <h3 className={`mt-6 text-sm font-medium ${
                  activeFeature === feature.id ? 'text-green-600' : 'text-slate-600'
                }`}>
                  <button 
                    role="tab" 
                    type="button" 
                    aria-selected={activeFeature === feature.id}
                    className="focus:outline-none"
                  >
                    <span className="absolute inset-0"></span>
                    {feature.title}
                  </button>
                </h3>
                <p className="mt-2 font-semibold text-xl text-slate-900">
                  {feature.headline}
                </p>
                <p className="mt-4 text-sm text-slate-600">
                  {feature.description}
                </p>
              </div>
            ))}
          </div>
          
          {/* Feature image container */}
          <div className="relative mt-20 overflow-hidden rounded-3xl bg-slate-200 px-14 py-16 xl:px-16">
            <div className="-mx-5 flex">
              {features.map((feature) => (
                <div 
                  key={feature.id}
                  role="tabpanel"
                  aria-hidden={activeFeature !== feature.id}
                  className={`px-5 transition duration-500 ease-in-out focus:outline-none ${
                    activeFeature === feature.id ? '' : 'opacity-60'
                  }`}
                  style={{ transform: 'translateX(0%)' }}
                >
                  <div className="w-[52.75rem] overflow-hidden rounded-xl bg-white shadow-lg shadow-slate-900/5 ring-1 ring-slate-500/10">
                    <img 
                      src={feature.image} 
                      alt={feature.imageAlt} 
                      className="w-full" 
                    />
                  </div>
                </div>
              ))}
            </div>
            <div className="pointer-events-none absolute inset-0 rounded-3xl ring-1 ring-inset ring-slate-900/10"></div>
          </div>
        </div>
      </div>
    </section>
  );
};
