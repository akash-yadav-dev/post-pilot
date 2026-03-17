export type AuthTokens = {
  accessToken: string;
  refreshToken: string;
};

export type AuthUser = {
  user_id: string;
  name: string;
  email: string;
};

type AuthPayload = {
  email: string;
  password: string;
};

type RegisterPayload = AuthPayload & {
  name: string;
};

type AuthApiResponse = {
  user_id: string;
  name: string;
  email: string;
  tokens: {
    access_token: string;
    refresh_token: string;
    token_type: string;
    expires_in: number;
  };
};

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ??
  "http://localhost:8080";

const ACCESS_COOKIE = "postpilot_access_token";
const REFRESH_COOKIE = "postpilot_refresh_token";

function setCookie(name: string, value: string, maxAgeSeconds: number) {
  if (typeof document === "undefined") {
    return;
  }

  document.cookie = `${name}=${encodeURIComponent(value)}; Max-Age=${maxAgeSeconds}; Path=/; SameSite=Lax`;
}

function getCookie(name: string) {
  if (typeof document === "undefined") {
    return null;
  }

  const prefix = `${name}=`;
  const parts = document.cookie.split(";");
  for (const raw of parts) {
    const item = raw.trim();
    if (item.startsWith(prefix)) {
      return decodeURIComponent(item.slice(prefix.length));
    }
  }

  return null;
}

function deleteCookie(name: string) {
  if (typeof document === "undefined") {
    return;
  }

  document.cookie = `${name}=; Max-Age=0; Path=/; SameSite=Lax`;
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    let message = "Request failed";
    try {
      const data = (await response.json()) as { error?: string };
      if (data.error) {
        message = data.error;
      }
    } catch {
      // ignore json parse error
    }
    throw new Error(message);
  }

  return (await response.json()) as T;
}

function storeTokens(tokens: AuthTokens) {
  setCookie(ACCESS_COOKIE, tokens.accessToken, 60 * 60);
  setCookie(REFRESH_COOKIE, tokens.refreshToken, 60 * 60 * 24 * 7);
}

export const AuthService = {
  getAccessToken() {
    return getCookie(ACCESS_COOKIE);
  },

  getRefreshToken() {
    return getCookie(REFRESH_COOKIE);
  },

  isAuthenticated() {
    return Boolean(this.getAccessToken());
  },

  clearTokens() {
    deleteCookie(ACCESS_COOKIE);
    deleteCookie(REFRESH_COOKIE);
  },

  async login(payload: AuthPayload): Promise<AuthUser> {
    const response = await fetch(`${API_BASE}/api/auth/login`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await handleResponse<AuthApiResponse>(response);

    storeTokens({
      accessToken: data.tokens.access_token,
      refreshToken: data.tokens.refresh_token,
    });

    return {
      user_id: data.user_id,
      name: data.name,
      email: data.email,
    };
  },

  async register(payload: RegisterPayload): Promise<AuthUser> {
    const response = await fetch(`${API_BASE}/api/auth/register`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await handleResponse<AuthApiResponse>(response);

    storeTokens({
      accessToken: data.tokens.access_token,
      refreshToken: data.tokens.refresh_token,
    });

    return {
      user_id: data.user_id,
      name: data.name,
      email: data.email,
    };
  },

  async me(): Promise<AuthUser> {
    const token = this.getAccessToken();
    if (!token) {
      throw new Error("Not authenticated");
    }

    const response = await fetch(`${API_BASE}/api/auth/me`, {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    return handleResponse<AuthUser>(response);
  },

  async logout(): Promise<void> {
    const refreshToken = this.getRefreshToken();

    if (refreshToken) {
      try {
        await fetch(`${API_BASE}/api/auth/logout`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });
      } catch {
        // noop, clear client state regardless
      }
    }

    this.clearTokens();
  },
};
