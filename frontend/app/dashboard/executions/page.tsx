"use client";

import { apiClient } from "@/lib/api";
import Link from "next/link";
import { useEffect, useState } from "react";

interface Execution {
  id: string;
  function_id: string;
  status: string;
  input: string;
  output: string;
  error?: string;
  logs: string;
  duration_ms: number;
  created_at: string;
  function?: {
    name: string;
    runtime: string;
  };
}

export default function ExecutionsPage() {
  const [executions, setExecutions] = useState<Execution[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [filter, setFilter] = useState<string>("all");
  const [selectedExecution, setSelectedExecution] = useState<Execution | null>(
    null
  );

  useEffect(() => {
    loadExecutions();
  }, []);

  const loadExecutions = async () => {
    try {
      const data = await apiClient.listExecutions();
      setExecutions(Array.isArray(data) ? data : []);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to load executions";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const filteredExecutions = executions.filter((ex) => {
    if (filter === "all") return true;
    return ex.status === filter;
  });

  if (isLoading) {
    return <div className="text-center py-8">Loading executions...</div>;
  }

  return (
    <div>
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Executions</h1>
        <p className="mt-2 text-sm text-gray-600">
          View and monitor function execution history
        </p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* Filters */}
      <div className="mb-6 flex gap-2">
        {["all", "pending", "running", "success", "failed"].map((status) => (
          <button
            key={status}
            onClick={() => setFilter(status)}
            className={`px-4 py-2 rounded-md font-medium ${
              filter === status
                ? "bg-indigo-600 text-white"
                : "bg-gray-100 text-gray-700 hover:bg-gray-200"
            }`}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </button>
        ))}
      </div>

      {filteredExecutions.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <div className="text-6xl mb-4">📊</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            No executions yet
          </h3>
          <p className="text-gray-600 mb-6">
            Execute a function to see its execution history here
          </p>
          <Link
            href="/dashboard"
            className="inline-block bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-md font-medium"
          >
            Go to Functions
          </Link>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Function
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Status
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Duration
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Created
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {filteredExecutions.map((execution) => (
                <tr
                  key={execution.id}
                  className="hover:bg-gray-50 cursor-pointer"
                  onClick={() => setSelectedExecution(execution)}
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {execution.function?.name || "Unknown"}
                    </div>
                    <div className="text-sm text-gray-500">
                      {execution.function?.runtime || ""}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span
                      className={`px-2 py-1 text-xs rounded-full ${
                        execution.status === "success"
                          ? "bg-green-100 text-green-800"
                          : execution.status === "failed"
                          ? "bg-red-100 text-red-800"
                          : execution.status === "running"
                          ? "bg-blue-100 text-blue-800"
                          : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {execution.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {execution.duration_ms}ms
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(execution.created_at).toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        setSelectedExecution(execution);
                      }}
                      className="text-indigo-600 hover:text-indigo-900"
                    >
                      View Details
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Execution Detail Modal */}
      {selectedExecution && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
          onClick={() => setSelectedExecution(null)}
        >
          <div
            className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="p-6">
              <div className="flex justify-between items-start mb-6">
                <div>
                  <h2 className="text-2xl font-bold text-gray-900">
                    Execution Details
                  </h2>
                  <p className="text-sm text-gray-500 mt-1">
                    ID: {selectedExecution.id}
                  </p>
                </div>
                <button
                  onClick={() => setSelectedExecution(null)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <svg
                    className="w-6 h-6"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                </button>
              </div>

              <div className="space-y-4">
                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    Status
                  </h3>
                  <span
                    className={`px-3 py-1 text-sm rounded-full ${
                      selectedExecution.status === "success"
                        ? "bg-green-100 text-green-800"
                        : selectedExecution.status === "failed"
                        ? "bg-red-100 text-red-800"
                        : selectedExecution.status === "running"
                        ? "bg-blue-100 text-blue-800"
                        : "bg-gray-100 text-gray-800"
                    }`}
                  >
                    {selectedExecution.status}
                  </span>
                </div>

                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    Input
                  </h3>
                  <pre className="bg-gray-50 p-3 rounded text-sm overflow-x-auto text-gray-900">
                    {selectedExecution.input || "{}"}
                  </pre>
                </div>

                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    Output
                  </h3>
                  <pre className="bg-gray-50 p-3 rounded text-sm overflow-x-auto text-gray-900">
                    {selectedExecution.output || "No output"}
                  </pre>
                </div>

                {selectedExecution.error && (
                  <div>
                    <h3 className="text-sm font-semibold text-gray-700 mb-2">
                      Error
                    </h3>
                    <pre className="bg-red-50 p-3 rounded text-sm overflow-x-auto text-red-900">
                      {selectedExecution.error}
                    </pre>
                  </div>
                )}

                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    Logs
                  </h3>
                  <pre className="bg-gray-900 text-gray-100 p-3 rounded text-sm overflow-x-auto max-h-64">
                    {selectedExecution.logs || "No logs"}
                  </pre>
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <span className="text-gray-500">Duration:</span>
                    <span className="ml-2 font-medium text-gray-900">
                      {selectedExecution.duration_ms}ms
                    </span>
                  </div>
                  <div>
                    <span className="text-gray-500">Created:</span>
                    <span className="ml-2 font-medium text-gray-900">
                      {new Date(selectedExecution.created_at).toLocaleString()}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
