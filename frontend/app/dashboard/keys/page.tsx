"use client";

import { apiClient } from "@/lib/api";
import { useEffect, useState } from "react";

interface APIKey {
  id: string;
  name: string;
  prefix: string;
  last_used?: string;
  expires_at?: string;
  created_at: string;
}

export default function APIKeysPage() {
  const [keys, setKeys] = useState<APIKey[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [isCreating, setIsCreating] = useState(false);
  const [newKeyName, setNewKeyName] = useState("");
  const [newKeyValue, setNewKeyValue] = useState<string | null>(null);
  const [showCreateForm, setShowCreateForm] = useState(false);

  useEffect(() => {
    loadKeys();
  }, []);

  const loadKeys = async () => {
    try {
      const data = await apiClient.listAPIKeys();
      setKeys(Array.isArray(data) ? data : []);
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to load API keys";
      setError(errorMessage);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newKeyName.trim()) return;

    setIsCreating(true);
    setError("");

    try {
      const result = await apiClient.createAPIKey(newKeyName);
      setNewKeyValue(result.key);
      setNewKeyName("");
      loadKeys();
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to create API key";
      setError(errorMessage);
    } finally {
      setIsCreating(false);
    }
  };

  const handleDelete = async (id: string, name: string) => {
    if (!confirm(`Are you sure you want to delete the key "${name}"?`)) return;

    try {
      await apiClient.deleteAPIKey(id);
      setKeys(keys.filter((k) => k.id !== id));
    } catch (err: unknown) {
      const errorMessage =
        err instanceof Error ? err.message : "Failed to delete API key";
      setError(errorMessage);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // Could add toast notification here
  };

  if (isLoading) {
    return <div className="text-center py-8">Loading API keys...</div>;
  }

  return (
    <div>
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">API Keys</h1>
            <p className="mt-2 text-sm text-gray-600">
              Manage API keys for programmatic access to your functions
            </p>
          </div>
          <button
            onClick={() => setShowCreateForm(true)}
            className="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md font-medium"
          >
            + Create API Key
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error}
        </div>
      )}

      {/* New Key Display */}
      {newKeyValue && (
        <div className="bg-green-50 border border-green-200 rounded-lg p-6 mb-6">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              <svg
                className="h-6 w-6 text-green-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <div className="ml-3 flex-1">
              <h3 className="text-lg font-semibold text-green-900">
                API Key Created Successfully
              </h3>
              <p className="mt-1 text-sm text-green-700">
                Save this key securely. It won't be shown again!
              </p>
              <div className="mt-4 bg-white rounded p-3 border border-green-300">
                <div className="flex items-center justify-between">
                  <code className="text-sm text-gray-900 break-all">
                    {newKeyValue}
                  </code>
                  <button
                    onClick={() => copyToClipboard(newKeyValue)}
                    className="ml-4 flex-shrink-0 text-indigo-600 hover:text-indigo-700 font-medium"
                  >
                    Copy
                  </button>
                </div>
              </div>
              <button
                onClick={() => setNewKeyValue(null)}
                className="mt-4 text-sm text-green-700 hover:text-green-800 font-medium"
              >
                I've saved it, dismiss this message
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Create Form Modal */}
      {showCreateForm && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
          onClick={() => setShowCreateForm(false)}
        >
          <div
            className="bg-white rounded-lg shadow-xl max-w-md w-full p-6"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              Create API Key
            </h2>
            <form onSubmit={handleCreate}>
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Key Name
                </label>
                <input
                  type="text"
                  value={newKeyName}
                  onChange={(e) => setNewKeyName(e.target.value)}
                  placeholder="e.g., Production API Key"
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 text-gray-900"
                  required
                />
                <p className="mt-1 text-xs text-gray-500">
                  A descriptive name to help you identify this key
                </p>
              </div>
              <div className="flex justify-end gap-3">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 font-medium"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isCreating}
                  className="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-md font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isCreating ? "Creating..." : "Create Key"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Keys List */}
      {keys.length === 0 ? (
        <div className="bg-white rounded-lg shadow p-12 text-center">
          <div className="text-6xl mb-4">🔑</div>
          <h3 className="text-xl font-semibold text-gray-900 mb-2">
            No API keys yet
          </h3>
          <p className="text-gray-600 mb-6">
            Create an API key to access your functions programmatically
          </p>
          <button
            onClick={() => setShowCreateForm(true)}
            className="inline-block bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-3 rounded-md font-medium"
          >
            Create Your First API Key
          </button>
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Name
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Key
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Created
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Last Used
                </th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              {keys.map((key) => (
                <tr key={key.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="text-sm font-medium text-gray-900">
                      {key.name}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <code className="text-sm text-gray-600 bg-gray-100 px-2 py-1 rounded">
                      {key.prefix}...
                    </code>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {new Date(key.created_at).toLocaleDateString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                    {key.last_used
                      ? new Date(key.last_used).toLocaleDateString()
                      : "Never"}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <button
                      onClick={() => handleDelete(key.id, key.name)}
                      className="text-red-600 hover:text-red-900"
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* API Usage Info */}
      <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-blue-900 mb-2">
          Using API Keys
        </h3>
        <p className="text-sm text-blue-700 mb-4">
          Include your API key in the Authorization header when making requests:
        </p>
        <pre className="bg-blue-900 text-blue-100 p-4 rounded text-sm overflow-x-auto">
          {`curl -X POST https://api.voltrun.com/api/functions/{id}/execute \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{"input": {"key": "value"}}'`}
        </pre>
      </div>
    </div>
  );
}
