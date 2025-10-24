"use client";

import { apiClient } from "@/lib/api";
import Link from "next/link";
import { useEffect, useState } from "react";

interface FunctionItem {
  id: string;
  name: string;
  description: string;
  runtime: string;
  status: string;
  created_at: string;
}

export default function DashboardPage() {
  const [functions, setFunctions] = useState<FunctionItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    loadFunctions();
  }, []);

  const loadFunctions = async () => {
    try {
      const data = await apiClient.listFunctions();
      console.log("API Response:", data); // Debug log
      // Handle both array and object responses
      if (Array.isArray(data)) {
        setFunctions(data);
      } else if (data && typeof data === "object" && "functions" in data) {
        const response = data as { functions?: FunctionItem[] };
        setFunctions(response.functions || []);
      } else {
        setFunctions([]);
      }
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to load functions";
      setError(errorMessage);
      setFunctions([]);
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading functions...</div>;
  }

  return (
    <div>
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Functions</h1>
            <p className="mt-2 text-sm text-gray-600">
              Manage and deploy your cloud functions
            </p>
          </div>
          <Link
            href="/dashboard/functions/new"
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md font-medium"
          >
            + New Function
          </Link>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {functions.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <div className="text-6xl mb-4">⚡</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            No functions yet
          </h3>
          <p className="text-gray-600 mb-6">
            Get started by creating your first cloud function
          </p>
          <Link
            href="/dashboard/functions/new"
            className="inline-block bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-md font-medium"
          >
            Create Function
          </Link>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
          {functions.map((func) => (
            <Link
              key={func.id}
              href={`/dashboard/functions/${func.id}`}
              className="bg-white rounded-lg shadow hover:shadow-md transition-shadow p-6"
            >
              <div className="flex items-start justify-between mb-3">
                <h3 className="text-lg font-semibold text-gray-900">
                  {func.name}
                </h3>
                <span
                  className={`px-2 py-1 text-xs rounded-full ${
                    func.status === "active"
                      ? "bg-green-100 text-green-800"
                      : "bg-gray-100 text-gray-800"
                  }`}
                >
                  {func.status}
                </span>
              </div>
              <p className="text-sm text-gray-600 mb-4 line-clamp-2">
                {func.description || "No description"}
              </p>
              <div className="flex items-center justify-between text-xs text-gray-500">
                <span className="bg-gray-100 px-2 py-1 rounded">
                  {func.runtime}
                </span>
                <span>{new Date(func.created_at).toLocaleDateString()}</span>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
