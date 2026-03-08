'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface Assignment {
    ID: string;
    Title: string;
    Description: string;
    Language: string;
    DueDate: string;
}

export default function Assignments() {
    const [assignments, setAssignments] = useState<Assignment[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState('');
    const router = useRouter();

    useEffect(() => {
        const fetchAssignments = async () => {
            try {
                const token = localStorage.getItem('token');
                if (!token) {
                    router.push('/login');
                    return;
                }

                const res = await fetch('http://localhost:8080/api/assignments', {
                    headers: { Authorization: `Bearer ${token}` },
                });

                if (res.ok) {
                    const data = await res.json();
                    setAssignments(Array.isArray(data) ? data : []);
                } else if (res.status === 401) {
                    localStorage.removeItem('token');
                    router.push('/login');
                } else {
                    throw new Error('Failed to fetch assignments');
                }
            } catch (err) {
                setError('Unable to fetch assignments from the API.');
            } finally {
                setIsLoading(false);
            }
        };

        fetchAssignments();
    }, [router]);

    const handleLogout = () => {
        localStorage.removeItem('token');
        router.push('/');
    };

    return (
        <div className="min-h-screen bg-[#fdfdfc] dark:bg-[#111110] text-gray-900 dark:text-gray-100 font-sans">
            
            <header className="border-b border-gray-200 dark:border-gray-800 bg-white dark:bg-[#111110]">
                <div className="max-w-5xl mx-auto px-6 h-14 flex items-center justify-between">
                    <div className="flex items-center gap-4">
                        <Link href="/" className="font-semibold text-gray-900 dark:text-white">
                            PESU Dashboard
                        </Link>
                        <span className="text-gray-300 dark:text-gray-700">|</span>
                        <span className="text-sm text-gray-600 dark:text-gray-400 hidden sm:block">Assignments</span>
                    </div>
                    
                    <button 
                        onClick={handleLogout}
                        className="text-sm text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors"
                    >
                        Log out
                    </button>
                </div>
            </header>

            <main className="max-w-5xl mx-auto px-6 py-10">
                <div className="mb-8 border-b border-gray-200 dark:border-gray-800 pb-4">
                    <h1 className="text-3xl font-bold tracking-tight mb-2">My Assignments</h1>
                    <p className="text-gray-600 dark:text-gray-400 text-sm">
                        View and manage exactly what is required for your courses. Submissions should be done via the CLI.
                    </p>
                </div>
                
                {isLoading ? (
                    <div className="animate-pulse flex space-x-4">
                       <div className="flex-1 space-y-4 py-1">
                         <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-3/4"></div>
                         <div className="space-y-2">
                           <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded"></div>
                           <div className="h-4 bg-gray-200 dark:bg-gray-800 rounded w-5/6"></div>
                         </div>
                       </div>
                    </div>
                ) : error ? (
                    <div className="bg-red-50 dark:bg-red-900/10 border-l-4 border-red-500 p-4 rounded text-sm text-red-700 dark:text-red-400 mb-6">
                        <p className="font-bold">Error loading data</p>
                        <p>{error} Make sure the local Go API is running on port 8080.</p>
                        <button onClick={() => window.location.reload()} className="mt-2 text-red-600 dark:text-red-500 font-medium hover:underline">
                            Try again
                        </button>
                    </div>
                ) : assignments.length === 0 ? (
                    <div className="py-12 text-center border-2 border-dashed border-gray-200 dark:border-gray-800 rounded-lg">
                        <p className="text-gray-500 dark:text-gray-400 font-medium">No assignments currently posted.</p>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Check back later for updates from instructors.</p>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {assignments.map((a) => {
                            const isPastDue = new Date(a.DueDate) < new Date();
                            const formattedDate = new Date(a.DueDate).toLocaleDateString(undefined, { 
                                year: 'numeric', month: 'short', day: 'numeric' 
                            });
                            
                            return (
                                <div 
                                    key={a.ID} 
                                    className="p-5 bg-white dark:bg-[#1a1a1a] border border-gray-200 dark:border-gray-800 rounded-md shadow-sm hover:border-gray-300 dark:hover:border-gray-700 transition-colors"
                                >
                                    <div className="flex flex-col sm:flex-row sm:justify-between sm:items-start mb-2 gap-2">
                                        <h2 className="text-lg font-bold text-gray-900 dark:text-white">
                                            {a.Title}
                                        </h2>
                                        
                                        <div className="flex items-center gap-2 text-xs font-semibold">
                                            <span className="bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300 px-2 py-0.5 rounded">
                                                {a.Language}
                                            </span>
                                            <span className={`px-2 py-0.5 rounded ${isPastDue ? 'bg-red-100 text-red-800 dark:bg-red-900/40 dark:text-red-300' : 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300'}`}>
                                                {isPastDue ? 'Past Due' : 'Active'}
                                            </span>
                                        </div>
                                    </div>
                                    
                                    <p className="text-gray-600 dark:text-gray-400 text-sm mb-4 leading-relaxed line-clamp-2">
                                        {a.Description}
                                    </p>
                                    
                                    <div className="flex items-center text-xs text-gray-500 dark:text-gray-500 bg-gray-50 dark:bg-[#111110] p-2 rounded inline-flex border border-gray-100 dark:border-gray-800">
                                        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 mr-1.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                        </svg>
                                        Due <strong className="ml-1 font-medium text-gray-700 dark:text-gray-300">{formattedDate}</strong>
                                    </div>
                                </div>
                            );
                        })}
                    </div>
                )}
            </main>
        </div>
    );
}
