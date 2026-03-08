import Link from "next/link";
import Image from "next/image";

export default function Home() {
  return (
    <div className="min-h-screen bg-[#fdfdfc] dark:bg-[#111110] text-gray-800 dark:text-gray-200 font-sans selection:bg-blue-100 dark:selection:bg-blue-900/40">
      
      <header className="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-[#111110]">
        <div className="max-w-5xl mx-auto px-6 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="bg-blue-600 text-white w-7 h-7 flex items-center justify-center font-bold text-sm rounded">
              P
            </div>
            <span className="font-semibold text-[15px] tracking-wide text-gray-900 dark:text-white">PESU CLI Dashboard</span>
          </div>
          <nav className="flex items-center gap-6 text-sm font-medium">
            <Link href="/login" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">
              Student Login
            </Link>
            <a href="https://github.com" target="_blank" rel="noreferrer" className="text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">
              Documentation
            </a>
          </nav>
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-20">
        <div className="mb-10">
          <h1 className="text-4xl sm:text-5xl font-bold tracking-tight text-gray-900 dark:text-white mb-6">
            A reliable way to manage and track your coding evaluations.
          </h1>
          <p className="text-lg text-gray-600 dark:text-gray-400 mb-8 max-w-2xl leading-relaxed">
            Welcome to the PESU CLI portal. This dashboard works alongside your local command-line tool, letting you safely track submissions and verify they executed properly in the sandboxed grading environment.
          </p>
          <div className="flex flex-wrap gap-4">
             <Link href="/login" className="bg-blue-600 hover:bg-blue-700 text-white px-5 py-2.5 rounded-md font-medium transition-colors border border-blue-700 shadow-sm text-sm">
               Log in to dashboard
             </Link>
             <code className="bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-300 px-5 py-2.5 rounded-md text-sm border border-gray-200 dark:border-gray-700 font-mono flex items-center">
               $ pesu login
             </code>
          </div>
        </div>

        <hr className="border-gray-200 dark:border-gray-800 my-16" />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
          
          <div>
            <div className="w-10 h-10 bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400 rounded flex items-center justify-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Language Support</h3>
            <p className="text-gray-600 dark:text-gray-400 leading-relaxed text-sm">
              The worker service handles both compiled and interpreted languages. Currently, C and Python are supported out of the box with strict memory and CPU limits.
            </p>
          </div>

          <div>
            <div className="w-10 h-10 bg-green-50 dark:bg-green-900/20 text-green-600 dark:text-green-400 rounded flex items-center justify-center mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-white mb-2">Task Queue</h3>
            <p className="text-gray-600 dark:text-gray-400 leading-relaxed text-sm">
              Submissions are offloaded to a Redis queue. The dashboard gives you insight into whether your submission is pending, executing, or fully graded.
            </p>
          </div>

        </div>
      </main>

      <footer className="border-t border-gray-200 dark:border-gray-800 bg-white dark:bg-[#111110] py-8 mt-auto">
        <div className="max-w-5xl mx-auto px-6 text-sm text-gray-500 flex flex-col sm:flex-row justify-between items-center gap-4">
          <p>© {new Date().getFullYear()} PESU. All rights reserved.</p>
          <div className="flex items-center gap-4">
             <Image src="/logoPesu.png" alt="PESU Logo" width={80} height={35} className="dark:invert opacity-60 grayscale hover:grayscale-0 transition-all" />
          </div>
        </div>
      </footer>
    </div>
  );
}
