"use client";

import { apiClient } from "@/lib/api";
import { useRouter } from "next/navigation";
import { use, useEffect, useState } from "react";

interface FunctionData {
  id: string;
  name: string;
  description: string;
  runtime: string;
  code: string;
  entry_point: string;
  memory_mb: number;
  timeout_sec: number;
  status: string;
  created_at: string;
  updated_at: string;
}

interface Execution {
  id: string;
  status: string;
  input: string;
  output: string;
  error?: string;
  logs: string;
  duration_ms: number;
  created_at: string;
}

export default function FunctionDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const router = useRouter();
  const [functionData, setFunctionData] = useState<FunctionData | null>(null);
  const [executions, setExecutions] = useState<Execution[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isExecuting, setIsExecuting] = useState(false);
  const [error, setError] = useState("");
  const [executionInput, setExecutionInput] = useState("{}");
  const [executionResult, setExecutionResult] = useState<any>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [editedCode, setEditedCode] = useState("");

  useEffect(() => {
    loadFunction();
    loadExecutions();
  }, [id]);

  const loadFunction = async () => {
    try {
      const data = await apiClient.getFunction(id);
      setFunctionData(data);
      setEditedCode(data.code);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to load function";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const loadExecutions = async () => {
    try {
      const data = await apiClient.listExecutions();
      // Filter executions for this function
      const filtered = Array.isArray(data)
        ? data.filter((ex: any) => ex.function_id === id)
        : [];
      setExecutions(filtered);
    } catch (err) {
      console.error("Failed to load executions:", err);
    }
  };

  const handleExecute = async () => {
    setIsExecuting(true);
    setError("");
    setExecutionResult(null);

    try {
      // Parse input JSON
      let inputData = {};
      try {
        inputData = JSON.parse(executionInput);
      } catch (e) {
        throw new Error("Invalid JSON input");
      }

      const result = await apiClient.executeFunction(id, inputData);
      setExecutionResult(result);

      // Reload executions after a short delay
      setTimeout(() => {
        loadExecutions();
      }, 1000);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to execute function";
      setError(errorMessage);
    } finally {
      setIsExecuting(false);
    }
  };

  const handleSave = async () => {
    if (!functionData) return;

    try {
      await apiClient.updateFunction(id, {
        code: editedCode,
      });
      setFunctionData({ ...functionData, code: editedCode });
      setIsEditing(false);
      setError("");
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to update function";
      setError(errorMessage);
    }
  };

  const handleDelete = async () => {
    if (!confirm("Are you sure you want to delete this function?")) return;

    try {
      await apiClient.deleteFunction(id);
      router.push("/dashboard");
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to delete function";
      setError(errorMessage);
    }
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading function...</div>;
  }

  if (!functionData) {
    return <div className="text-center py-8">Function not found</div>;
  }

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8 flex justify-between items-start">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            {functionData.name}
          </h1>
          <p className="mt-2 text-sm text-gray-600">
            {functionData.description || "No description"}
          </p>
          <div className="mt-2 flex items-center gap-4 text-sm text-gray-500">
            <span className="bg-gray-100 px-2 py-1 rounded">
              {functionData.runtime}
            </span>
            <span>{functionData.memory_mb} MB</span>
            <span>{functionData.timeout_sec}s timeout</span>
            <span
              className={`px-2 py-1 rounded ${
                functionData.status === "active"
                  ? "bg-green-100 text-green-800"
                  : "bg-gray-100 text-gray-800"
              }`}
            >
              {functionData.status}
            </span>
          </div>
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => router.push(`/dashboard/functions/${id}/edit`)}
            className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium"
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md font-medium"
          >
            Delete
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* Code Display/Editor */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-semibold text-gray-900">Function Code</h2>
          {!isEditing ? (
            <button
              onClick={() => setIsEditing(true)}
              className="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Edit Code
            </button>
          ) : (
            <div className="flex gap-2">
              <button
                onClick={() => {
                  setEditedCode(functionData.code);
                  setIsEditing(false);
                }}
                className="px-3 py-1 border border-gray-300 rounded text-gray-700 hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                onClick={handleSave}
                className="px-3 py-1 bg-indigo-600 hover:bg-indigo-700 text-white rounded"
              >
                Save
              </button>
            </div>
          )}
        </div>
        <div className="relative">
          <textarea
            value={isEditing ? editedCode : functionData.code}
            onChange={(e) => setEditedCode(e.target.value)}
            readOnly={!isEditing}
            className="w-full h-96 font-mono text-sm p-4 bg-gray-900 text-white border border-gray-700 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
      </div>

      {/* Execute Function */}
      <div className="bg-white rounded-lg shadow p-6 mb-8">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Test Function
        </h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Input (JSON)
            </label>
            <textarea
              value={executionInput}
              onChange={(e) => setExecutionInput(e.target.value)}
              className="w-full h-32 font-mono text-sm p-3 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
              placeholder='{"key": "value"}'
            />
          </div>
          <button
            onClick={handleExecute}
            disabled={isExecuting}
            className="w-full px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isExecuting ? "Executing..." : "Execute Function"}
          </button>

          {executionResult && (
            <div className="mt-4 p-4 bg-gray-50 rounded-md border border-gray-200">
              <h3 className="text-sm font-semibold text-gray-700 mb-2">
                Result
              </h3>
              <pre className="text-sm text-gray-900 overflow-x-auto">
                {JSON.stringify(executionResult, null, 2)}
              </pre>
            </div>
          )}
        </div>
      </div>

      {/* Execution History */}
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-lg font-semibold text-gray-900 mb-4">
          Recent Executions
        </h2>
        {executions.length === 0 ? (
          <p className="text-gray-600 text-sm">No executions yet</p>
        ) : (
          <div className="space-y-2">
            {executions.slice(0, 10).map((execution) => (
              <div
                key={execution.id}
                className="p-4 border border-gray-200 rounded-md hover:bg-gray-50 cursor-pointer"
                onClick={() =>
                  router.push(`/dashboard/executions/${execution.id}`)
                }
              >
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-4">
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
                    <span className="text-sm text-gray-600">
                      {execution.duration_ms}ms
                    </span>
                  </div>
                  <span className="text-xs text-gray-500">
                    {new Date(execution.created_at).toLocaleString()}
                  </span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
