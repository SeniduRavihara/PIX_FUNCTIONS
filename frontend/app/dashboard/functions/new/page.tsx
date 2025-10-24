"use client";

import { apiClient } from "@/lib/api";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function NewFunctionPage() {
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState("");

  const [formData, setFormData] = useState({
    name: "",
    description: "",
    runtime: "nodejs20",
    code: `// Your function code here
exports.handler = async (event) => {
  console.log('Event:', event);
  
  return {
    statusCode: 200,
    body: JSON.stringify({
      message: 'Hello from VoltRun!',
      input: event
    })
  };
};`,
    entryPoint: "index.handler",
    memoryMb: 128,
    timeoutSec: 30,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsSubmitting(true);

    try {
      const newFunction = await apiClient.createFunction({
        name: formData.name,
        description: formData.description,
        runtime: formData.runtime,
        code: formData.code,
        entry_point: formData.entryPoint,
        memory_mb: formData.memoryMb,
        timeout_sec: formData.timeoutSec,
      });

      router.push(`/dashboard/functions/${newFunction.id}`);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to create function";
      setError(errorMessage);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-5xl">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Create Function</h1>
        <p className="mt-2 text-sm text-gray-600">
          Deploy a new serverless function
        </p>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Basic Information */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Basic Information
          </h2>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Function Name *
              </label>
              <input
                type="text"
                required
                value={formData.name}
                onChange={(e) =>
                  setFormData({ ...formData, name: e.target.value })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900 placeholder-gray-400"
                placeholder="my-function"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Description
              </label>
              <textarea
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900 placeholder-gray-400"
                placeholder="What does this function do?"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Runtime *
              </label>
              <select
                required
                value={formData.runtime}
                onChange={(e) =>
                  setFormData({ ...formData, runtime: e.target.value })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
              >
                <option value="nodejs20">Node.js 20</option>
                <option value="nodejs18">Node.js 18</option>
                <option value="python311">Python 3.11</option>
                <option value="python39">Python 3.9</option>
                <option value="go121">Go 1.21</option>
              </select>
            </div>
          </div>
        </div>

        {/* Code Editor */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Function Code
          </h2>

          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Entry Point
              </label>
              <input
                type="text"
                value={formData.entryPoint}
                onChange={(e) =>
                  setFormData({ ...formData, entryPoint: e.target.value })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900 placeholder-gray-400"
                placeholder="index.handler"
              />
              <p className="mt-1 text-xs text-gray-500">
                Format: filename.functionName (e.g., index.handler)
              </p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Code *
              </label>
              <textarea
                required
                value={formData.code}
                onChange={(e) =>
                  setFormData({ ...formData, code: e.target.value })
                }
                rows={20}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 font-mono text-sm text-gray-900 placeholder-gray-400"
                placeholder="Enter your function code..."
              />
            </div>
          </div>
        </div>

        {/* Configuration */}
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">
            Configuration
          </h2>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Memory (MB)
              </label>
              <input
                type="number"
                min="128"
                max="10240"
                step="64"
                value={formData.memoryMb}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    memoryMb: parseInt(e.target.value),
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
              />
              <p className="mt-1 text-xs text-gray-500">128 MB - 10,240 MB</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Timeout (seconds)
              </label>
              <input
                type="number"
                min="1"
                max="900"
                value={formData.timeoutSec}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    timeoutSec: parseInt(e.target.value),
                  })
                }
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
              />
              <p className="mt-1 text-xs text-gray-500">1 - 900 seconds</p>
            </div>
          </div>
        </div>

        {/* Actions */}
        <div className="flex justify-end space-x-4">
          <button
            type="button"
            onClick={() => router.back()}
            className="px-6 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isSubmitting ? "Creating..." : "Create Function"}
          </button>
        </div>
      </form>
    </div>
  );
}
