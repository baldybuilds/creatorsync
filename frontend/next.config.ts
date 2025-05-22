import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'images.unsplash.com'
      },
      {
        protocol: 'https',
        hostname: 'static-cdn.jtvnw.net'  // Twitch static content CDN
      },
      {
        protocol: 'https',
        hostname: 'vod-secure.twitch.tv'  // Twitch VOD CDN
      },
      {
        protocol: 'https',
        hostname: 'i.ytimg.com'          // YouTube thumbnails
      },
      {
        protocol: 'https',
        hostname: 'img.youtube.com'       // Alternative YouTube image domain
      }
    ],
  }
};

export default nextConfig;
