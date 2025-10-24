import Link from "next/link";

export default function Home() {
  return (
    <div className="min-h-screen bg-linear-to-br from-blue-50 via-indigo-50 to-purple-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <header className="flex justify-between items-center mb-20">
          <div className="text-3xl font-bold text-indigo-600">⚡ VoltRun</div>
          <div className="space-x-4">
            <Link
              href="/dashboard"
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md font-medium"
            >
              Go to Dashboard
            </Link>
          </div>
        </header>
        <div className="text-center max-w-4xl mx-auto">
          <h1 className="text-6xl font-extrabold text-gray-900 mb-6">
            Execute Functions at
            <span className="text-indigo-600"> Lightning Speed</span>
          </h1>
          <p className="text-xl text-gray-600 mb-10">
            Deploy and run your code in isolated Firecracker MicroVMs. Secure,
            scalable, and fast.
          </p>
          <div className="flex justify-center gap-4 mb-20">
            <Link
              href="/register"
              className="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-4 rounded-lg text-lg font-semibold"
            >
              Start Building
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
