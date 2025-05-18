import { Inter } from 'next/font/google';
import '../styles/globals.css';
import { ClerkProvider } from "@clerk/nextjs";
import { dark } from '@clerk/themes';
import { ThemeProvider } from '@/components/theme-provider';
import { Toaster } from '@/components/ui/toaster';
import { PostHogProvider } from '@/app/providers';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'CreatorSync - Content Management for Creators',
  description: 'Manage and optimize your content creation workflow',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <PostHogProvider>
      <ClerkProvider
        appearance={{
          baseTheme: dark,
          variables: {
            colorPrimary: "#6366f1",
          },
        }}
        publishableKey={process.env.NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY}
      >
        <html lang="en" suppressHydrationWarning>
          <body className={`${inter.className} min-h-screen bg-background`}>
            <ThemeProvider
              attribute="class"
              defaultTheme="system"
              enableSystem
              disableTransitionOnChange
            >
              {children}
              <Toaster />
            </ThemeProvider>
          </body>
        </html>
      </ClerkProvider>
    </PostHogProvider>
  );
}