import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';

export const CaseStudiesSection = () => {
  const { t, i18n } = useTranslation('pages');
  const lang = i18n.language;
  
  // Define case studies
  const caseStudies = [
    {
      id: 'company-a',
      image: '/case-studies/company-a.webp',
      title: t('homepage.caseStudies.items.companyA.title'),
      description: t('homepage.caseStudies.items.companyA.description'),
      slug: 'company-a'
    },
    {
      id: 'company-b',
      image: '/case-studies/company-b.webp',
      title: t('homepage.caseStudies.items.companyB.title'),
      description: t('homepage.caseStudies.items.companyB.description'),
      slug: 'company-b'
    },
    {
      id: 'company-c',
      image: '/case-studies/company-c.webp',
      title: t('homepage.caseStudies.items.companyC.title'),
      description: t('homepage.caseStudies.items.companyC.description'),
      slug: 'company-c'
    }
  ];
  
  return (
    <div className="bg-white pb-32">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
            {t('homepage.caseStudies.title')}
          </h2>
          <p className="mt-2 text-lg leading-8 text-gray-600">
            {t('homepage.caseStudies.subtitle')}
          </p>
        </div>
        
        <div className="mx-auto mt-16 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-20 lg:mx-0 lg:max-w-none lg:grid-cols-3">
          {caseStudies.map((study) => (
            <article key={study.id} className="flex flex-col items-start justify-between">
              <div className="relative w-full">
                <img 
                  src={study.image} 
                  alt={study.title}
                  className="aspect-[16/9] w-full rounded-2xl bg-gray-100 object-cover sm:aspect-[2/1] lg:aspect-[3/2]"
                />
                <div className="absolute inset-0 rounded-2xl ring-1 ring-inset ring-gray-900/10"></div>
              </div>
              
              <div className="max-w-xl">
                <div className="group relative">
                  <h3 className="mt-3 text-lg font-semibold leading-6 text-gray-900 group-hover:text-gray-600">
                    <Link to={`/${lang}/case-studies/${study.slug}`} className="">
                      <span className="absolute inset-0"></span>
                      {study.title}
                    </Link>
                  </h3>
                  <p className="mt-5 line-clamp-3 text-sm leading-6 text-gray-600">
                    {study.description}
                  </p>
                </div>
              </div>
            </article>
          ))}
        </div>
      </div>
    </div>
  );
};
