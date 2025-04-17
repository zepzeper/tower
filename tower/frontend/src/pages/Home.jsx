import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { CaseStudiesSection } from './Home/CaseStudiesSection';
import { FeatureSection } from './Home/FeatureSection';
import { SecondaryFeatureSection } from './Home/SecondaryFeatureSection';
import { 
  Zap, 
  ArrowRight, 
  ChevronLeft,
  ChevronRight
} from 'lucide-react';

const HomePage = () => {
  const { t, i18n } = useTranslation('pages');
  const lang = i18n.language;
  
  // Customer logos
  const customerLogos = [
    { id: 1, name: 'Company A', logo: '/logos/company-a.webp' },
    { id: 2, name: 'Company B', logo: '/logos/company-b.webp' },
    { id: 3, name: 'Company C', logo: '/logos/company-c.webp' },
    { id: 4, name: 'Company D', logo: '/logos/company-d.webp' },
    { id: 5, name: 'Company E', logo: '/logos/company-e.webp' }
  ];
  
  // Logos slide controller
  const [currentSlide, setCurrentSlide] = useState(0);
  const totalSlides = 4;
  
  const nextSlide = () => {
    setCurrentSlide((prev) => (prev + 1) % totalSlides);
  };
  
  const prevSlide = () => {
    setCurrentSlide((prev) => (prev - 1 + totalSlides) % totalSlides);
  };
  
  const goToSlide = (slideIndex) => {
    setCurrentSlide(slideIndex);
  };

  return (
    <main>
      {/* Hero Section */}
      <div className="pb-16 pt-20 text-center lg:pt-32 mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h1 className="mx-auto max-w-5xl font-bold text-5xl tracking-tight text-slate-900 sm:text-7xl">
          {t('homepage.hero.titleStart')} {' '}
          <span className="relative whitespace-nowrap text-green-600">
            <svg aria-hidden="true" viewBox="0 0 418 42" className="absolute left-0 top-2/3 h-[0.58em] w-full fill-green-300/70" preserveAspectRatio="none">
              <path d="M203.371.916c-26.013-2.078-76.686 1.963-124.73 9.946L67.3 12.749C35.421 18.062 18.2 21.766 6.004 25.934 1.244 27.561.828 27.778.874 28.61c.07 1.214.828 1.121 9.595-1.176 9.072-2.377 17.15-3.92 39.246-7.496C123.565 7.986 157.869 4.492 195.942 5.046c7.461.108 19.25 1.696 19.17 2.582-.107 1.183-7.874 4.31-25.75 10.366-21.992 7.45-35.43 12.534-36.701 13.884-2.173 2.308-.202 4.407 4.442 4.734 2.654.187 3.263.157 15.593-.78 35.401-2.686 57.944-3.488 88.365-3.143 46.327.526 75.721 2.23 130.788 7.584 19.787 1.924 20.814 1.98 24.557 1.332l.066-.011c1.201-.203 1.53-1.825.399-2.335-2.911-1.31-4.893-1.604-22.048-3.261-57.509-5.556-87.871-7.36-132.059-7.842-23.239-.254-33.617-.116-50.627.674-11.629.54-42.371 2.494-46.696 2.967-2.359.259 8.133-3.625 26.504-9.81 23.239-7.825 27.934-10.149 28.304-14.005.417-4.348-3.529-6-16.878-7.066Z"></path>
            </svg>
            <span className="relative">{t('homepage.hero.titleHighlight')}</span>
          </span>{' '}
          {t('homepage.hero.titleEnd')}
        </h1>
        
        <p className="mx-auto mt-6 max-w-2xl text-lg tracking-tight text-slate-700">
          {t('homepage.hero.subtitle')}
        </p>
        
        <div className="mt-10 flex justify-center gap-x-6">
          <Link 
            to={`/${lang}/register`} 
            className="bg-green-600 text-white hover:bg-green-700 focus-visible:outline-green-900 group inline-flex items-center justify-center rounded-full py-2 px-4 text-sm font-semibold focus:outline-none focus-visible:outline-2 focus-visible:outline-offset-2"
          >
            {t('homepage.hero.primaryCta')}
          </Link>
          
          <Link 
            to={`/${lang}/login`} 
            className="group inline-flex ring-1 items-center justify-center rounded-full py-2 px-4 text-sm focus:outline-none ring-slate-200 text-slate-700 hover:text-slate-900 hover:ring-slate-300 active:bg-slate-100 active:text-slate-600 focus-visible:outline-green-600 focus-visible:ring-slate-300"
          >
            {t('homepage.hero.secondaryCta')}
          </Link>
        </div>
        
        <div className="mt-10 -mb-10 flex justify-center">
          <Link to={`/${lang}/partners`} className="mx-0 sm:mx-auto lg:mx-0 h-6 sm:h-6 md:h-6 lg:h-6">
            <img src="/img/partner-badge.webp" height="32" alt="Partner Badge" className="h-8" />
          </Link>
        </div>
        
        {/* Logos Section */}
        <div className="bg-white mt-36 lg:mt-44">
          <div className="mx-auto max-w-7xl px-6 lg:px-8">
            <h2 className="text-center text-lg font-bold leading-8 text-gray-900">
              {t('homepage.logos.title')}
            </h2>
            
            <div className="relative mx-auto mt-10 overflow-visible">
              <div className="overflow-hidden px-8 md:px-12">
                <div 
                  className="flex transition-transform duration-500 ease-in-out" 
                  style={{ transform: `translateX(-${currentSlide * 100}%)` }}
                >
                  {/* First slide */}
                  <div className="w-full flex-shrink-0">
                    <div className="mx-auto grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 items-center gap-x-4 sm:gap-x-6 md:gap-x-8 gap-y-6 md:gap-y-10 w-full">
                      {customerLogos.map((logo) => (
                        <div key={logo.id} className="col-span-1 flex justify-center items-center px-2 sm:px-4 min-h-[60px] sm:min-h-[80px]">
                          <div className="h-12 w-full bg-gray-100 rounded-md flex items-center justify-center text-gray-700 font-medium">
                            {logo.name}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                  
                  {/* Additional slides would be here */}
                </div>
              </div>
              
              {/* Navigation buttons */}
              <button 
                className="absolute left-0 top-1/2 -translate-y-1/2 bg-white/90 p-1.5 sm:p-2 rounded-full shadow-md z-20 hover:bg-white"
                aria-label="Previous slide"
                onClick={prevSlide}
              >
                <ChevronLeft className="h-4 w-4 sm:h-5 sm:w-5 text-gray-700" />
              </button>
              
              <button 
                className="absolute right-0 top-1/2 -translate-y-1/2 bg-white/90 p-1.5 sm:p-2 rounded-full shadow-md z-20 hover:bg-white"
                aria-label="Next slide"
                onClick={nextSlide}
              >
                <ChevronRight className="h-4 w-4 sm:h-5 sm:w-5 text-gray-700" />
              </button>
              
              {/* Pagination dots */}
              <div className="flex justify-center mt-4">
                {Array.from({ length: totalSlides }).map((_, index) => (
                  <button
                    key={index}
                    className={`mx-1 h-2 w-2 rounded-full ${
                      currentSlide === index ? 'bg-green-600' : 'bg-gray-300'
                    }`}
                    aria-label={`Go to slide ${index + 1}`}
                    onClick={() => goToSlide(index)}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <FeatureSection />
      <SecondaryFeatureSection />
      <CaseStudiesSection />
      
    </main>
  );
};

export default HomePage;
