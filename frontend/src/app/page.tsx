"use client";

import { usePostHogUser } from '../hooks/usePostHogUser';
import { Hero } from '../components/Hero';
import { Features } from '../components/Features';
import { About } from '../components/About';
import { Testimonials } from '../components/Testimonials';
import { Navbar } from '../components/ui/navbar';
import { Footer } from '../components/ui/footer';

const LandingPage = () => {
  usePostHogUser();

  return (
    <div className="min-h-screen bg-[var(--bg-primary)] text-[var(--text-primary)] overflow-x-hidden relative">
      {/* Unified Background Elements */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        {/* Subtle gradient overlay */}
        <div className="absolute inset-0 bg-gradient-to-br from-brand-500/5 via-transparent to-accent-500/10" />

        {/* High-resolution gradient orbs with improved blending */}
        <div
          className="absolute opacity-30 animate-float"
          style={{
            top: '15%',
            left: '20%',
            width: '60rem',
            height: '60rem',
            background: 'radial-gradient(circle, rgba(99,102,241,0.2) 0%, rgba(99,102,241,0) 70%)',
            filter: 'blur(60px)',
            animationDelay: '0s',
            transform: 'translate3d(0, 0, 0)'
          }}
        />
        <div
          className="absolute opacity-30 animate-float"
          style={{
            bottom: '10%',
            right: '15%',
            width: '50rem',
            height: '50rem',
            background: 'radial-gradient(circle, rgba(236,72,153,0.15) 0%, rgba(236,72,153,0) 70%)',
            filter: 'blur(60px)',
            animationDelay: '2s',
            transform: 'translate3d(0, 0, 0)'
          }}
        />
        <div
          className="absolute opacity-30 animate-float"
          style={{
            top: '60%',
            left: '30%',
            width: '45rem',
            height: '45rem',
            background: 'radial-gradient(circle, rgba(139,92,246,0.15) 0%, rgba(139,92,246,0) 70%)',
            filter: 'blur(60px)',
            animationDelay: '4s',
            transform: 'translate3d(0, 0, 0)'
          }}
        />
        <div
          className="absolute opacity-30 animate-float"
          style={{
            bottom: '40%',
            right: '25%',
            width: '55rem',
            height: '55rem',
            background: 'radial-gradient(circle, rgba(14,165,233,0.15) 0%, rgba(14,165,233,0) 70%)',
            filter: 'blur(60px)',
            animationDelay: '6s',
            transform: 'translate3d(0, 0, 0)'
          }}
        />

        {/* Additional subtle noise texture for depth */}
        <div
          className="absolute inset-0 opacity-[0.03] mix-blend-overlay"
          style={{
            backgroundImage: 'url("data:image/svg+xml,%3Csvg viewBox=\'0 0 200 200\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cfilter id=\'noiseFilter\'%3E%3CfeTurbulence type=\'fractalNoise\' baseFrequency=\'0.65\' numOctaves=\'3\' stitchTiles=\'stitch\'/%3E%3C/filter%3E%3Crect width=\'100%\' height=\'100%\' filter=\'url(%23noiseFilter)\' /%3E%3C/svg%3E")',
            backgroundSize: '200px 200px',
          }}
        />
      </div>

      {/* Navigation */}
      <Navbar />

      {/* Main Content */}
      <div className="relative z-10">
        <Hero />

        {/* Features Section */}
        <div id="features">
          <Features />
        </div>

        {/* About Section */}
        <div id="about">
          <About />
        </div>

        {/* Testimonials Section */}
        {/* <div id="testimonials">
          <Testimonials />
        </div> */}

        {/* Contact Section - Using id for navigation */}
        <div id="contact">
          {/* Footer also serves as contact section */}
          <Footer />
        </div>
      </div>
    </div>
  );
}

export default LandingPage;
