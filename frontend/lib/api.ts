const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api";

export interface User {
  id: string;
  email: string;
  name: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  name: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor() {
    this.baseURL = API_BASE_URL;
    // Load token from localStorage on client side
    if (typeof window !== "undefined") {
      this.token = localStorage.getItem("token");
    }
  }

  setToken(token: string | null) {
    this.token = token;
    if (typeof window !== "undefined") {
      if (token) {
        localStorage.setItem("token", token);
      } else {
        localStorage.removeItem("token");
      }
    }
  }

  getToken() {
    return this.token;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (this.token) {
      headers["Authorization"] = `Bearer ${this.token}`;
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, {
      ...options,
      headers,
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ error: "Unknown error" }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  }

  // Auth endpoints
  async register(data: RegisterRequest): Promise<AuthResponse> {
    return this.request<AuthResponse>("/auth/register", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await this.request<AuthResponse>("/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
    this.setToken(response.token);
    return response;
  }

  logout() {
    this.setToken(null);
  }

  async getCurrentUser(): Promise<{ user: User }> {
    return this.request<{ user: User }>("/auth/me");
  }

  // Function endpoints
  async listFunctions() {
    return this.request("/functions");
  }

  async getFunction(id: string) {
    return this.request(`/functions/${id}`);
  }

  async createFunction(data: any) {
    return this.request("/functions", {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async updateFunction(id: string, data: any) {
    return this.request(`/functions/${id}`, {
      method: "PUT",
      body: JSON.stringify(data),
    });
  }

  async deleteFunction(id: string) {
    return this.request(`/functions/${id}`, {
      method: "DELETE",
    });
  }

  async executeFunction(id: string, input: any) {
    return this.request(`/functions/${id}/execute`, {
      method: "POST",
      body: JSON.stringify({ input }),
    });
  }

  // Execution endpoints
  async listExecutions() {
    return this.request("/executions");
  }

  async getExecution(id: string) {
    return this.request(`/executions/${id}`);
  }

  async getExecutionLogs(id: string) {
    return this.request(`/executions/${id}/logs`);
  }

  // API Key endpoints
  async listAPIKeys() {
    return this.request("/keys");
  }

  async createAPIKey(name: string) {
    return this.request("/keys", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
  }

  async deleteAPIKey(id: string) {
    return this.request(`/keys/${id}`, {
      method: "DELETE",
    });
  }
}

export const apiClient = new ApiClient();
